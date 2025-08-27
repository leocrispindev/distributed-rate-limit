package createtokenbucket

import (
	"concurrency-hazelcast/internal/core/domain"
	"context"
	"errors"

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
		return nil, errors.New("bucket already exists")
	}

	newBucket := *domain.NewClientBucket()
	newBucket.Name = bucket.Name
	newBucket.ApiId = uuid.New().String()

	if err := uc.repo.Set(ctx, clientKey, true); err != nil {
		return nil, errors.New("error creating bucket")

	}

	bucketKey := domain.BucketPrefix + newBucket.ApiId

	if err := uc.repo.Set(ctx, bucketKey, newBucket); err != nil {
		return nil, errors.New("error creating bucket")
	}

	return &domain.BucketResponse{
		Name: newBucket.Name,
		ID:   newBucket.ApiId,
	}, nil
}
