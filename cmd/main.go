package main

import (
	"concurrency-hazelcast/internal/delivery"

	"github.com/gin-gonic/gin"
)

const (
	MaxTokens     = 100 // limite por minuto
	RefillRate    = 100 // tokens por minuto
	RefillSeconds = 60
)

func main() {
	router := gin.Default()

	delivery.SetupRoutes(router)

	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}
