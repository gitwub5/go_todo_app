package store

import (
	"errors"

	"github.com/gitwub5/go_todo_app/entity"
)

var (
	Tasks       = &TaskStore{Tasks: map[entity.TaskID]*entity.Task{}} // TaskStore 구조체를 사용하여 Task를 저장하는 저장소 생성
	ErrNotFound = errors.New("not found")
)

type TaskStore struct {
	// 동적 확인용이므로 일부러 export 하고 있다.
	LastID entity.TaskID
	Tasks  map[entity.TaskID]*entity.Task
}

// Add 메서드는 Task를 저장소에 추가한다.
func (ts *TaskStore) Add(t *entity.Task) (entity.TaskID, error) {
	ts.LastID++
	t.ID = ts.LastID
	ts.Tasks[t.ID] = t
	return t.ID, nil
}

// Get 메서드는 Task를 저장소에서 가져온다.
func (ts *TaskStore) Get(id entity.TaskID) (*entity.Task, error) {
	if ts, ok := ts.Tasks[id]; ok {
		return ts, nil
	}
	return nil, ErrNotFound
}

// All 메서드는 정렬이 끝난 Task 목록를 반환한다.
func (ts *TaskStore) All() entity.Tasks {
	tasks := make([]*entity.Task, len(ts.Tasks))
	for i, t := range ts.Tasks {
		tasks[i-1] = t
	}
	return tasks
}
