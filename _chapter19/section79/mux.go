package main

import (
	"context"
	"net/http"

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
	r := store.Repository{Clocker: clock.RealClocker{}}

	// POST /tasks 요청을 처리하는 핸들러 등록 (비즈니스 로직과 데이터베이스 처리를 분리)
	at := &handler.AddTask{
		Service:   &service.AddTask{DB: db, Repo: &r},
		Validator: v,
	}
	mux.Post("/tasks", at.ServeHTTP)

	// GET /tasks 요청 처리하는 핸들러 등록 (비즈니스 로직과 데이터베이스 처리를 분리)
	lt := &handler.ListTask{
		Service: &service.ListTask{DB: db, Repo: &r},
	}
	mux.Get("/tasks", lt.ServeHTTP)

	// POST /register 요청을 처리하는 핸들러 등록
	ru := &handler.RegisterUser{
		Service:   &service.RegisterUser{DB: db, Repo: &r},
		Validator: v,
	}
	mux.Post("/register", ru.ServeHTTP)

	return mux, cleanup, nil
}
