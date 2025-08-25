package server

import (
	"log"
	"vds/routes"
	"vds/utils"

	"github.com/gin-gonic/gin"
)

func StartServer() {
	utils.LoadEnv()

	serverPort := utils.GetEnv("SERVER_PORT", "")

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	routes.ModelRoutes(r)
	log.Println("Server started on port", serverPort)
	r.Run(":" + serverPort)
}
