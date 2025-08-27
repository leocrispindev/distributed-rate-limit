package domain_test

import (
	"concurrency-hazelcast/internal/core/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ClientTokensBucketTestSuite struct {
	suite.Suite
	bucket *domain.Bucket
}

func TestClientTokensBucketTestSuite(t *testing.T) {
	suite.Run(t, new(ClientTokensBucketTestSuite))
}

func (suite *ClientTokensBucketTestSuite) TestAcquireTokenShouldConsumeToken() {
	suite.bucket = domain.NewClientBucket()

	suite.bucket.Tokens = 3

	acquired := suite.bucket.AcquireToken()

	assert.True(suite.T(), acquired)
	assert.Equal(suite.T(), 2, suite.bucket.Tokens)
}

func (suite *ClientTokensBucketTestSuite) TestAcquireTokenShouldNotConsumeTokenWhenEmpty() {
	suite.bucket = domain.NewClientBucket()

	suite.bucket.Tokens = 0

	acquired := suite.bucket.AcquireToken()

	assert.False(suite.T(), acquired)
	assert.Equal(suite.T(), 0, suite.bucket.Tokens)
}

func (suite *ClientTokensBucketTestSuite) TestToByteArrayShouldMarshalBucket() {
	suite.bucket = domain.NewClientBucket()

	suite.bucket.Tokens = 3
	suite.bucket.LastRefillTime = 1234567890

	data, err := suite.bucket.ToByteArray()

	assert.NoError(suite.T(), err)
	assert.JSONEq(suite.T(), `{"tokens":3,"last_refill_time":1234567890}`, string(data))
}
