include .env

MIGRATE_DB_URL=postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

migrate\:init:
	PGPASSWORD=$(DB_PASSWORD) psql -U$(DB_USER) -d postgres -h $(DB_HOST) -p $(DB_PORT) -c "create database $(DB_NAME);"

migrate\:drop:
	PGPASSWORD=$(DB_PASSWORD) psql -U$(DB_USER) -d postgres -h $(DB_HOST) -p $(DB_PORT) -c "drop database if exists $(DB_NAME) with (force);"

migrate\:up:
	migrate -database "$(MIGRATE_DB_URL)" -path migrations up $(version)

migrate\:down:
	migrate -database "$(MIGRATE_DB_URL)" -path migrations down $(version)


migrate\:reset: 
	$(MAKE) migrate:drop 
	$(MAKE) migrate:init 
	$(MAKE) migrate:up 

# make migrate:down version=1
# make migrate:down version=1
# migrate create -ext sql -dir migrations -seq create_users