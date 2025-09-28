package delivery

import (
	"github.com/leocrispindev/distributed-rate-limit/internal/delivery/middleware"
	"github.com/leocrispindev/distributed-rate-limit/internal/handler"
	"github.com/leocrispindev/distributed-rate-limit/internal/infrastructure/repository/hazelcast"
	"github.com/leocrispindev/distributed-rate-limit/internal/infrastructure/repository/prometheus"
	createtokenbucket "github.com/leocrispindev/distributed-rate-limit/internal/usecase/bucket/createTokenBucket"
	"github.com/leocrispindev/distributed-rate-limit/internal/usecase/ratelimit"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	//Metric
	prometheusRepository := prometheus.NewPrometheusRepository()

	//Hazelcast NoSQL
	hzClient, _ := hazelcast.NewHazelcastClient()
	NoSQLRepository := hazelcast.NewHazelcastRepository(hzClient)
	//Hazelcast Lock
	lockRepository := hazelcast.NewHazelcastLockRepository(hzClient)

	// Create Token Bucket
	createUseCase := createtokenbucket.NewCreateTokenBucketUseCase(NoSQLRepository)
	createBucketHandler := handler.NewCreateBucketHandler(createUseCase)

	//Middleware
	rateLimitUseCase := ratelimit.NewRateLimitUseCase(NoSQLRepository, prometheusRepository, lockRepository)

	authorizationMiddleware := middleware.NewAuthorizationMiddleware(rateLimitUseCase)

	router.GET("/example", authorizationMiddleware.Middleware(), handler.HelloHandler)
	router.POST("/bucket", createBucketHandler.Handler)

	//metrics
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	//health
	router.GET("/health", handler.HealthHandler)

}
