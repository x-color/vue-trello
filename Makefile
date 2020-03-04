run:
	(cd web; npm run build:dev) &
	echo "" > db/sqlite.db
	DB_PATH=db/sqlite.db go run main.go