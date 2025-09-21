package domain

import (
	"encoding/json"
	"time"
)

const (
	BucketPrefix       = "bucket_"
	ClientBucketPrefix = "client_"
	maxTokens          = 100
	refillRate         = 30
)

type BucketRequest struct {
	Name string `json:"name"`
}

type BucketResponse struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type Bucket struct {
	Name           string `json:"name"`
	ApiId          string `json:"api_id"`
	Tokens         int    `json:"tokens"`
	LastRefillTime int64  `json:"last_refill_time"`
}

func NewClientBucket() *Bucket {

	return &Bucket{
		Tokens:         maxTokens,
		LastRefillTime: time.Now().UnixMilli(),
	}
}

func (bucket *Bucket) hasAvailableTokens() bool {
	return bucket.Tokens > 0
}

func (bucket *Bucket) consumeToken() {
	if bucket.Tokens > 0 {
		bucket.Tokens--
	}
}

func (bucket *Bucket) AcquireToken() bool {
	bucket.refillIfNeeded()

	if bucket.hasAvailableTokens() {
		bucket.consumeToken()
		return true
	}
	return false
}

func (bucket *Bucket) refillIfNeeded() {
	now := time.Now().UnixMilli()
	if bucket.shouldRefill(now) {
		bucket.Tokens = maxTokens
		bucket.LastRefillTime = now
	}
}

func (bucket *Bucket) shouldRefill(now int64) bool {
	elapsed := now - bucket.LastRefillTime
	return elapsed >= refillRate*1000
}

func (bucket *Bucket) ToByteArray() ([]byte, error) {
	return json.Marshal(bucket)
}
