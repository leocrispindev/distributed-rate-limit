package delivery

import (
	"concurrency-hazelcast/internal/handler"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {

	handler := handler.NewHandler()

	// Define your routes here
	router.GET("/example", handler.HelloHandler)

}
