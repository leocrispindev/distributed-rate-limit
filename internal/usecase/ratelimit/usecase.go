package ratelimit

import (
	"concurrency-hazelcast/internal/core/domain"
	"context"
	"encoding/json"
	"log"
)

const ()

type RateLimit interface {
	AllowAccess(ctx context.Context, clientId string) bool
}

type RateLimitUseCase struct {
	repository        domain.NoSQLRepository
	metricsRepository domain.MetricsRepository
	lockRepository    domain.LockRepository
}

func NewRateLimitUseCase(rep domain.NoSQLRepository, metrics domain.MetricsRepository, lockRep domain.LockRepository) *RateLimitUseCase {
	return &RateLimitUseCase{
		repository:        rep,
		metricsRepository: metrics,
		lockRepository:    lockRep,
	}
}

func (uc *RateLimitUseCase) AllowAccess(ctx context.Context, clientId string) (bool, error) {

	key := uc.getBucketKey(clientId)

	uc.lockRepository.Lock(ctx, key)
	defer uc.lockRepository.Unlock(ctx, key)

	bucket, err := uc.repository.Get(ctx, key)
	if err != nil {
		log.Println("error on get tokens for for=[" + key + "]")
		return false, err
	}

	if bucket == nil {
		return false, domain.BucketNotFoundError("bucket not found")
	}

	var clientBucket *domain.Bucket
	json.Unmarshal(bucket.([]byte), &clientBucket)
	if clientBucket == nil {
		log.Println("error on parse bucket for for=[" + key + "]")
		return false, nil
	}

	uc.metricsRepository.CountMetric(clientBucket.Name)

	avaliableToken := clientBucket.AcquireToken()
	if !avaliableToken {
		return false, nil
	}

	payload, _ := clientBucket.ToByteArray()

	err = uc.repository.Set(ctx, key, payload)
	if err != nil {
		log.Println("error on save bucket for for=[" + key + "]")
		return false, domain.ErrorOnUpdateBucket("error on save bucket for for=[" + key + "]")
	}

	return avaliableToken, nil

}

func (uc *RateLimitUseCase) getBucketKey(clientId string) string {
	return domain.BucketPrefix + clientId
}
