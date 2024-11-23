package main

import (
	"context"
	"net/http"

	"github.com/gitwub5/go_todo_app/auth"
	"github.com/gitwub5/go_todo_app/clock"
	"github.com/gitwub5/go_todo_app/config"
	"github.com/gitwub5/go_todo_app/handler"
	"github.com/gitwub5/go_todo_app/service"
	"github.com/gitwub5/go_todo_app/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

// context.Context와 *config.Config를 인자로 받고, http.Handler와 cleanup 함수를 반환
func NewMux(ctx context.Context, cfg *config.Config) (http.Handler, func(), error) {
	mux := chi.NewRouter()

	// /health 요청을 처리하는 핸들러 등록
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, _ = w.Write([]byte(`{"status": "ok"}`))
	})

	// 유효성 검사기 생성
	v := validator.New()

	// 데이터베이스 연결 및 정리 함수 생성
	db, cleanup, err := store.New(ctx, cfg)
	if err != nil {
		return nil, cleanup, err
	}

	// 실제 시간 시계 사용하여 Repository를 생성
	clocker := clock.RealClocker{}
	r := store.Repository{Clocker: clocker}
	// Redis 클라이언트 생성
	rcli, err := store.NewKVS(ctx, cfg)
	if err != nil {
		return nil, cleanup, err
	}
	// JWTer 생성
	jwter, err := auth.NewJWTer(rcli, clocker)
	if err != nil {
		return nil, cleanup, err
	}

	// POST /register 요청을 처리하는 핸들러 등록
	ru := &handler.RegisterUser{
		Service:   &service.RegisterUser{DB: db, Repo: &r},
		Validator: v,
	}
	mux.Post("/register", ru.ServeHTTP)

	// POST /login 요청을 처리하는 핸들러 등록
	l := &handler.Login{
		Service: &service.Login{
			DB:             db,
			Repo:           &r,
			TokenGenerator: jwter,
		},
		Validator: v,
	}
	mux.Post("/login", l.ServeHTTP)

	// POST /tasks 요청을 처리하는 핸들러
	at := &handler.AddTask{
		Service:   &service.AddTask{DB: db, Repo: &r},
		Validator: v,
	}
	// GET /tasks 요청 처리하는 핸들러
	lt := &handler.ListTask{
		Service: &service.ListTask{DB: db, Repo: &r},
	}

	mux.Route("/tasks", func(r chi.Router) {
		r.Use(handler.AuthMiddleware(jwter)) // /tasks 하위 모든 요청에 대해 인증 미들웨어 적용
		r.Post("/", at.ServeHTTP)            // POST /tasks 요청을 처리하는 핸들러 등록
		r.Get("/", lt.ServeHTTP)             // GET /tasks 요청 처리하는 핸들러 등록
	})

	// /admin 권한 사용자만 접속할 수 있는 엔드포인트
	mux.Route("/admin", func(r chi.Router) {
		r.Use(handler.AuthMiddleware(jwter), handler.AdminMiddleware)
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			_, _ = w.Write([]byte(`{"message": "admin only"}`))
		})
	})

	return mux, cleanup, nil
}
