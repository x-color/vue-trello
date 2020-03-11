build-frontend:
	@cd web && npm run build:prod

build-backend: clean-db
	@if [ ! -d ./dist ]; then \
		mkdir dist; \
	fi
	@go build -o dist/server main.go

build: build-frontend build-backend
	@echo "Built SPA and API server"

clean-db:
	@if [ ! -d ./db ]; then \
		mkdir db; \
	fi
	@echo "" > db/sqlite.db

run-dev:
	@(cd web && npm run build:dev) &
	@DB_PATH=db/sqlite.db go run main.go

run: build
	@./dist/server