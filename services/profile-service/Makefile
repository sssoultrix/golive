DOCKER := cmd.exe /c docker

# Override if needed: make migrate-up DOCKER=/path/to/docker
COMPOSE := $(DOCKER) compose --profile tools
MIGRATE_DB ?= postgres://profile:profile@postgres:5432/profile?sslmode=disable
MIGRATE_FLAGS := -path=/migrations -database=$(MIGRATE_DB)
MIGRATE := $(COMPOSE) run --rm migrate $(MIGRATE_FLAGS)

.PHONY: up down migrate-up migrate-down migrate-version migrate-force migrate-create

up:
	$(COMPOSE) up -d postgres profile-service

down:
	$(COMPOSE) down

migrate-up:
	$(MIGRATE) up

migrate-down:
	$(MIGRATE) down 1

migrate-version:
	$(MIGRATE) version

# Usage: make migrate-force VERSION=0
migrate-force:
	$(MIGRATE) force $(VERSION)

# Usage: make migrate-create NAME=init_schema
migrate-create:
	@test -n "$(NAME)" || (echo "usage: make migrate-create NAME=your_migration_name" && exit 1)
	$(COMPOSE) run --rm migrate create -ext sql -dir /migrations -seq $(NAME)
