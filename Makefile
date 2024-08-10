.PHONY build:
build:
	go build -o ./cmd/sso/sso ./cmd/sso/main.go

.PHONY run:
run:
	./cmd/sso/sso

.PHONY all:
all: build run