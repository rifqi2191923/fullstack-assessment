backend:
	cd backend && go run ./cmd/server

infra:
	docker compose up -d mysql redis

infra-down:
	docker compose down
