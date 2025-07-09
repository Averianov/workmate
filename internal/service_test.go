package internal

import (
	"context"
	"testing"
	"time"
	"workmate/pkg"
)

func TestCreateTask(t *testing.T) {
	service := NewService()
	ctx := context.Background()

	// Тест на успешное создание задачи
	task, err := service.CreateTask(ctx, "Test Task")
	if err != nil {
		t.Fatalf("CreateTask() returned error: %v", err)
	}
	if task.Id == "" {
		t.Error("Task ID is empty")
	}
	if task.Name != "Test Task" {
		t.Errorf("Expected task name 'Test Task', got '%s'", task.Name)
	}
	if task.Status != pkg.TaskStatusPending {
		t.Errorf("Expected task status '%s', got '%s'", pkg.TaskStatusPending, task.Status)
	}

	// Проверяем, что задача появилась в хранилище
	tasks, _ := service.GetTasks(ctx)
	if len(tasks) != 1 {
		t.Errorf("Expected 1 task in store, got %d", len(tasks))
	}

	// Тест на создание задачи с пустым именем
	_, err = service.CreateTask(ctx, "")
	if err == nil {
		t.Error("CreateTask() with empty name should return error")
	}
}

func TestGetTask(t *testing.T) {
	service := NewService()
	ctx := context.Background()

	// Создаем тестовую задачу
	task, _ := service.CreateTask(ctx, "Test Task")

	// Тест на получение существующей задачи
	retrievedTask, err := service.GetTask(ctx, task.Id)
	if err != nil {
		t.Fatalf("GetTask() returned error: %v", err)
	}
	if retrievedTask.Id != task.Id {
		t.Errorf("Expected task ID '%s', got '%s'", task.Id, retrievedTask.Id)
	}
	if retrievedTask.Name != task.Name {
		t.Errorf("Expected task name '%s', got '%s'", task.Name, retrievedTask.Name)
	}

	// Тест на получение несуществующей задачи
	_, err = service.GetTask(ctx, "non-existent-id")
	if err == nil {
		t.Error("GetTask() with non-existent ID should return error")
	}
}

func TestGetTasks(t *testing.T) {
	service := NewService()
	ctx := context.Background()

	// Проверяем, что изначально список пуст
	tasks, err := service.GetTasks(ctx)
	if err != nil {
		t.Fatalf("GetTasks() returned error: %v", err)
	}
	if len(tasks) != 0 {
		t.Errorf("Expected empty task list, got %d tasks", len(tasks))
	}

	// Создаем несколько задач
	service.CreateTask(ctx, "Task 1")
	service.CreateTask(ctx, "Task 2")

	// Проверяем, что все задачи возвращаются
	tasks, err = service.GetTasks(ctx)
	if err != nil {
		t.Fatalf("GetTasks() returned error: %v", err)
	}
	if len(tasks) != 2 {
		t.Errorf("Expected 2 tasks, got %d tasks", len(tasks))
	}
}

func TestDeleteTask(t *testing.T) {
	service := NewService()
	ctx := context.Background()

	// Создаем тестовую задачу
	task, _ := service.CreateTask(ctx, "Test Task")

	// Тест на удаление существующей задачи
	err := service.DeleteTask(ctx, task.Id)
	if err != nil {
		t.Fatalf("DeleteTask() returned error: %v", err)
	}

	// Проверяем, что задача удалена
	_, err = service.GetTask(ctx, task.Id)
	if err == nil {
		t.Error("GetTask() after deletion should return error")
	}

	// Тест на удаление несуществующей задачи
	err = service.DeleteTask(ctx, "non-existent-id")
	if err == nil {
		t.Error("DeleteTask() with non-existent ID should return error")
	}
}

func TestGetTaskResult(t *testing.T) {
	service := NewService()
	ctx := context.Background()

	// Создаем тестовую задачу
	task, _ := service.CreateTask(ctx, "Test Task")

	// Пытаемся получить результат незавершенной задачи
	_, err := service.GetTaskResult(ctx, task.Id)
	if err == nil {
		t.Error("GetTaskResult() for pending task should return error")
	}

	// Имитируем завершение задачи
	task.Status = pkg.TaskStatusCompleted
	task.Result = "Test result"
	service.store.UpdateTask(task)

	// Получаем результат завершенной задачи
	completedTask, err := service.GetTaskResult(ctx, task.Id)
	if err != nil {
		t.Fatalf("GetTaskResult() returned error: %v", err)
	}
	if completedTask.Result != "Test result" {
		t.Errorf("Expected result 'Test result', got '%s'", completedTask.Result)
	}
}

func TestSimulateIOTask(t *testing.T) {
	service := NewService()
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	// Создаем тестовую задачу
	task, _ := service.CreateTask(ctx, "Test Task")

	// Проверяем, что задача изначально в статусе pending
	retrievedTask, _ := service.GetTask(ctx, task.Id)
	if retrievedTask.Status != pkg.TaskStatusPending {
		t.Errorf("Expected task status '%s', got '%s'", pkg.TaskStatusPending, retrievedTask.Status)
	}

	// Запускаем симуляцию задачи
	service.simulateIOTask(ctx, task.Id)

	// Проверяем, что задача перешла в статус running или failed
	retrievedTask, _ = service.GetTask(ctx, task.Id)
	if retrievedTask.Status != pkg.TaskStatusRunning && retrievedTask.Status != pkg.TaskStatusFailed {
		t.Errorf("Expected task status '%s' or '%s', got '%s'",
			pkg.TaskStatusRunning, pkg.TaskStatusFailed, retrievedTask.Status)
	}
}
