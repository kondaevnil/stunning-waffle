.PHONY: build run test docker-build docker-up docker-down clean

build:
	go build -o bin/main cmd/main.go

run:
	go run cmd/main.go

test:
	go test -v ./...

docker-build:
	docker-compose build

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

docker-restart: docker-down docker-build docker-up

clean:
	rm -rf bin/
	docker-compose down -v
	docker system prune -f

db-connect:
	docker exec -it ecom_postgres psql -U postgres -d ecom

init:
	go mod init vk/ecom
	go mod tidy

deps:
	go mod download
	go mod tidy