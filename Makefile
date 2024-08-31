.PHONY: up down build swaggen

up:
	docker compose up -d

down:
	docker compose down

build:
	docker compose build

swaggen:
	docker run --rm -v .:/code ghcr.io/swaggo/swag:latest init -d internal/handler,./  -g /http.go

migrate:
	docker run -v internal/repository/postgres/migrations:/migrations --network host migrate/migrate \
	-path=/migrations/ -database postgres://db:5432/database up