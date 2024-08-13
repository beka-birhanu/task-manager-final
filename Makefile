build:
	@go build -o bin/task-manager ./main.go

test:
	@test_packages=$$(./scripts/filter_test_packages.sh); \
	go test -v -cover $$test_packages

coverage:
	@test_packages=$$(./scripts/filter_test_packages.sh); \
	go test -v -coverprofile=coverage.out $$test_packages; \
	go tool cover -html=coverage.out

run: build
	@./bin/task-manager
