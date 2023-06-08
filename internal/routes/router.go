package routes

import (
	"web-app/internal/logger"

	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	r := gin.New()
	r.Use(logger.GinLog(), gin.Recovery())

	return r
}
