package main

import (
	"database/sql"
	"fmt"

	db "go.opentelemetry.io/otel/example/CRUD/db/sqlc"
	trace "go.opentelemetry.io/otel/example/CRUD/trace"
)

func main() {
	conn, _ := sql.Open("config.DBDriver", "config.DBSource")
	database := trace.NewDBTXTrace(conn)
	store := db.NewStore(database)
	fmt.Println(store)
	// store.CreateUser()
}
