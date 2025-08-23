package domain

import "context"

type NoSQLRepository interface {
	Set(ctx context.Context, key string, value interface{}) error
	Get(ctx context.Context, key string) (interface{}, error)
	SetWithTTL(ctx context.Context, key string, value interface{}, ttlSeconds int) error
}
