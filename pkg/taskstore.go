package pkg

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

// TaskStore хранит задачи в памяти
type TaskStore struct {
	mu    sync.RWMutex
	tasks map[string]*InternalTask
}

// NewTaskStore создает новое хранилище задач
func NewTaskStore() *TaskStore {
	return &TaskStore{
		tasks: make(map[string]*InternalTask),
	}
}

// GetTasks - Получить список всех задач
func (s *TaskStore) GetTasks() (tasks []InternalTask) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tasks = make([]InternalTask, 0, len(s.tasks))
	for _, task := range s.tasks {
		tasks = append(tasks, *task)
	}

	return
}

// CreateTask - Создать новую задачу (только сохранение в хранилище)
func (s *TaskStore) CreateTask(taskName string) (task InternalTask, err error) {
	if taskName == "" {
		err = fmt.Errorf("%s", TaskErrorNameRequired)
		return
	}

	newTask := InternalTask{
		Id:        uuid.New().String(),
		Name:      taskName,
		Status:    TaskStatusPending,
		CreatedAt: time.Now(),
	}

	s.mu.Lock()
	s.tasks[newTask.Id] = &newTask
	s.mu.Unlock()
	task = newTask

	return
}

// GetTask - Получить задачу по ID
func (s *TaskStore) GetTask(taskId string) (task InternalTask, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	targetTask, exists := s.tasks[taskId]
	if !exists {
		err = fmt.Errorf("%s", TaskErrorNotFound)
		return
	}
	task = *targetTask
	return
}

// UpdateTask - Обновить задачу
func (s *TaskStore) UpdateTask(task InternalTask) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.tasks[task.Id]; !exists {
		return fmt.Errorf("%s", TaskErrorNotFound)
	}

	s.tasks[task.Id] = &task
	return nil
}

// DeleteTask - Удалить задачу
func (s *TaskStore) DeleteTask(taskId string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.tasks[taskId]; !exists {
		return fmt.Errorf("%s", TaskErrorNotFound)
	}

	delete(s.tasks, taskId)
	return nil
}
