VERSION=$(shell git describe --tags --always 2>/dev/null || echo "dev")

.PHONY: run
# run locally
run:
	go run ./cmd/app/ -conf configs/config.local.yaml

.PHONY: start
# start docker container locally
start:
	docker compose build && docker compose up -d

.PHONY: stop
# stop docker container locally
stop:
	docker compose down

.PHONY: build
# build executable file
build:
	mkdir -p bin/ && go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/ ./...

.PHONY: test
# run tests
test:
	go test -v -count=1 ./...

.PHONY: lint
# run linter
lint:
	golangci-lint run ./...

.PHONY: race
# run tests with race
race:
	go test -v -race -count=10 ./...

.PHONY: tidy
# go mod tidy
tidy:
	go mod tidy

# show help
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help
