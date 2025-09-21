package hazelcast

import (
	"context"
	"log"

	"github.com/hazelcast/hazelcast-go-client"
)

type HazelcastLockRepository struct {
	hzClientMap *hazelcast.Map
}

func NewHazelcastLockRepository(hzClient *hazelcast.Client) *HazelcastLockRepository {
	hzClientMap, err := hzClient.GetMap(context.TODO(), mapName)
	if err != nil {
		log.Fatalf("Error on get clients map: %v", err)
	}

	return &HazelcastLockRepository{
		hzClientMap: hzClientMap,
	}
}

func (hz *HazelcastLockRepository) Lock(ctx context.Context, key string) (string, error) {

	err := hz.hzClientMap.Lock(ctx, key)
	if err != nil {
		log.Printf("Erro ao obter lock '%s': %v", key, err)
		return "", err
	}

	log.Printf("Lock acquire with success '%s'", key)
	return "", nil
}

func (hz *HazelcastLockRepository) Unlock(ctx context.Context, key string) error {
	err := hz.hzClientMap.Unlock(ctx, key)
	if err != nil {
		log.Printf("Erro ao liberar lock '%s': %v", key, err)
		return err
	}

	return nil
}
