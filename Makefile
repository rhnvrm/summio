build: build-frontend build-backend

build-backend:
	go build -o summio.bin main.go

build-frontend:
	cd frontend && npm run build

.PHONY: build build-backend build-frontend