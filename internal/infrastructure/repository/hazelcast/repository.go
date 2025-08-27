package hazelcast

import (
	"context"
	"log"
	"time"

	"github.com/hazelcast/hazelcast-go-client"
)

const (
	mapName = "buckets"
)

type HazelcastRepository struct {
	hzClientMap *hazelcast.Map
}

func NewHazelcastClient() (*hazelcast.Client, error) {
	config := hazelcast.Config{}

	config.Cluster.Network.SetAddresses("127.0.0.1:5701")

	// Conecta ao cluster
	client, err := hazelcast.StartNewClientWithConfig(context.TODO(), config)
	if err != nil {
		log.Fatalf("Erro ao conectar ao Hazelcast: %v", err)
	}

	return client, nil

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

func (hz *HazelcastRepository) Delete(ctx context.Context, key string) error {
	err := hz.hzClientMap.Delete(ctx, key)
	return err
}
