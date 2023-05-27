package server

import (
	"github.com/gin-gonic/gin"
	db "go.opentelemetry.io/otel/example/CRUD/db/sqlc"
)

func (server *Server) CreateUser(ctx *gin.Context) {
	server.store.CreateUser(ctx, db.CreateUserParams{})
}
