package pkg

import "time"

// TaskStatus - статус задачи
const (
	TaskStatusPending   = "pending"
	TaskStatusRunning   = "running"
	TaskStatusCompleted = "completed"
	TaskStatusFailed    = "failed"
)

// TaskError - ошибка задачи
const (
	TaskErrorNameRequired = "Task name is required"
	TaskErrorNotFound     = "Task not found"
	TaskErrorNotCompleted = "Task is not completed yet"
)

// InternalTask - внутренняя сущность задачи
type InternalTask struct {
	Id         string    `json:"id"`
	Name       string    `json:"name"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"createdAt"`
	StartedAt  time.Time `json:"startedAt,omitempty"`
	FinishedAt time.Time `json:"finishedAt,omitempty"`
	Result     string    `json:"result,omitempty"`
	Error      string    `json:"error,omitempty"`
	Duration   string    `json:"duration,omitempty"`
}
