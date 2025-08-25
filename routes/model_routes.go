package routes

import (
	"vds/handler"

	"github.com/gin-gonic/gin"
)

func ModelRoutes(r *gin.Engine) {
	r.GET("/models", handler.GetModelHandler)
	r.GET("/models/count", handler.GetModelCountHandler)
}
