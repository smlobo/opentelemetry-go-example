
TARGET = opentelemetry-go-example

build:
	go build -o bin/${TARGET} cmd/${TARGET}/main.go

module:
	rm -f go.mod go.sum
	go mod init ${TARGET}
	go mod tidy

.PHONY: docker
docker:
	docker build -t otlp-go-example-backend:$(tag) -f docker/Dockerfile.backend .
	docker tag otlp-go-example-backend:$(tag) otlp-go-example-backend:latest
	docker build -t otlp-go-example-frontend:$(tag) -f docker/Dockerfile.frontend .
	docker tag otlp-go-example-frontend:$(tag) otlp-go-example-frontend:latest

clean:
	rm -rf bin
	rm -f go.mod go.sum
	rm *traces.txt
