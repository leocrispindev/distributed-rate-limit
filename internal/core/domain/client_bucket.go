package domain

import (
	"encoding/json"
	"time"
)

const (
	maxTokens  = 5
	refillRate = 1000 //miliseconds
)

type ClientTokensBucket struct {
	Tokens         int   `json:"tokens"`
	LastRefillTime int64 `json:"last_refill_time"`
}

func NewClientBucket() *ClientTokensBucket {

	return &ClientTokensBucket{
		Tokens:         maxTokens,
		LastRefillTime: time.Now().UnixMilli(),
	}
}

func (bucket *ClientTokensBucket) hasAvailableTokens() bool {
	return bucket.Tokens > 0
}

func (bucket *ClientTokensBucket) consumeToken() {
	if bucket.Tokens > 0 {
		bucket.Tokens--
	}
}

func (bucket *ClientTokensBucket) AcquireToken() bool {
	bucket.refillIfNeeded()

	if bucket.hasAvailableTokens() {
		bucket.consumeToken()
		return true
	}
	return false
}

func (bucket *ClientTokensBucket) refillIfNeeded() {
	now := time.Now().UnixMilli()
	if bucket.shouldRefill(now) {
		bucket.Tokens = maxTokens
		bucket.LastRefillTime = now
	}
}

func (bucket *ClientTokensBucket) shouldRefill(now int64) bool {
	elapsed := now - bucket.LastRefillTime
	return elapsed > refillRate
}

func (bucket *ClientTokensBucket) ToByteArray() ([]byte, error) {
	return json.Marshal(bucket)
}
