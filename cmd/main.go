package main

import (
	"concurrency-hazelcast/internal/delivery"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	delivery.SetupRoutes(router)

	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}
