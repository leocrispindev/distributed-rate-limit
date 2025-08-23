package hazelcast

import (
	"context"
	"log"
	"time"

	"github.com/hazelcast/hazelcast-go-client"
)

const (
	mapName = "client-buckets"
)

type HazelcastRepository struct {
	hzClientMap *hazelcast.Map
}

func NewHazelcastRepository(hzClient *hazelcast.Client) *HazelcastRepository {
	hzClientMap, err := hzClient.GetMap(context.TODO(), mapName)
	if err != nil {
		log.Fatalf("Error on get clients map: %v", err)
	}

	return &HazelcastRepository{
		hzClientMap: hzClientMap,
	}
}

func (hz *HazelcastRepository) Get(ctx context.Context, key string) (interface{}, error) {
	return hz.hzClientMap.Get(ctx, key)
}

func (hz *HazelcastRepository) Set(ctx context.Context, key string, value interface{}) error {
	return hz.hzClientMap.Set(ctx, key, value)
}

func (hz *HazelcastRepository) SetWithTTL(ctx context.Context, key string, value interface{}, ttlSeconds int) error {
	ttl := time.Duration(ttlSeconds) * time.Second

	return hz.hzClientMap.SetWithTTL(ctx, key, value, ttl)
}
