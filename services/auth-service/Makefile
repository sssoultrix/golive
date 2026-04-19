# In WSL, Linux `docker` often uses /var/run/docker.sock (empty if Desktop integration is off).
# Prefer `docker.exe` so commands hit Docker Desktop on Windows.
WSL := $(shell grep -qi microsoft /proc/version 2>/dev/null && echo 1)

DOCKER :=
ifeq ($(WSL),1)
  ifneq ($(shell command -v docker.exe 2>/dev/null),)
    DOCKER := docker.exe
  else
    _win_docker := /mnt/c/Program Files/Docker/Docker/resources/bin/docker.exe
    ifneq ($(shell test -x '$(_win_docker)' && echo 1),)
      DOCKER := $(_win_docker)
    endif
  endif
endif
ifeq ($(DOCKER),)
  _which := $(shell command -v docker 2>/dev/null)
  ifneq ($(_which),)
    DOCKER := $(_which)
  endif
endif
ifeq ($(DOCKER),)
  _which := $(shell command -v docker.exe 2>/dev/null)
  ifneq ($(_which),)
    DOCKER := $(_which)
  endif
endif
ifeq ($(DOCKER),)
  _win_docker := /mnt/c/Program Files/Docker/Docker/resources/bin/docker.exe
  ifneq ($(shell test -x '$(_win_docker)' && echo 1),)
    DOCKER := $(_win_docker)
  endif
endif
ifeq ($(DOCKER),)
  DOCKER := docker
endif

# Override if needed: make migrate-up DOCKER=/path/to/docker
COMPOSE := "$(DOCKER)" compose --profile tools
MIGRATE_DB ?= postgres://auth:auth@postgres:5432/auth?sslmode=disable
MIGRATE_FLAGS := -path=/migrations -database=$(MIGRATE_DB)
MIGRATE := $(COMPOSE) run --rm migrate $(MIGRATE_FLAGS)

.PHONY: up down migrate-up migrate-down migrate-version migrate-force migrate-create cleanup-refresh-tokens

up:
	$(COMPOSE) up -d postgres redis auth-service

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

cleanup-refresh-tokens:
	$(COMPOSE) run --rm cleanup
