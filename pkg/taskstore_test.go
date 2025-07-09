package pkg

import (
	"testing"
	"time"
)

func TestCreateTask(t *testing.T) {
	store := NewTaskStore()

	// Тест на успешное создание задачи
	task, err := store.CreateTask("Test Task")
	if err != nil {
		t.Fatalf("CreateTask() returned error: %v", err)
	}
	if task.Id == "" {
		t.Error("Task ID is empty")
	}
	if task.Name != "Test Task" {
		t.Errorf("Expected task name 'Test Task', got '%s'", task.Name)
	}
	if task.Status != TaskStatusPending {
		t.Errorf("Expected task status '%s', got '%s'", TaskStatusPending, task.Status)
	}
	if task.CreatedAt.IsZero() {
		t.Error("Task CreatedAt is zero")
	}

	// Тест на создание задачи с пустым именем
	_, err = store.CreateTask("")
	if err == nil {
		t.Error("CreateTask() with empty name should return error")
	}
	if err != nil && err.Error() != TaskErrorNameRequired {
		t.Errorf("Expected error '%s', got '%s'", TaskErrorNameRequired, err.Error())
	}
}

func TestGetTask(t *testing.T) {
	store := NewTaskStore()

	// Создаем тестовую задачу
	task, _ := store.CreateTask("Test Task")

	// Тест на получение существующей задачи
	retrievedTask, err := store.GetTask(task.Id)
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
	_, err = store.GetTask("non-existent-id")
	if err == nil {
		t.Error("GetTask() with non-existent ID should return error")
	}
	if err != nil && err.Error() != TaskErrorNotFound {
		t.Errorf("Expected error '%s', got '%s'", TaskErrorNotFound, err.Error())
	}
}

func TestUpdateTask(t *testing.T) {
	store := NewTaskStore()

	// Создаем тестовую задачу
	task, _ := store.CreateTask("Test Task")

	// Изменяем задачу
	task.Status = TaskStatusRunning
	task.StartedAt = time.Now()

	// Тест на обновление существующей задачи
	err := store.UpdateTask(task)
	if err != nil {
		t.Fatalf("UpdateTask() returned error: %v", err)
	}

	// Проверяем, что задача обновилась
	updatedTask, _ := store.GetTask(task.Id)
	if updatedTask.Status != TaskStatusRunning {
		t.Errorf("Expected task status '%s', got '%s'", TaskStatusRunning, updatedTask.Status)
	}
	if updatedTask.StartedAt.IsZero() {
		t.Error("Task StartedAt is zero after update")
	}

	// Тест на обновление несуществующей задачи
	nonExistentTask := InternalTask{Id: "non-existent-id"}
	err = store.UpdateTask(nonExistentTask)
	if err == nil {
		t.Error("UpdateTask() with non-existent ID should return error")
	}
	if err != nil && err.Error() != TaskErrorNotFound {
		t.Errorf("Expected error '%s', got '%s'", TaskErrorNotFound, err.Error())
	}
}

func TestDeleteTask(t *testing.T) {
	store := NewTaskStore()

	// Создаем тестовую задачу
	task, _ := store.CreateTask("Test Task")

	// Тест на удаление существующей задачи
	err := store.DeleteTask(task.Id)
	if err != nil {
		t.Fatalf("DeleteTask() returned error: %v", err)
	}

	// Проверяем, что задача удалена
	_, err = store.GetTask(task.Id)
	if err == nil {
		t.Error("GetTask() after deletion should return error")
	}

	// Тест на удаление несуществующей задачи
	err = store.DeleteTask("non-existent-id")
	if err == nil {
		t.Error("DeleteTask() with non-existent ID should return error")
	}
	if err != nil && err.Error() != TaskErrorNotFound {
		t.Errorf("Expected error '%s', got '%s'", TaskErrorNotFound, err.Error())
	}
}

func TestGetTasks(t *testing.T) {
	store := NewTaskStore()

	// Проверяем, что изначально список пуст
	tasks := store.GetTasks()
	if len(tasks) != 0 {
		t.Errorf("Expected empty task list, got %d tasks", len(tasks))
	}

	// Создаем несколько задач
	store.CreateTask("Task 1")
	store.CreateTask("Task 2")
	store.CreateTask("Task 3")

	// Проверяем, что все задачи возвращаются
	tasks = store.GetTasks()
	if len(tasks) != 3 {
		t.Errorf("Expected 3 tasks, got %d tasks", len(tasks))
	}
}
