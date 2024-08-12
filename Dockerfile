FROM golang:latest AS build

WORKDIR /sso

COPY . .

RUN go mod download && go mod verify

RUN go build -o ./sso ./cmd/sso/main.go

FROM scratch AS prod

WORKDIR /sso

COPY --from=build /sso/sso /sso/

COPY /app_configs/dev.yaml .

ENTRYPOINT [ "/sso/sso", "c", "/sso/dev.yaml"]