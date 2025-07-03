package openapi

import (
	"workmate/pkg"
)

// Маппинг из внутреннего типа сервиса в API DTO
func MapInternalTaskToAPI(task pkg.InternalTask) Task {
	return Task{
		Id:         task.Id,
		Status:     task.Status,
		CreatedAt:  task.CreatedAt,
		StartedAt:  task.StartedAt,
		FinishedAt: task.FinishedAt,
		Result:     task.Result,
		Error:      task.Error,
		Duration:   task.Duration,
	}
}

// Маппинг списка задач из внутренних типов сервиса в API
func MapInternalTasksToAPI(tasks []pkg.InternalTask) []Task {
	result := make([]Task, len(tasks))
	for i, task := range tasks {
		result[i] = MapInternalTaskToAPI(task)
	}
	return result
}
