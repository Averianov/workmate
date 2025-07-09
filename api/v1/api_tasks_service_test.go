package openapi

import (
	"context"
	"testing"
	"workmate/pkg"
)

// Вспомогательная функция для проверки кода ответа
func assertResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d, got %d", expected, actual)
	}
}

// Вспомогательная функция для проверки ошибки
func assertError(t *testing.T, expected string, resp ImplResponse) {
	errResp, ok := resp.Body.(ErrorResponse)
	if !ok {
		t.Fatalf("Expected ErrorResponse, got %T", resp.Body)
	}
	if errResp.Error != expected {
		t.Errorf("Expected error '%s', got '%s'", expected, errResp.Error)
	}
}

func TestCreateTask(t *testing.T) {
	service := NewTasksAPIService()
	ctx := context.Background()

	// Тест на успешное создание задачи
	resp, err := service.CreateTask(ctx, CreateTaskRequest{Name: "Test Task"})
	if err != nil {
		t.Fatalf("CreateTask() returned error: %v", err)
	}
	assertResponseCode(t, 201, resp.Code)

	taskResp, ok := resp.Body.(TaskResponse)
	if !ok {
		t.Fatalf("Expected TaskResponse, got %T", resp.Body)
	}
	if taskResp.Task.Status != pkg.TaskStatusPending {
		t.Errorf("Expected task status '%s', got '%s'", pkg.TaskStatusPending, taskResp.Task.Status)
	}

	// Тест на создание задачи с пустым именем
	resp, err = service.CreateTask(ctx, CreateTaskRequest{Name: ""})
	if err != nil {
		t.Fatalf("CreateTask() returned error: %v", err)
	}
	assertResponseCode(t, 400, resp.Code)
	assertError(t, pkg.TaskErrorNameRequired, resp)
}

func TestGetTask(t *testing.T) {
	service := NewTasksAPIService()
	ctx := context.Background()

	// Создаем тестовую задачу
	createResp, _ := service.CreateTask(ctx, CreateTaskRequest{Name: "Test Task"})
	taskResp := createResp.Body.(TaskResponse)
	taskId := taskResp.Task.Id

	// Тест на получение существующей задачи
	resp, err := service.GetTask(ctx, taskId)
	if err != nil {
		t.Fatalf("GetTask() returned error: %v", err)
	}
	assertResponseCode(t, 200, resp.Code)

	getTaskResp, ok := resp.Body.(TaskResponse)
	if !ok {
		t.Fatalf("Expected TaskResponse, got %T", resp.Body)
	}
	if getTaskResp.Task.Id != taskId {
		t.Errorf("Expected task ID '%s', got '%s'", taskId, getTaskResp.Task.Id)
	}

	// Тест на получение несуществующей задачи
	resp, err = service.GetTask(ctx, "non-existent-id")
	if err != nil {
		t.Fatalf("GetTask() returned error: %v", err)
	}
	assertResponseCode(t, 404, resp.Code)
	assertError(t, pkg.TaskErrorNotFound, resp)
}

func TestGetTasks(t *testing.T) {
	service := NewTasksAPIService()
	ctx := context.Background()

	// Проверяем, что изначально список пуст
	resp, err := service.GetTasks(ctx)
	if err != nil {
		t.Fatalf("GetTasks() returned error: %v", err)
	}
	assertResponseCode(t, 200, resp.Code)

	listResp, ok := resp.Body.(TaskListResponse)
	if !ok {
		t.Fatalf("Expected TaskListResponse, got %T", resp.Body)
	}

	initialCount := len(listResp.Tasks)

	// Создаем несколько задач
	service.CreateTask(ctx, CreateTaskRequest{Name: "Task 1"})
	service.CreateTask(ctx, CreateTaskRequest{Name: "Task 2"})

	// Проверяем, что все задачи возвращаются
	resp, err = service.GetTasks(ctx)
	if err != nil {
		t.Fatalf("GetTasks() returned error: %v", err)
	}
	assertResponseCode(t, 200, resp.Code)

	listResp, ok = resp.Body.(TaskListResponse)
	if !ok {
		t.Fatalf("Expected TaskListResponse, got %T", resp.Body)
	}
	if len(listResp.Tasks) != initialCount+2 {
		t.Errorf("Expected %d tasks, got %d tasks", initialCount+2, len(listResp.Tasks))
	}
}

func TestDeleteTask(t *testing.T) {
	service := NewTasksAPIService()
	ctx := context.Background()

	// Создаем тестовую задачу
	createResp, _ := service.CreateTask(ctx, CreateTaskRequest{Name: "Test Task"})
	taskResp := createResp.Body.(TaskResponse)
	taskId := taskResp.Task.Id

	// Тест на удаление существующей задачи
	resp, err := service.DeleteTask(ctx, taskId)
	if err != nil {
		t.Fatalf("DeleteTask() returned error: %v", err)
	}
	assertResponseCode(t, 204, resp.Code)

	// Проверяем, что задача удалена
	getResp, _ := service.GetTask(ctx, taskId)
	assertResponseCode(t, 404, getResp.Code)

	// Тест на удаление несуществующей задачи
	resp, err = service.DeleteTask(ctx, "non-existent-id")
	if err != nil {
		t.Fatalf("DeleteTask() returned error: %v", err)
	}
	assertResponseCode(t, 404, resp.Code)
	assertError(t, pkg.TaskErrorNotFound, resp)
}

func TestGetTaskResult(t *testing.T) {
	service := NewTasksAPIService()
	ctx := context.Background()

	// Создаем тестовую задачу
	createResp, _ := service.CreateTask(ctx, CreateTaskRequest{Name: "Test Task"})
	taskResp := createResp.Body.(TaskResponse)
	taskId := taskResp.Task.Id

	// Пытаемся получить результат незавершенной задачи
	resp, err := service.GetTaskResult(ctx, taskId)
	if err != nil {
		t.Fatalf("GetTaskResult() returned error: %v", err)
	}
	assertResponseCode(t, 425, resp.Code)
	assertError(t, pkg.TaskErrorNotCompleted, resp)

	// Тест на получение результата несуществующей задачи
	resp, err = service.GetTaskResult(ctx, "non-existent-id")
	if err != nil {
		t.Fatalf("GetTaskResult() returned error: %v", err)
	}
	assertResponseCode(t, 404, resp.Code)
	assertError(t, pkg.TaskErrorNotFound, resp)
}
