package server

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/example/CRUD/controllers"
)

func initRoute() *gin.Engine {
	r := gin.Default()

	r.POST("/users", controllers.CreateUser)
	r.GET("/users/:username", controllers.GetUser)
	return r
}
