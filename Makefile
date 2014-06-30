NO_COLOR=\033[0m
TEXT_COLOR=\033[1m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m
DEPS = $(go list -f '{{range .TestImports}}{{.}} {{end}}' ./...)
TEST_DIR=tests/test-dirmode
FILE_MODE_TEST_DIR=tests/test-filemode
FILE_MODE_CONFIG=tests/test-filemode/test-flight-config.json
COMMON_OPTS=-race

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
	@echo "$(TEXT_COLOR) test-version: tests the version command $(NO_COLOR)"
	@echo "$(TEXT_COLOR) test-check: tests the check command $(NO_COLOR)"
	@echo "$(TEXT_COLOR) test-check-s: tests the check command using filemode in test dir $(NO_COLOR)"
	@echo "$(TEXT_COLOR) test-build: tests the build command using test dir $(NO_COLOR)"
	@echo "$(TEXT_COLOR) test-build-s: tests the build command using filemode in test dir $(NO_COLOR)"
	@echo "$(TEXT_COLOR) test-template: tests the template command using test dir $(NO_COLOR)"
	@echo "$(TEXT_COLOR) test-template-s: tests the template command using filemode in test dir $(NO_COLOR)"
	@echo "$(TEXT_COLOR) test-launch: tests the launch command using test dir $(NO_COLOR)"
	@echo "$(TEXT_COLOR) test-launch-f: tests the launch command using force in test dir $(NO_COLOR)"
	@echo "$(TEXT_COLOR) test-launch-f-s: tests the launch command using force and with single file mode in test dir$(NO_COLOR)"
	@echo "$(TEXT_COLOR) test-ground: tests the ground command using test dir $(NO_COLOR)"
	@echo "$(TEXT_COLOR) test-ground-s: tests the ground command using filemode in test dir $(NO_COLOR)"
	@echo "$(TEXT_COLOR) test-destroy: tests the destroy command using test dir $(NO_COLOR)"
	@echo "$(TEXT_COLOR) test-destroy-s: tests the destroy command using filemode in test dir $(NO_COLOR)"
	@echo "$(TEXT_COLOR) test-images: tests the images command using test dir $(NO_COLOR)"
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
	go test $(COMMON_OPTS) ./...

lint:
	@echo "$(OK_COLOR)==> Linting $(NO_COLOR)"
	golint .

test-version:
	@echo "$(OK_COLOR)==> Testing Version $(NO_COLOR)"
	go run $(COMMON_OPTS) flight.go version

test-check:
	@echo "$(OK_COLOR)==> Testing Check$(NO_COLOR)"
	go run $(COMMON_OPTS) flight.go check -s -d $(TEST_DIR)

test-check-s:
	@echo "$(OK_COLOR)==> Testing Check with FileMode set$(NO_COLOR)"
	go run $(COMMON_OPTS) flight.go check -s -d $(FILE_MODE_TEST_DIR)

test-build:
	@echo "$(OK_COLOR)==> Testing Build $(NO_COLOR)"
	go run $(COMMON_OPTS) flight.go build -d $(TEST_DIR)

test-build-s:
	@echo "$(OK_COLOR)==> Testing Build with FileMode set$(NO_COLOR)"
	go run $(COMMON_OPTS) flight.go -c $(FILE_MODE_CONFIG) -s -d $(FILE_MODE_TEST_DIR) build

test-template:
	@echo "$(OK_COLOR)==> Testing Template set$(NO_COLOR)"
	go run $(COMMON_OPTS) flight.go template -d $(TEST_DIR)

test-template-s:
	@echo "$(OK_COLOR)==> Testing Template with FileMode set$(NO_COLOR)"
	go run $(COMMON_OPTS) flight.go template -s -d $(FILE_MODE_TEST_DIR)

test-launch:
	@echo "$(OK_COLOR)==> Testing Launch $(NO_COLOR)"
	go run $(COMMON_OPTS) flight.go launch -d $(TEST_DIR)

test-launch-f:
	@echo "$(OK_COLOR)==> Testing Launch with Force set$(NO_COLOR)"
	go run $(COMMON_OPTS) flight.go launch -f -d $(TEST_DIR)

test-launch-f-s:
	@echo "$(OK_COLOR)==> Testing Launch with Force set$(NO_COLOR)"
	go run $(COMMON_OPTS) flight.go launch -f -s -d $(FILE_MODE_TEST_DIR)

test-ground:
	@echo "$(OK_COLOR)==> Testing Ground $(NO_COLOR)"
	go run $(COMMON_OPTS) flight.go ground -d $(TEST_DIR)

test-ground-s:
	@echo "$(OK_COLOR)==> Testing Ground with FileMode set$(NO_COLOR)"
	go run $(COMMON_OPTS) flight.go ground -s -d $(FILE_MODE_TEST_DIR)

test-destroy:
	@echo "$(OK_COLOR)==> Testing Destroy $(NO_COLOR)"
	go run $(COMMON_OPTS) flight.go destroy -d $(TEST_DIR)

test-destroy-s:
	@echo "$(OK_COLOR)==> Testing Destroy with FileMode set$(NO_COLOR)"
	go run $(COMMON_OPTS) flight.go destroy -s -d $(FILE_MODE_TEST_DIR)

test-images:
	@echo "$(OK_COLOR)==> Testing $(NO_COLOR)"
	go run $(COMMON_OPTS) flight.go images -d $(TEST_DIR)

all: format lint test
