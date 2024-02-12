# Define the docker-compose command
DC := cd development && docker-compose

# Services names
MYSQL_SERVICE := mysql
REDIS_SERVICE := redis
MIGRATION_SERVICE := migration

# Paths
MIGRATE_UP_PATH := migration/migrate_up.sql
MIGRATE_DOWN_PATH := migration/migrate_down.sql

# Default to help
default: help

# bundle
bundle:
	file_bundle -v

# Start services
compose-up:
	$(DC) up -d

# Stop services
compose-down:
	$(DC) down

compose-re: compose-down compose-up

# Display logs
compose-logs:
	$(DC) logs

# Run database migrations
db-migrate-up:
	$(DC) exec $(MYSQL_SERVICE) sh -c 'mysql -u root -pexample memorianexus < $(MIGRATE_UP_PATH)'

# Rollback database migrations
db-migrate-down:
	$(DC) exec $(MYSQL_SERVICE) sh -c 'mysql -u root -pexample memorianexus < $(MIGRATE_DOWN_PATH)'

# Clean up environment including volumes
compose-reset:
	$(DC) down --volumes

# Show help
help:
	@echo "Available commands:"
	@echo "  make compose-up : Start all services with docker-compose"
	@echo "  make compose-down : Stop all services and remove containers"
	@echo "  make compose-logs : Fetch logs for all services"
	@echo "  make db-migrate-up : Perform database migrations"
	@echo "  make db-migrate-down : Rollback database migrations"
	@echo "  make compose-reset : Stop all services and remove data"

dev:
	MEMORIA_NUXUS_ENV="dev" go run ./cmd

.PHONY: default compose-up compose-down compose-re compose-logs db-migrate-up db-migrate-down compose-reset help bundle dev