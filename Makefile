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

.PHONY migrations_up:
migrations_up:
	migrate -path migrations/ -database postgres://root:dfcz123@localhost:5432/Users_Info?sslmode=disable up

.PHONY migrations_down:
migrations_down:
	migrate -path migrations/ -database postgres://root:dfcz123@localhost:5432/Users_Info?sslmode=disable down