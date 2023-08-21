
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
