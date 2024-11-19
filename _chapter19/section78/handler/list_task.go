package handler

import (
	"net/http"

	"github.com/gitwub5/go_todo_app/entity"
)

// ListTask는 Task 목록을 반환하는 핸들러이다.
type ListTask struct {
	Service ListTasksService
}

type task struct {
	ID     entity.TaskID     `json:"id"`
	Title  string            `json:"title"`
	Status entity.TaskStatus `json:"status"`
}

// ServeHTTP는 HTTP 요청을 처리하는 메서드로, ListTask 핸들러의 엔트리 포인트이다. (GET /tasks)
func (lt *ListTask) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tasks, err := lt.Service.ListTasks(ctx)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}
	// 등록이 끝난 모든 Task 목록을 JSON 응답으로 변환한다.
	rsp := []task{}
	for _, t := range tasks {
		rsp = append(rsp, task{
			ID:     t.ID,
			Title:  t.Title,
			Status: t.Status,
		})
	}
	RespondJSON(ctx, w, rsp, http.StatusOK)
}
