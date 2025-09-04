package delivery

import (
	"concurrency-hazelcast/internal/delivery/middleware"
	"concurrency-hazelcast/internal/handler"
	"concurrency-hazelcast/internal/infrastructure/repository/hazelcast"
	"concurrency-hazelcast/internal/infrastructure/repository/prometheus"
	createtokenbucket "concurrency-hazelcast/internal/usecase/bucket/createTokenBucket"
	"concurrency-hazelcast/internal/usecase/ratelimit"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	//Metric
	prometheusRepository := prometheus.NewPrometheusRepository()

	//Hazelcast NoSQL
	hzClient, _ := hazelcast.NewHazelcastClient()
	NoSQLRepository := hazelcast.NewHazelcastRepository(hzClient)

	// Create Token Bucket
	createUseCase := createtokenbucket.NewCreateTokenBucketUseCase(NoSQLRepository)
	createBucketHandler := handler.NewCreateBucketHandler(createUseCase)

	//Middleware
	rateLimitUseCase := ratelimit.NewRateLimitUseCase(NoSQLRepository, prometheusRepository)

	authorizationMiddleware := middleware.NewAuthorizationMiddleware(rateLimitUseCase)

	router.GET("/example", authorizationMiddleware.Middleware(), handler.HelloHandler)
	router.POST("/bucket", createBucketHandler.Handler)

	//metrics
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

}
