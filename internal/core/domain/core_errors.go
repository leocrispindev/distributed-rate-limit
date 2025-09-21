package domain

type BucketNotFoundError string

func (e BucketNotFoundError) Error() string {
	return string(e)
}

type LockError string

func (e LockError) Error() string {
	return string(e)
}

type UnlockError string

func (e UnlockError) Error() string {
	return string(e)
}

type BucketAlreadyExistsError string

func (e BucketAlreadyExistsError) Error() string {
	return string(e)
}

type ErrorOnUpdateBucket string

func (e ErrorOnUpdateBucket) Error() string {
	return string(e)
}
