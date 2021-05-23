ifneq (,$(wildcard ./.env))
    include .env
    export
endif

dep:
	go mod tidy

lint:
	docker run --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:v1.37.1 golangci-lint run

fmt:
	go fmt ./src/...

test: fmt
	go test ./src/...

migration:
	echo -n "Key: "; read MIGRATION_NAME; \
	docker run --rm -v $(shell pwd)/src/migrations:/migrations migrate/migrate create -ext sql -dir /migrations -seq "$$MIGRATION_NAME"
	sudo chown $$USER:$$USER ./src/migrations/*

migrates:
	docker run --rm -v $(shell pwd)/src/migrations:/migrations --network host migrate/migrate \
        -path=/migrations -database mysql://$(DATABASE_USER):$(DATABASE_PASSWORD)@/$(DATABASE_NAME) up

build: fmt lint test
	docker-compose -f docker/docker-compose.yml build

up:
	docker-compose -f docker/docker-compose.yml up -d

down:
	docker-compose -f docker/docker-compose.yml down

logs:
	docker-compose -f docker/docker-compose.yml logs

api_tests: up
	docker run -v $(shell pwd)/api-tests:/app --network host postman/newman run --global-var url=localhost:${PORT} /app/go-hackaton.postman_collection.json
