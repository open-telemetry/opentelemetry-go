package server

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/example/CRUD/trace"
)

func (server *Server) setupRoutes() {
	r := gin.Default()
	r.Use(otelgin.Middleware(trace.ServiceName))

	r.POST("/users", server.CreateUser)
	r.GET("/users/:username", server.GetUser)

	server.routes = r
}
