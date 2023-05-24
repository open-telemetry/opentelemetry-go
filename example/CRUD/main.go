package main

import (
	"database/sql"

	db "go.opentelemetry.io/otel/example/CRUD/db/sqlc"
	"go.opentelemetry.io/otel/example/CRUD/server"
	trace "go.opentelemetry.io/otel/example/CRUD/trace"
)

func main() {
	
	conn, _ := sql.Open("config.DBDriver", "config.DBSource")
	database := trace.NewDBTXTrace(conn)
	store := db.NewStore(database)

	server := server.NewServer(&store)
	server.Start()
}
