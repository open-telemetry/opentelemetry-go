package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "go.opentelemetry.io/otel/example/CRUD/db/sqlc"
)

func (server *Server) CreateUser(ctx *gin.Context) {
	var user db.CreateUserParams
	err := ctx.Bind(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	u, err := server.store.CreateUser(ctx, user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Can not create an user"})
		return
	}
	ctx.JSON(http.StatusOK, u)

}
