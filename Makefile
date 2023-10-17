build: build-frontend build-backend

build-backend:
	go build -o summio.bin main.go

build-frontend:
	cd frontend && npm run build

release:
	@echo "Last tag: ${git describe --tags --abbrev=0}"
	@read -p "Enter Version Name:" version
	@read -p "Enter Release message:" release_message 
	git tag -a $$version -m "$$release_message"
	git push origin $$version
	goreleaser release

.PHONY: build build-backend build-frontend
.ONESHELL: release