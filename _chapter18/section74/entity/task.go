package entity

import (
	"time"
)

type TaskID int64      // Task의 ID를 나타내는 타입
type TaskStatus string // Task의 상태를 나타내는 타입

// TaskStatus 상수
const (
	TaskStatusTodo  TaskStatus = "todo"
	TaskStatusDoing TaskStatus = "doing"
	TaskStatusDone  TaskStatus = "done"
)

// Task 구조체는 할 일을 나타내는 구조체이다.
type Task struct {
	ID       TaskID     `json:"id" db:"id"`
	Title    string     `json:"title" db:"title"`
	Status   TaskStatus `json:"status" db:"status"`
	Created  time.Time  `json:"created" db:"created"`
	Modified time.Time  `json:"modified" db:"modified"`
}

// Tasks는 Task의 슬라이스이다.
type Tasks []*Task
