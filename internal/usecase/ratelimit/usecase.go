package ratelimit

import (
	"concurrency-hazelcast/internal/core/domain"
	"context"
	"log"
)

const ()

type RateLimit interface {
	AllowAccess(ctx context.Context, clientId string) bool
}

type RateLimitUseCase struct {
	repository domain.NoSQLRepository
}

func NewRateLimitUseCase(rep domain.NoSQLRepository) *RateLimitUseCase {
	return &RateLimitUseCase{repository: rep}
}

func (uc *RateLimitUseCase) AllowAccess(ctx context.Context, clientId string) bool {
	key := uc.getBucketKey(clientId)

	bucket, err := uc.repository.Get(ctx, key)
	if err != nil {
		log.Println("error on get tokens for for=[" + key + "]")
		return false
	}

	clientBucket := bucket.(domain.Bucket)

	avaliableToken := clientBucket.AcquireToken()
	if !avaliableToken {
		return false
	}

	/*data, err := clientBucket.ToByteArray()
	if err != nil {
		log.Println("error on marshal bucket for for=[" + key + "]")
		return false
	}*/

	err = uc.repository.Set(ctx, key, clientBucket)
	if err != nil {
		log.Println("error on save bucket for for=[" + key + "]")
		return false
	}

	return avaliableToken

}

func (uc *RateLimitUseCase) getBucketKey(clientId string) string {
	return domain.BucketPrefix + clientId
}
