package server

import (
	"log"
	"vds/config"
	"vds/routes"
	"vds/utils"

	"github.com/gin-gonic/gin"
)

func StartServer() {
	utils.LoadEnv()

	serverPort := utils.GetEnv("SERVER_PORT", "")

	username := utils.GetEnv("DB_USERNAME", "")
	password := utils.GetEnv("DB_PASSWORD", "")
	host := utils.GetEnv("DB_HOST", "")
	port := utils.GetEnv("DB_PORT", "")
	schema := utils.GetEnv("DB_SCHEMA", "")

	err := config.ConnectDB(username, password, host, port, schema)
	if err != nil {
		log.Fatal(err)
	}
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	routes.ModelRoutes(r)
	log.Println("Server started on port", serverPort)
	r.Run(":" + serverPort)
}
