.PHONY: build test clean k6

# The name of the binary to be built
BINARY_NAME=basic-auth-to-oauth2-transformer

# This will take the Go files and compile them into a binary
build:
	go build -o out/$(BINARY_NAME) cmd/main.go

# This will run the tests in the project
test:
	go test ./...

cover:
	go test ./... -coverprofile=tmp/c.out
	go tool cover -html="tmp/c.out"

# This will remove the binary
clean:
	go clean
	rm out/$(BINARY_NAME)

# This will start the Docker Compose services
docker-up:
	docker-compose -f docker/docker-compose.yaml up -d

docker-down:
	docker-compose -f docker/docker-compose.yaml down

docker-build:
	docker build -f docker/Dockerfile -t basic-auth-to-oauth2-transformer .

k6:
	k6 run k6/test.js
