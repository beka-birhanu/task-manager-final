build:
	@go build -o bin/task-manager ./main.go

test:
	go test -v ./...

run: build
	@./bin/task-manager
