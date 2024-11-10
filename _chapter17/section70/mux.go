package main

import (
	"net/http"

	"github.com/gitwub5/go_todo_app/handler"
	"github.com/gitwub5/go_todo_app/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

// 유연한 라우팅을 위해 go-chi/chi 패키지를 사용하여 라우터를 생성하도록 변경
func NewMux() http.Handler {
	mux := chi.NewRouter()

	// /health 요청을 처리하는 핸들러 등록
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, _ = w.Write([]byte(`{"status": "ok"}`))
	})

	v := validator.New()

	// /tasks 엔드포인트를 처리하는 AddTask 핸들러 등록
	mux.Handle("/tasks", &handler.AddTask{Store: store.Tasks, Validator: v})

	// POST /tasks 요청을 처리하는 핸들러 등록
	at := &handler.AddTask{Store: store.Tasks, Validator: v}
	mux.Post("/tasks", at.ServeHTTP)

	// GET /tasks 요청을 처리하는 핸들러 등록
	lt := &handler.ListTask{Store: store.Tasks}
	mux.Get("/tasks", lt.ServeHTTP)

	// 설정된 라우터를 반환
	return mux
}
