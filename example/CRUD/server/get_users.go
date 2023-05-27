package server

import "github.com/gin-gonic/gin"

func (server *Server) GetUser(ctx *gin.Context) {
	server.store.GetUser(ctx, "")
}
	