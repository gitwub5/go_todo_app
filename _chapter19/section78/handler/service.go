package handler

import (
	"context"

	"github.com/gitwub5/go_todo_app/entity"
)

// handler 패키지로부터 비즈니스 로직과 데이터베이스 처리를 제외시키기 위해 서비스 인터페이스를 정의한다.

//go:generate go run github.com/matryer/moq -out moq_test.go . ListTasksService AddTaskService
type ListTasksService interface {
	ListTasks(ctx context.Context) (entity.Tasks, error)
}
type AddTaskService interface {
	AddTask(ctx context.Context, title string) (*entity.Task, error)
}
