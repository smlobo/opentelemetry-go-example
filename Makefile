
TARGET = opentelemetry-go-example
REGISTRY = localhost:32000/

build:
	go build -o bin/${TARGET} cmd/${TARGET}/main.go

module:
	rm -f go.mod go.sum
	go mod init ${TARGET}
	go mod tidy

.PHONY: docker
docker:
	docker build -t ${REGISTRY}otlp-go-example-backend:$(tag) -f docker/Dockerfile.backend .
	docker tag ${REGISTRY}otlp-go-example-backend:$(tag) ${REGISTRY}otlp-go-example-backend:latest
	docker build -t ${REGISTRY}otlp-go-example-frontend:$(tag) -f docker/Dockerfile.frontend .
	docker tag ${REGISTRY}otlp-go-example-frontend:$(tag) ${REGISTRY}otlp-go-example-frontend:latest

clean:
	rm -rf bin
	rm -f go.mod go.sum
	rm *traces.txt
