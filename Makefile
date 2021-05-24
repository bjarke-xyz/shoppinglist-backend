.PHONY: clean swag build

include .env

APP_NAME = shoppinglist-backend
BUILD_DIR = $(PWD)/build
MIGRATIONS_FOLDER = $(PWD)/platform/migrations
DOCKER_MIGRATION = docker run --rm --user "$$(id -u):$$(id -g)" -v $(MIGRATIONS_FOLDER):/migrations --network host migrate/migrate

clean:
	rm -rf $(BUILD_DIR)

build: clean swag
	go build -o $(BUILD_DIR)/$(APP_NAME) main.go

swag:
	swag init

migrate.up:
	$(DOCKER_MIGRATION) -path /migrations -database ${DB_SERVER_URL} up

migrate.down:
	$(DOCKER_MIGRATION) -path /migrations -database ${DB_SERVER_URL} down -all

# Usage: make ... version=<version>
migrate.force:
	$(DOCKER_MIGRATION) -path /migrations -database ${DB_SERVER_URL} force $(version)

# Usage: make ... force=<true | false>
migrate.drop:
	$(DOCKER_MIGRATION) -path /migrations -database ${DB_SERVER_URL} drop -f=$(force)

# Usage: make ... name=<name of migration>
migrate.create:
	$(DOCKER_MIGRATION) create -dir /migrations -ext sql $(name)
