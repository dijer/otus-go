BIN := "./bin/calendar"
DOCKER_IMG="calendar:develop"
PROTO_DIR="./api"
GENERATE_DIR="./internal/pb"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

.PHONY: build run build-img run-img version test lint generate dev rabbit-up rabbit-down

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar
	go build -v -o $(BIN)_scheduler -ldflags "$(LDFLAGS)" ./cmd/calendar_scheduler
	go build -v -o $(BIN)_sender -ldflags "$(LDFLAGS)" ./cmd/calendar_sender

run: build
	$(BIN) -config ./configs/config.toml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

test:
	go test -race ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.57.2

lint: install-lint-deps
	GOGC=50 golangci-lint run -v ./...

generate:
	rm -rf $(GENERATE_DIR)
	mkdir -p $(GENERATE_DIR)
	protoc -I=$(PROTO_DIR) \
		-I=$(PROTO_DIR)/googleapis \
		--go_out=$(GENERATE_DIR) \
		--go-grpc_out=$(GENERATE_DIR) \
		--grpc-gateway_out=$(GENERATE_DIR) \
		$(PROTO_DIR)/*.proto
	go generate -v ./...

dev:
	go run ./cmd/calendar/main.go

rabbit-up:
	docker-compose up -d

rabbit-down:
	docker-compose down