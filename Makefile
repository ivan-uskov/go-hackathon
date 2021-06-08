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

hackathon_migration:
	echo -n "Key: "; read MIGRATION_NAME; \
	docker run --rm -v $(shell pwd)/src/hackathonservice/migrations:/migrations migrate/migrate create -ext sql -dir /migrations -seq "$$MIGRATION_NAME"
	sudo chown $$USER:$$USER ./src/hackathonservice/migrations/*

hackathon_migrates:
	docker run --rm -v $(shell pwd)/src/hackathonservice/migrations:/migrations --network host migrate/migrate \
        -path=/migrations -database "$(HACKATHON_DATABASE_DRIVER)://$(HACKATHON_DATABASE_USER):$(HACKATHON_DATABASE_PASSWORD)@tcp(localhost:3370)/$(HACKATHON_DATABASE_NAME)" up

scoring_migration:
	echo -n "Key: "; read MIGRATION_NAME; \
	docker run --rm -v $(shell pwd)/src/scoringservice/migrations:/migrations migrate/migrate create -ext sql -dir /migrations -seq "$$MIGRATION_NAME"
	sudo chown $$USER:$$USER ./src/scoringservice/migrations/*

scoring_migrates:
	docker run --rm -v $(shell pwd)/src/scoringservice/migrations:/migrations --network host migrate/migrate \
        -path=/migrations -database "$(SCORING_DATABASE_DRIVER)://$(SCORING_DATABASE_USER):$(SCORING_DATABASE_PASSWORD)@tcp(localhost:3371)/$(SCORING_DATABASE_NAME)" up

build: fmt lint test proto
	docker-compose -f docker/docker-compose.yml build --force-rm --parallel

up:
	docker-compose -f docker/docker-compose.yml up -d

down:
	docker-compose -f docker/docker-compose.yml down

logs:
	docker-compose -f docker/docker-compose.yml logs

api_tests: up
	docker run --rm -v $(shell pwd)/api-tests:/app --network host postman/newman run --global-var url=localhost:${HACKATHON_HTTP_PORT} /app/hackathonservice.postman_collection.json
	docker run --rm -v $(shell pwd)/api-tests:/app --network host postman/newman run --global-var url=localhost:${SCORING_HTTP_PORT} /app/scoringservice.postman_collection.json

proto:
	docker run --rm -v "$(shell pwd)/api/hackathonservice:/app" ivanuskov/go-protobuf-builder *.proto
	docker run --rm -v "$(shell pwd)/api/scoringservice:/app" ivanuskov/go-protobuf-builder *.proto
	sudo chown $$USER:$$USER -R ./api
