package handler

import (
	"net/http"
	"vds/service"

	"github.com/gin-gonic/gin"
)

type ApiResponse struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func GetModelHandler(c *gin.Context) {
	models, err := service.GetUniqueModels()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ApiResponse{
			Error: "Error fetching models",
		})
		return
	}
	c.JSON(http.StatusOK, ApiResponse{
		Message: "Models fetched successfully",
		Data: map[string]interface{}{
			"count":  len(models),
			"models": models,
		},
	})
}

func GetModelCountHandler(c *gin.Context) {
	models, err := service.GetUniqueModelsCount()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ApiResponse{
			Error: "Error fetching models count",
		})
		return
	}
	c.JSON(http.StatusOK, ApiResponse{
		Message: "Models count fetched successfully",
		Data:    models,
	})
}
