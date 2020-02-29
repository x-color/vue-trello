run:
	echo "" > db/sqlite.db
	go build main.go
	DB_PATH=db/sqlite.db ./main