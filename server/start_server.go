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

	serverPort := utils.GetEnv("SERVER_PORT", "8080")

	username := utils.GetEnv("DB_USERNAME", "public_user")
	password := utils.GetEnv("DB_PASSWORD", "public@123")
	host := utils.GetEnv("DB_HOST", "dmt1.intellicar.in")
	port := utils.GetEnv("DB_PORT", "3306")
	schema := utils.GetEnv("DB_SCHEMA", "dmt")

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
