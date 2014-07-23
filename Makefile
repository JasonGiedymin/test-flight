NO_COLOR=\033[0m
TEXT_COLOR=\033[1m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m

DEPS = $(go list -f '{{range .TestImports}}{{.}} {{end}}' ./...)
TEST_DIR=tests/test-dirmode/example-playbook
FILE_MODE_TEST_DIR=tests/test-filemode/example-playbook
FILE_MODE_CONFIG=tests/test-filemode/example-playbook/test-flight-config.json
COMMON_OPTS=-race
COMMON_FLIGHT_OPTS=-vvvv
PACKAGE=github.com/JasonGiedymin/test-flight
PATH_SRC=$(GOPATH)/src/github.com/JasonGiedymin/test-flight
PATH_PKG=$(GOPATH)/pkg/*/github.com/JasonGiedymin/test-flight

help:
	@echo "$(OK_COLOR)-----------------------Commands:----------------------$(NO_COLOR)"
	@echo "$(TEXT_COLOR) help-cmd:   this help listing $(NO_COLOR)"
	@echo "$(TEXT_COLOR) link:       symlinks this repo to gopath $(NO_COLOR)"
	@echo "$(TEXT_COLOR) install:    installs test-flight (via go) $(NO_COLOR)"
	@echo "$(TEXT_COLOR) uninstall:  uninstalls test-flight (via go) $(NO_COLOR)"
	@echo "$(TEXT_COLOR) reinstall:  uninstalls and then installs test-flight (via go) $(NO_COLOR)"
	@echo "$(TEXT_COLOR) dev-setup   set up developer environment $(NO_COLOR)"
	@echo "$(TEXT_COLOR) deps:       install dependencies $(NO_COLOR)"
	@echo "$(TEXT_COLOR) updatedeps: update dependencies $(NO_COLOR)"
	@echo "$(TEXT_COLOR) format:     formats the code $(NO_COLOR)"
	@echo "$(TEXT_COLOR) test:       tests code $(NO_COLOR)"
	@echo "$(TEXT_COLOR) lint:       lints code $(NO_COLOR)"
	@echo ""
	@echo "$(TEXT_COLOR) == Docker helpers == $(NO_COLOR)"
	@echo "$(TEXT_COLOR) docker-clean: removes all stopped containers and untagged images $(NO_COLOR)"
	@echo ""
	@echo "$(TEXT_COLOR) == Commands requiring docker endpoint == $(NO_COLOR)"
	@echo "$(TEXT_COLOR) test-version: tests the version command $(NO_COLOR)"
	@echo "$(TEXT_COLOR) test-check: tests the check command $(NO_COLOR)"
	@echo "$(TEXT_COLOR) test-check-s: tests the check command using filemode in test dir $(NO_COLOR)"
	@echo "$(TEXT_COLOR) test-build: tests the build command using test dir $(NO_COLOR)"
	@echo "$(TEXT_COLOR) test-build-s: tests the build command using filemode in test dir $(NO_COLOR)"
	@echo "$(TEXT_COLOR) test-template: tests the template command using test dir $(NO_COLOR)"
	@echo "$(TEXT_COLOR) test-template-f-s: tests the template command using filemode in test dir $(NO_COLOR)"
	@echo "$(TEXT_COLOR) test-launch: tests the launch command using test dir $(NO_COLOR)"
	@echo "$(TEXT_COLOR) test-launch-f: tests the launch command using force in test dir $(NO_COLOR)"
	@echo "$(TEXT_COLOR) test-launch-f-s: tests the launch command using force and with single file mode in test dir$(NO_COLOR)"
	@echo "$(TEXT_COLOR) test-ground: tests the ground command using test dir $(NO_COLOR)"
	@echo "$(TEXT_COLOR) test-ground-f-s: tests the ground command using filemode in test dir $(NO_COLOR)"
	@echo "$(TEXT_COLOR) test-destroy: tests the destroy command using test dir $(NO_COLOR)"
	@echo "$(TEXT_COLOR) test-destroy-s: tests the destroy command using filemode in test dir $(NO_COLOR)"
	@echo "$(TEXT_COLOR) test-images: tests the images command using test dir $(NO_COLOR)"
	@echo "$(TEXT_COLOR) test-images-s: tests the images command using filemode in test dir $(NO_COLOR)"
	@echo "$(OK_COLOR)------------------------------------------------------$(NO_COLOR)"

help-cmd:
	@echo "$(OK_COLOR)==> Testing Help $(PATH_SRC) $(NO_COLOR)"
	go run $(COMMON_OPTS) test-flight.go $(COMMON_FLIGHT_OPTS) --help

link:
	@echo "$(OK_COLOR)==> Symlinking project to $(PATH_SRC) $(NO_COLOR)"
	@ln -vFfsn $(shell pwd) $(PATH_SRC)

	@echo "$(OK_COLOR)==> Creating local test dir: $(shell pwd)/.test-flight $(NO_COLOR)"
	@mkdir -p $(shell pwd)/.test-flight
	@echo "$(OK_COLOR)==> Symlinking test-flight-config.json into .test-flight/ $(NO_COLOR)"
	@ln -vFfsn $(shell pwd)/test-flight-config.json $(shell pwd)/.test-flight/
	@echo "$(OK_COLOR)==> Symlinking templates into .test-flight $(NO_COLOR)"
	@ln -vFfsn $(shell pwd)/templates $(shell pwd)/.test-flight/

	@echo "$(OK_COLOR)==> Creating home test dir: $(HOME)/.test-flight $(NO_COLOR)"
	@mkdir -p $(HOME)/.test-flight
	@echo "$(OK_COLOR)==> Symlinking test-flight-config.json into $(HOME)/.test-flight/ $(NO_COLOR)"
	@ln -vFfsn $(shell pwd)/test-flight-config.json $(HOME)/.test-flight/
	@echo "$(OK_COLOR)==> Symlinking templates into $(HOME)/.test-flight $(NO_COLOR)"
	@ln -vFfsn $(shell pwd)/templates $(HOME)/.test-flight/

install-pkg:
	@echo "$(OK_COLOR)==> Installing Test-Flight $(PATH_SRC) $(NO_COLOR)"
	@go install $(PACKAGE)

install: install-pkg

uninstall:
	@echo "$(OK_COLOR)==> Uninstalling Test-Flight $(PATH_PKG) $(NO_COLOR)"
	@if [ -d $(PATH_PKG) ]; then rm -R $(PATH_PKG); fi;

reinstall: uninstall install

dev-setup: deps link

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
	go test -v $(COMMON_OPTS) ./...

lint:
	@echo "$(OK_COLOR)==> Linting $(NO_COLOR)"
	golint .

docker-clean:
	@echo "$(OK_COLOR)==> Removes all stopped docker containers $(NO_COLOR)"
	docker ps -a | grep 'Exited' |  awk '{print $$1}' | xargs docker rm
	@echo "$(OK_COLOR)==> Removes all untagged docker images $(NO_COLOR)"
	docker images | grep '<none>' |  awk '{print $$3}' | xargs docker rmi

test-version:
	@echo "$(OK_COLOR)==> Testing Version $(NO_COLOR)"
	go run $(COMMON_OPTS) test-flight.go $(COMMON_FLIGHT_OPTS) version

test-check:
	@echo "$(OK_COLOR)==> Testing Check $(NO_COLOR)"
	go run $(COMMON_OPTS) test-flight.go $(COMMON_FLIGHT_OPTS) -d $(TEST_DIR) check

test-check-s:
	@echo "$(OK_COLOR)==> Testing Check with FileMode set $(NO_COLOR)"
	go run $(COMMON_OPTS) test-flight.go $(COMMON_FLIGHT_OPTS) -c $(FILE_MODE_CONFIG) -s -d $(FILE_MODE_TEST_DIR) check

test-build:
	@echo "$(OK_COLOR)==> Testing Build $(NO_COLOR)"
	go run $(COMMON_OPTS) test-flight.go $(COMMON_FLIGHT_OPTS) build -d $(TEST_DIR)

test-build-s:
	@echo "$(OK_COLOR)==> Testing Build with FileMode set $(NO_COLOR)"
	go run $(COMMON_OPTS) test-flight.go $(COMMON_FLIGHT_OPTS) -c $(FILE_MODE_CONFIG) -s -d $(FILE_MODE_TEST_DIR) build

test-template:
	@echo "$(OK_COLOR)==> Testing Template set $(NO_COLOR)"
	go run $(COMMON_OPTS) test-flight.go $(COMMON_FLIGHT_OPTS) -d $(TEST_DIR) template

test-template-f-s:
	@echo "$(OK_COLOR)==> Testing Template with FileMode set $(NO_COLOR)"
	go run $(COMMON_OPTS) test-flight.go $(COMMON_FLIGHT_OPTS) -f -s -d $(FILE_MODE_TEST_DIR) template

test-launch:
	@echo "$(OK_COLOR)==> Testing Launch $(NO_COLOR)"
	go run $(COMMON_OPTS) test-flight.go $(COMMON_FLIGHT_OPTS) -d $(TEST_DIR) launch

test-launch-f:
	@echo "$(OK_COLOR)==> Testing Launch with Force set $(NO_COLOR)"
	go run $(COMMON_OPTS) test-flight.go $(COMMON_FLIGHT_OPTS) -f -d $(TEST_DIR) launch

test-launch-f-s:
	@echo "$(OK_COLOR)==> Testing Launch with Force set and using FileMode $(NO_COLOR)"
	go run $(COMMON_OPTS) test-flight.go $(COMMON_FLIGHT_OPTS) -c $(FILE_MODE_CONFIG) -f -s -d $(FILE_MODE_TEST_DIR) launch

test-ground:
	@echo "$(OK_COLOR)==> Testing Ground $(NO_COLOR)"
	go run $(COMMON_OPTS) test-flight.go $(COMMON_FLIGHT_OPTS) -d $(TEST_DIR) ground

test-ground-f-s:
	@echo "$(OK_COLOR)==> Testing Ground with FileMode set $(NO_COLOR)"
	go run $(COMMON_OPTS) test-flight.go $(COMMON_FLIGHT_OPTS) -c $(FILE_MODE_CONFIG) -s -d $(FILE_MODE_TEST_DIR) ground

test-destroy:
	@echo "$(OK_COLOR)==> Testing Destroy $(NO_COLOR)"
	go run $(COMMON_OPTS) test-flight.go $(COMMON_FLIGHT_OPTS) destroy -d $(TEST_DIR)

test-destroy-f-s:
	@echo "$(OK_COLOR)==> Testing Destroy with FileMode set $(NO_COLOR)"
	go run $(COMMON_OPTS) test-flight.go $(COMMON_FLIGHT_OPTS) -c $(FILE_MODE_CONFIG) -s -d $(FILE_MODE_TEST_DIR) destroy

test-images:
	@echo "$(OK_COLOR)==> Testing Images $(NO_COLOR)"
	go run $(COMMON_OPTS) test-flight.go $(COMMON_FLIGHT_OPTS) -d $(TEST_DIR) images

test-images-s:
	@echo "$(OK_COLOR)==> Testing Images with FileMode set $(NO_COLOR)"
	go run $(COMMON_OPTS) test-flight.go $(COMMON_FLIGHT_OPTS) -c $(FILE_MODE_CONFIG) -s -d $(FILE_MODE_TEST_DIR) images

all: format lint test
