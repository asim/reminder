.PHONY: dev dev-backend dev-frontend

# Development mode
dev:
	@trap 'kill $$(jobs -p)' EXIT; \
	make dev-backend & \
	make dev-frontend & \
	wait

dev-backend:
	go run . --serve

dev-frontend:
	cd web && pnpm run dev

# Building
build:
	go build -o reminder main.go

run:
	./reminder

