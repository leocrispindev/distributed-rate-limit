package middleware

import (
	"concurrency-hazelcast/internal/core/domain"
	"concurrency-hazelcast/internal/usecase/ratelimit"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthorizationMiddleware struct {
	ratelimitUseCase *ratelimit.RateLimitUseCase
}

func NewAuthorizationMiddleware(ratelimitUseCase *ratelimit.RateLimitUseCase) *AuthorizationMiddleware {
	return &AuthorizationMiddleware{ratelimitUseCase: ratelimitUseCase}
}

func (m *AuthorizationMiddleware) Middleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		clientId := ctx.GetHeader("X-Api-Id")
		if clientId == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Api ID is required"})
			ctx.Abort()
			return
		}

		allowed, err := m.ratelimitUseCase.AllowAccess(ctx.Request.Context(), clientId)

		if err != nil {
			if errors.Is(err, domain.BucketNotFoundError("bucket not found")) {
				ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				ctx.Abort()
				return
			}

			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			ctx.Abort()
			return

		}

		if !allowed {
			ctx.JSON(429, gin.H{"error": "Rate limit exceeded"})
			ctx.Abort()
			return
		}

		ctx.Next()
	}

}
