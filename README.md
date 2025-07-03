# WorkMate Task Manager

API сервис для управления длительными I/O задачами с поддержкой HTTP/2.

## Описание

WorkMate - это HTTP/2 сервис, который позволяет:
- Создавать длительные I/O задачи (выполняются 3-5 минут)
- Отслеживать статус выполнения задач
- Получать результаты выполненных задач
- Удалять задачи
- Получать список всех задач

Все данные хранятся в памяти сервиса.

## Технологии

- **Go 1.21**
- **HTTP/2** с TLS
- **OpenAPI 3.0** спецификация
- **OpenAPI Generator** для генерации серверного кода
- **gorilla/mux** - HTTP роутер
- **UUID** - генерация уникальных идентификаторов

## Структура проекта

```
workmate/
├── api/contract/v1/
│   └── workmate.yaml         # OpenAPI спецификация
├── api/v1/                   # Сгенерированный код
│   ├── common.go             # Маппинг сущностей для текущей версии апи (создан в ручную)
│   ├── api_tasks_service.go  # Контроллер для бизнес-логики (правится в ручную)
│   ├── api_tasks.go          # HTTP handlers
│   ├── model_*.go            # Модели данных
│   └── routers.go            # Роутинг
├── cmd/  
│   ├── launcher/
│         └── main.go         # Лаунчер для сервера  
│   ├── swagger/
│         └── main.go         # Лаунчер для сваггера
├── internal/  
│   ├── service.go            # Бизнес-логика   
├── pkg/  
│   ├── common.go             # Сущности и константы   
│   ├── taskstore.go          # Хранилище в памяти       
├── script/  
│   ├── gen-certs.sh          # Скрипт генерации сертификатов
├── go.mod                    # Go модуль
├── Makefile                  # Команды сборки
└── README.md                 # Документация
```

## Установка и запуск

### Предварительные требования

- Go 1.21 или выше
- npm (для openapi-generator-cli)
- openssl (для генерации сертификатов, опционально)
- Make (опционально)

### Установка openapi-generator-cli

```bash
npm install -g @openapitools/openapi-generator-cli
```

### Установка зависимостей

```bash
cd workmate
make deps
```

### Генерация серверного кода (при изменении OpenAPI спецификации)

```bash
make generate
# или напрямую:
openapi-generator-cli generate -g go-server -i ./v1/workmate.yaml --git-repo-id workmate -o ./generated
```

### Генерация TLS сертификатов для HTTP/2 (опционально)

```bash
make certs
# или:
chmod +x gen-certs.sh
./gen-certs.sh
```

### Запуск сервиса

#### С HTTP/2 (требуются сертификаты):
```bash
make certs  # генерация сертификатов (опционально)
make run    # запуск сервера
```

### Сборка исполняемого файла

```bash
make build
./workmate
```

## API Документация

Для изучения и тестирования API используйте встроенный Swagger UI интерфейс.

### Запуск Swagger UI

```bash
make swagger
```

### Доступные адреса

- ** Swagger UI интерфейс**: http://localhost:8088/
- ** Альтернативный URL**: http://localhost:8088/swagger/
- ** OpenAPI спецификация**: http://localhost:8088/swagger.yaml
- ** Health check**: http://localhost:8088/health

### API Endpoints

**Базовый URL для API:** `https://localhost:8080/api/v1`

Доступные операции:
- **POST** `/tasks` - Создать новую задачу
- **GET** `/tasks` - Получить список всех задач
- **GET** `/tasks/{taskId}` - Получить информацию о задаче
- **DELETE** `/tasks/{taskId}` - Удалить задачу
- **GET** `/tasks/{taskId}/result` - Получить результат задачи

** Полную документацию, примеры запросов и ответов смотрите в Swagger UI интерфейсе.**

### Статусы задач

- `pending` - задача создана, но еще не начала выполняться
- `running` - задача выполняется
- `completed` - задача успешно завершена
- `failed` - задача завершилась с ошибкой

## OpenAPI спецификация

OpenAPI 3.0 спецификация находится в файле `v1/workmate.yaml` и используется для генерации серверного кода через [OpenAPI Generator](https://openapi-generator.tech/).

## Примеры использования

### Создание задачи
```bash
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Content-Type: application/json" \
  -d '{"name": "Process data"}'
```

### Получение списка задач
```bash
curl http://localhost:8080/api/v1/tasks
```

### Проверка статуса задачи
```bash
curl http://localhost:8080/api/v1/tasks/{taskId}
```

### Получение результата
```bash
curl http://localhost:8080/api/v1/tasks/{taskId}/result
```

### Удаление задачи
```bash
curl -X DELETE http://localhost:8080/api/v1/tasks/{taskId}
```

## Проверка HTTP/2

Для проверки, что сервер работает с HTTP/2:

```bash
# С сертификатами
curl -I https://localhost:8080/api/v1/tasks --http2 -k

# Проверить версию протокола
curl -I https://localhost:8080/api/v1/tasks --http2 -k -v 2>&1 | grep "HTTP/2"
```

## Разработка

### Изменение API

1. Отредактируйте `v1/workmate.yaml`
2. Перегенерируйте код: `make generate`
3. Скопируйте нужные файлы из `generated/go/` в `go/`
4. Реализуйте новую логику в `go/api_tasks_service.go`

### Запуск тестов
```bash
make test
```
