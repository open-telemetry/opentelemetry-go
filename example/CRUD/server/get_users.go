package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (server *Server) GetUser(ctx *gin.Context) {
	username := ctx.Params.ByName("username")
	user, err := server.store.GetUser(ctx, username)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, user)
}
