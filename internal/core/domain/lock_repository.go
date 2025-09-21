package domain

import "context"

type LockRepository interface {
	Lock(ctx context.Context, lockId string) (string, error)
	Unlock(ctx context.Context, lockId string) error
}
