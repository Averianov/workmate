package internal

import (
	"context"
	"fmt"
	"math/rand"
	"time"
	"workmate/pkg"
)

type Service struct {
	store *pkg.TaskStore
}

func NewService() *Service {
	return &Service{
		store: pkg.NewTaskStore(),
	}
}

// simulateIOTask симулирует длительную I/O операцию
func (s *Service) simulateIOTask(ctx context.Context, taskId string) {
	// Получаем задачу из хранилища
	task, err := s.store.GetTask(taskId)
	if err != nil {
		return
	}

	// Обновляем статус на running
	now := time.Now()
	task.StartedAt = now
	task.Status = pkg.TaskStatusRunning

	// Сохраняем изменения
	if err := s.store.UpdateTask(task); err != nil {
		return
	}

	// Симулируем работу от 3 до 5 минут
	duration := time.Duration(rand.Intn(121)+180) * time.Second

	select {
	case <-time.After(duration):
		// Задача успешно завершена
		finishedAt := time.Now()
		task.FinishedAt = finishedAt
		task.Status = pkg.TaskStatusCompleted
		task.Result = fmt.Sprintf("Task completed successfully after %s", duration.Round(time.Second))
		task.Duration = finishedAt.Sub(task.StartedAt).Round(time.Second).String()
		s.store.UpdateTask(task)

	case <-ctx.Done():
		// Задача была отменена
		finishedAt := time.Now()
		task.FinishedAt = finishedAt
		task.Status = pkg.TaskStatusFailed
		task.Error = "Task was cancelled"
		task.Duration = finishedAt.Sub(task.StartedAt).Round(time.Second).String()
		s.store.UpdateTask(task)
	}
}

// GetTasks - Получить список всех задач с обновленной продолжительностью
func (s *Service) GetTasks(ctx context.Context) ([]pkg.InternalTask, error) {
	tasks := s.store.GetTasks()

	// Обновляем duration для запущенных задач
	for i := range tasks {
		if tasks[i].Status == pkg.TaskStatusRunning && !tasks[i].StartedAt.IsZero() {
			duration := time.Since(tasks[i].StartedAt)
			tasks[i].Duration = duration.Round(time.Second).String()
			// Обновляем в хранилище для консистентности
			s.store.UpdateTask(tasks[i])
		}
	}

	return tasks, nil
}

// CreateTask - Создать новую задачу и запустить её выполнение
func (s *Service) CreateTask(ctx context.Context, taskName string) (task pkg.InternalTask, err error) {
	task, err = s.store.CreateTask(taskName)
	if err != nil {
		return
	}

	// Запускаем задачу в фоне
	go s.simulateIOTask(context.Background(), task.Id)

	return
}

// GetTask - Получить информацию о задаче с обновленной продолжительностью
func (s *Service) GetTask(ctx context.Context, taskId string) (task pkg.InternalTask, err error) {
	task, err = s.store.GetTask(taskId)
	if err != nil {
		return
	}

	// Обновляем duration если задача еще выполняется
	if task.Status == pkg.TaskStatusRunning && !task.StartedAt.IsZero() {
		duration := time.Since(task.StartedAt)
		task.Duration = duration.Round(time.Second).String()
		// Обновляем в хранилище
		s.store.UpdateTask(task)
	}

	return
}

// DeleteTask - Удалить задачу
func (s *Service) DeleteTask(ctx context.Context, taskId string) error {
	return s.store.DeleteTask(taskId)
}

// GetTaskResult - Получить результат задачи (бизнес-логика проверки готовности)
func (s *Service) GetTaskResult(ctx context.Context, taskId string) (task pkg.InternalTask, err error) {
	task, err = s.store.GetTask(taskId)
	if err != nil {
		return
	}

	if task.Status != pkg.TaskStatusCompleted && task.Status != pkg.TaskStatusFailed {
		err = fmt.Errorf("%s", pkg.TaskErrorNotCompleted)
		return
	}

	return
}
