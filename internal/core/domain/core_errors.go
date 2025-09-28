package domain

import "errors"

var (
	ErrBucketNotFound      = errors.New("bucket not found")
	ErrBucketAlreadyExists = errors.New("bucket already exists")
	ErrLockFailed          = errors.New("failed to acquire lock")
	ErrUnlockFailed        = errors.New("failed to release lock")
	ErrBucketUpdateFailed  = errors.New("failed to update bucket")
	ErrSetKeyFailed        = errors.New("failed to set key in repository")
)

type CoreError struct {
	Type    error
	Message string
	Err     error
}

func (e *CoreError) Error() string {
	if e.Err != nil {
		return e.Type.Error() + ": " + e.Message + " -> " + e.Err.Error()
	}
	return e.Type.Error() + ": " + e.Message
}

func (e *CoreError) Unwrap() error {
	return e.Type
}

func NewBucketNotFoundError(message string, err error) *CoreError {
	return &CoreError{
		Type:    ErrBucketNotFound,
		Message: message,
		Err:     err,
	}
}

func NewBucketAlreadyExistsError(message string, err error) *CoreError {
	return &CoreError{
		Type:    ErrBucketAlreadyExists,
		Message: message,
		Err:     err,
	}
}

func NewLockFailedError(message string, err error) *CoreError {
	return &CoreError{
		Type:    ErrLockFailed,
		Message: message,
		Err:     err,
	}
}

func NewUnlockFailedError(message string, err error) *CoreError {
	return &CoreError{
		Type:    ErrUnlockFailed,
		Message: message,
		Err:     err,
	}
}

func NewBucketUpdateFailedError(message string, err error) *CoreError {
	return &CoreError{
		Type:    ErrBucketUpdateFailed,
		Message: message,
		Err:     err,
	}
}

func NewSetKeyFailedError(message string, err error) *CoreError {
	return &CoreError{
		Type:    ErrSetKeyFailed,
		Message: message,
		Err:     err,
	}
}
