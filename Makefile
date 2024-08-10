.PHONY build:
build:
	go build -o ./cmd/sso/bin/sso_app ./cmd/sso/main.go

.PHONY run:
run:
	./cmd/sso/bin/sso_app -c ./app_configs/dev.yaml

.PHONY all:
all: build run

.PHONY clean:
clean:
	rm -rf ./cmd/sso/bin/;