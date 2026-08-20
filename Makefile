.PHONY: build up down logs migrate seed

build:
	docker compose build --no-cache backend frontend

up:
	docker compose up -d --force-recreate

down:
	docker compose down

logs:
	docker compose logs -f backend frontend

migrate:
	docker compose exec backend go run cmd/migrate/main.go

seed:
	docker compose exec backend go run cmd/seed/main.go

dev-backend:
	cd backend && go run cmd/server/main.go

dev-frontend:
	cd frontend && npm run dev
