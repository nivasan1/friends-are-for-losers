TRACKER_ADDR ?=	wss://localhost:8080

get-linter:
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(golangci_version)


lint:
	@echo "--> Running linter"
	@golangci-lint run

lint-fix:
	@echo "--> Running linter + fixing shitty code"
	@golangci-lint run --fix

tidy:
	@echo "--> Running go mod tidy"
	@go mod tidy

generate-mocks:
	@echo "--> Generating mocks"
	@go generate ./...

workspace-init:
	@echo "--> Initializing workspace"
	@go work init
	@go work edit -use ./

test:
	@echo "--> Running tests"
	@go clean -testcache
	@go test -v -race ./...

run-tracker:
	@echo "--> Running tracker"
	@go run cmd/tracker/main.go --addr $(TRACKER_ADDR)
