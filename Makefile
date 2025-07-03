#Makefile

.PHONY: all build run deps generate certs api-test unit-test

all: deps build

deps:
	go mod download
	go mod tidy

generate:
	openapi-generator-cli generate -g go-server -i ./api/contract/v1/workmate.yaml --git-repo-id workmate -o ./generated

build:
	go build -o workmate ./cmd/launcher/main.go

run:
	go run ./cmd/launcher/main.go

swagger:
	go run ./cmd/swagger/main.go

certs:
	chmod +x ./script/gen-certs.sh
	./script/gen-certs.sh

api-test:
	chmod +x ./script/test-api.sh
	./script/test-api.sh	

unit-test:
	go test -v ./... 
