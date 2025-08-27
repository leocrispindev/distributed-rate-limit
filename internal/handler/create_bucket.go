package handler

import (
	"concurrency-hazelcast/internal/core/domain"
	createtokenbucket "concurrency-hazelcast/internal/usecase/bucket/createTokenBucket"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateBucketHandler struct {
	createUseCase *createtokenbucket.CreateTokenBucketUseCase
}

func NewCreateBucketHandler(createUseCase *createtokenbucket.CreateTokenBucketUseCase) *CreateBucketHandler {
	return &CreateBucketHandler{createUseCase: createUseCase}
}

func (h *CreateBucketHandler) Handler(ctx *gin.Context) {
	var bucket domain.BucketRequest
	if err := ctx.ShouldBindJSON(&bucket); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if bucket.Name == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Bucket name is required"})
		return

	}

	resp, err := h.createUseCase.Create(ctx.Request.Context(), bucket)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, resp)

}
