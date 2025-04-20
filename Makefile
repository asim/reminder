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
	cd web && pnpm run build
	find app/html -mindepth 1 -not -name '.gitignore' -delete
	cp -r web/build/client/* app/html/
	go build -o reminder .
