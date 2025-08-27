package middleware

import (
	"concurrency-hazelcast/internal/usecase/ratelimit"
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

		allowed := m.ratelimitUseCase.AllowAccess(ctx.Request.Context(), clientId)
		if !allowed {
			ctx.JSON(429, gin.H{"error": "Rate limit exceeded"})
			ctx.Abort()
			return
		}

		ctx.Next()
	}

}
