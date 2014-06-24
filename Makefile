NO_COLOR=\033[0m
TEXT_COLOR=\033[1m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m
DEPS = $(go list -f '{{range .TestImports}}{{.}} {{end}}' ./...)
TEST_DIR=test

help:
	@echo "$(OK_COLOR)-----------------------Commands:----------------------$(NO_COLOR)"
	@echo "$(TEXT_COLOR) help:       this help listing $(NO_COLOR)"
	@echo "$(TEXT_COLOR) deps:       install dependencies $(NO_COLOR)"
	@echo "$(TEXT_COLOR) updatedeps: update dependencies $(NO_COLOR)"
	@echo "$(TEXT_COLOR) format:     formats the code $(NO_COLOR)"
	@echo "$(TEXT_COLOR) test:       tests code $(NO_COLOR)"
	@echo "$(TEXT_COLOR) lint:       lints code $(NO_COLOR)"
	@echo ""
	@echo "$(TEXT_COLOR) Commands requiring docker endpoint: $(NO_COLOR)"
	@echo ""
	@echo "$(TEXT_COLOR) run-launch: executes launch command using test dir $(NO_COLOR)"
	@echo "$(TEXT_COLOR) run-ground: executes ground command using test dir $(NO_COLOR)"
	@echo "$(TEXT_COLOR) run-destroy: executes destroy command using test dir $(NO_COLOR)"
	@echo "$(TEXT_COLOR) run-images: executes images command using test dir $(NO_COLOR)"
	@echo "$(OK_COLOR)------------------------------------------------------$(NO_COLOR)"

deps:
	@echo "$(OK_COLOR)==> Installing dependencies $(NO_COLOR)"
	@go get -d -v ./...
	@echo $(DEPS) | xargs -n1 go get -d

updatedeps:
	@echo "$(OK_COLOR)==> Updating all dependencies $(NO_COLOR)"
	@go get -d -v -u ./...
	@echo $(DEPS) | xargs -n1 go get -d -v -u

format:
	@echo "$(OK_COLOR)==> Formatting $(NO_COLOR)"
	go fmt ./...

test: deps
	@echo "$(OK_COLOR)==> Testing $(NO_COLOR)"
	go test ./...

lint:
	@echo "$(OK_COLOR)==> Linting $(NO_COLOR)"
	golint .

run-launch:
	@echo "$(OK_COLOR)==> Testing Launch $(NO_COLOR)"
	go run flight.go launch -d $(TEST_DIR)

run-ground:
	@echo "$(OK_COLOR)==> Testing Ground $(NO_COLOR)"
	go run flight.go ground -d $(TEST_DIR)

run-destroy:
	@echo "$(OK_COLOR)==> Testing Destroy $(NO_COLOR)"
	go run flight.go destroy -d $(TEST_DIR)

run-images:
	@echo "$(OK_COLOR)==> Testing $(NO_COLOR)"
	go run flight.go images -d $(TEST_DIR)

all: format lint test
