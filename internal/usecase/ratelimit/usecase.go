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
	repository        domain.NoSQLRepository
	metricsRepository domain.MetricsRepository
}

func NewRateLimitUseCase(rep domain.NoSQLRepository, metrics domain.MetricsRepository) *RateLimitUseCase {
	return &RateLimitUseCase{
		repository:        rep,
		metricsRepository: metrics,
	}
}

func (uc *RateLimitUseCase) AllowAccess(ctx context.Context, clientId string) bool {

	key := uc.getBucketKey(clientId)

	bucket, err := uc.repository.Get(ctx, key)
	if err != nil {
		log.Println("error on get tokens for for=[" + key + "]")
		return false
	}

	clientBucket := bucket.(domain.Bucket)

	defer uc.metricsRepository.CountMetric(clientBucket.Name)

	avaliableToken := clientBucket.AcquireToken()
	if !avaliableToken {
		return false
	}

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
