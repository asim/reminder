.PHONY: dev dev-backend dev-frontend build setup

setup:
	cd web && pnpm install

# Development
dev:
	@trap 'kill $$(jobs -p)' EXIT; \
	make dev-backend & \
	make dev-frontend & \
	wait

dev-backend:
	go run . --serve --web

dev-frontend:
	cd web && pnpm run dev

# Building
build:
	cd web && pnpm run build
	find app/dist -mindepth 1 -not -name '.gitignore' -delete
	cp -r web/build/client/* app/dist/
	go build -o reminder .
