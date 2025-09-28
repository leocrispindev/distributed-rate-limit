package createtokenbucket

import (
	"context"
	"errors"

	"github.com/leocrispindev/distributed-rate-limit/internal/core/domain"

	"github.com/google/uuid"
)

type CreateTokenBucketUseCase struct {
	repo domain.NoSQLRepository
}

func NewCreateTokenBucketUseCase(repo domain.NoSQLRepository) *CreateTokenBucketUseCase {
	return &CreateTokenBucketUseCase{repo: repo}
}

func (uc *CreateTokenBucketUseCase) Create(ctx context.Context, bucket domain.BucketRequest) (*domain.BucketResponse, error) {
	if bucket.Name == "" {
		return nil, errors.New("bucket name is required")
	}

	clientKey := domain.ClientBucketPrefix + bucket.Name

	exists, err := uc.repo.Get(ctx, clientKey)
	if err != nil {
		return nil, errors.New("error checking if bucket exists")
	}

	if exists != nil {
		return nil, domain.NewBucketAlreadyExistsError(bucket.Name, nil)
	}

	newBucket := *domain.NewClientBucket()
	newBucket.Name = bucket.Name
	newBucket.ApiId = uuid.New().String()

	if err := uc.repo.Set(ctx, clientKey, true); err != nil {
		return nil, err

	}

	bucketKey := domain.BucketPrefix + newBucket.ApiId

	payload, _ := newBucket.ToByteArray()

	if err := uc.repo.Set(ctx, bucketKey, payload); err != nil {
		return nil, err
	}

	return &domain.BucketResponse{
		Name: newBucket.Name,
		ID:   newBucket.ApiId,
	}, nil
}
