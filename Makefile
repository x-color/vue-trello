build-frontend:
	@cd web && npm run build:prod

build-backend:
	@if [ ! -d ./dist ]; then \
		mkdir dist; \
	fi
	@go build -o dist/server main.go

build: build-frontend build-backend
	@echo "Built SPA and API server"

clean-db:
	@echo "" > db/sqlite.db
	@echo "Reset DB"

run-dev:
	@(cd web && npm run build:dev) &
	@DB_PATH=db/sqlite.db go run main.go

run: build
	@./dist/server