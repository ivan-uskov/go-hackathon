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

hackaton_migration:
	echo -n "Key: "; read MIGRATION_NAME; \
	docker run --rm -v $(shell pwd)/src/cmd/hackatonservice/migrations:/migrations migrate/migrate create -ext sql -dir /migrations -seq "$$MIGRATION_NAME"
	sudo chown $$USER:$$USER ./src/cmd/hackatonservice/migrations/*

hackaton_migrates:
	docker run --rm -v $(shell pwd)/src/cmd/hackatonservice/migrations:/migrations --network host migrate/migrate \
        -path=/migrations -database "$(HACKATON_DATABASE_DRIVER)://$(HACKATON_DATABASE_USER):$(HACKATON_DATABASE_PASSWORD)@tcp(localhost:3370)/$(HACKATON_DATABASE_NAME)" up

scoring_migration:
	echo -n "Key: "; read MIGRATION_NAME; \
	docker run --rm -v $(shell pwd)/src/cmd/scoringservice/migrations:/migrations migrate/migrate create -ext sql -dir /migrations -seq "$$MIGRATION_NAME"
	sudo chown $$USER:$$USER ./src/cmd/scoringservice/migrations/*

scoring_migrates:
	docker run --rm -v $(shell pwd)/src/cmd/scoringservice/migrations:/migrations --network host migrate/migrate \
        -path=/migrations -database "$(SCORING_DATABASE_DRIVER)://$(SCORING_DATABASE_USER):$(SCORING_DATABASE_PASSWORD)@tcp(localhost:3371)/$(SCORING_DATABASE_NAME)" up

build: fmt lint test
	docker-compose -f docker/docker-compose.yml build

up:
	docker-compose -f docker/docker-compose.yml up -d

down:
	docker-compose -f docker/docker-compose.yml down

logs:
	docker-compose -f docker/docker-compose.yml logs

api_tests: up
	docker run -v $(shell pwd)/api-tests:/app --network host postman/newman run --global-var url=localhost:${HACKATON_PORT} /app/go-hackaton.postman_collection.json
	docker run -v $(shell pwd)/api-tests:/app --network host postman/newman run --global-var url=localhost:${SCORING_PORT} /app/scoringservice.postman_collection.json
