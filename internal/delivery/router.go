package delivery

import (
	"concurrency-hazelcast/internal/delivery/middleware"
	"concurrency-hazelcast/internal/handler"
	"concurrency-hazelcast/internal/infrastructure/repository/hazelcast"
	createtokenbucket "concurrency-hazelcast/internal/usecase/bucket/createTokenBucket"
	"concurrency-hazelcast/internal/usecase/ratelimit"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {

	hzClient, _ := hazelcast.NewHazelcastClient()
	NoSQLRepository := hazelcast.NewHazelcastRepository(hzClient)

	// Create Token Bucket
	createUseCase := createtokenbucket.NewCreateTokenBucketUseCase(NoSQLRepository)
	createBucketHandler := handler.NewCreateBucketHandler(createUseCase)

	//Middleware
	rateLimitUseCase := ratelimit.NewRateLimitUseCase(NoSQLRepository)

	authorizationMiddleware := middleware.NewAuthorizationMiddleware(rateLimitUseCase)

	router.GET("/example", authorizationMiddleware.Middleware(), handler.HelloHandler)
	router.POST("/bucket", createBucketHandler.Handler)

}
