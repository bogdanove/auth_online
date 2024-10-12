LOCAL_BIN:=$(CURDIR)/bin

install-golangci-lint:
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.60.3

lint:
	golangci-lint run ./... --config .golangci.pipeline.yaml

build:
	GOOS=linux GOARCH=amd64 go build -o service_linux cmd/main.go

docker-build-and-push:
	docker buildx build --no-cache --platform linux/amd64 -t cr.selcloud.ru/course/auth:v0.0.1 .
	docker login -u token -p CRgAAAAAdtXKiWnNLiGAWxJ8LGJLQXERjYGJ6mkf cr.selcloud.ru/course
	docker push cr.selcloud.ru/course/auth:v0.0.1