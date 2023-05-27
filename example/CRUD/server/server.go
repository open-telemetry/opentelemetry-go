package server

import (
	"github.com/gin-gonic/gin"
	db "go.opentelemetry.io/otel/example/CRUD/db/sqlc"
)

type Server struct {
	routes *gin.Engine
	store  db.Store
}

func NewServer(store db.Store) *Server {
	server := &Server{}
	server.store = store
	server.setupRoutes()
	return server
}

func (server *Server) Start() {
	server.routes.Run()
}
