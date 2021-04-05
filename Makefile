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

test:
	go test ./src/...

migration:
	echo -n "Key: "; read MIGRATION_NAME; \
	docker run --rm -v $(shell pwd)/src/migrations:/migrations migrate/migrate create -ext sql -dir /migrations -seq "$$MIGRATION_NAME"
	sudo chown $$USER:$$USER ./src/migrations/*

migrates:
	docker run --rm -v $(shell pwd)/src/migrations:/migrations --network host migrate/migrate \
        -path=/migrations -database mysql://$(DATABASE_USER):$(DATABASE_PASSWORD)@/$(DATABASE_NAME) up

build: fmt lint
	docker-compose -f docker/docker-compose.yml build

up:
	docker-compose -f docker/docker-compose.yml up

down:
	docker-compose -f docker/docker-compose.yml down