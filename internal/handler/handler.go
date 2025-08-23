package handler

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hazelcast/hazelcast-go-client"
)

const (
	BucketMaxSize     = 5
	BucketRefreshRate = 1 // em segundos
)

type Handler struct {
	clientMap *hazelcast.Map
}

type BucketTokens struct {
	AvailableTokens int   `json:"available_tokens"`
	LastRefillTime  int64 `json:"last_refill_time"`
}

func NewHandler() *Handler {

	config := hazelcast.Config{}

	config.Cluster.Network.SetAddresses("127.0.0.1:5701")

	// Conecta ao cluster
	client, err := hazelcast.StartNewClientWithConfig(context.TODO(), config)
	if err != nil {
		log.Fatalf("Erro ao conectar ao Hazelcast: %v", err)
	}

	clientMap, err := client.GetMap(context.TODO(), "my-distributed-map")

	if err != nil {
		log.Fatalf("Erro ao obter o mapa distribuído: %v", err)
	}

	return &Handler{
		clientMap: clientMap,
	}
}

// HelloHandler responds with a greeting message.
func (h *Handler) HelloHandler(c *gin.Context) {

	ctx := c.Request.Context()
	clientId := c.GetHeader("X-Client-ID")
	if clientId == "" {
		c.JSON(400, gin.H{"error": "X-Client-ID header is required"})
		return
	}

	//O algoritmo utilizado é Bucket Token, algortimo para controle de taxa de requisições
	//a cada segundo o bucket é recarregado com 100 tokens
	// se existir token disponivel, consome um token e retorna a mensagem
	// se não existir token, retorna 429 Too Many Requests

	bucket, err := h.clientMap.Get(ctx, clientId)

	if err != nil {
		log.Printf("Erro ao obter tokens do cliente %s: %v", clientId, err)
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	var bucketTokens BucketTokens
	if bucket == nil {
		// Inicializa o bucket para o cliente
		bucketTokens = BucketTokens{
			AvailableTokens: BucketMaxSize - 1, // consome um token
			LastRefillTime:  time.Now().UnixMilli(),
		}

		value, _ := json.Marshal(bucketTokens)

		err = h.clientMap.Set(ctx, clientId, value)
		if err != nil {
			log.Printf("Erro ao inicializar tokens do cliente %s: %v", clientId, err)
			c.JSON(500, gin.H{"error": "Internal server error"})
			return
		}

		c.JSON(200, gin.H{"message": "Hello, World!"})
		return

	}
	// Converte o valor armazenado no mapa para BucketTokens
	err = json.Unmarshal(bucket.([]byte), &bucketTokens)
	if err != nil {
		log.Printf("Erro ao desserializar tokens do cliente %s: %v", clientId, err)
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	timeNow := time.Now().UnixMilli() // converte para milisegundos
	elapsed := timeNow - bucketTokens.LastRefillTime

	if elapsed > BucketRefreshRate*1000 { // Se passou mais de 1 segundo desde o ultimo consumo
		// Recarrega o bucket
		bucketTokens.AvailableTokens = BucketMaxSize
		bucketTokens.LastRefillTime = timeNow
	}

	if bucketTokens.AvailableTokens > 0 {
		bucketTokens.AvailableTokens--
		err = h.clientMap.Set(ctx, clientId, bucketTokens)
		if err != nil {
			log.Printf("Erro ao atualizar tokens do cliente %s: %v", clientId, err)
			c.JSON(500, gin.H{"error": "Internal server error"})
			return
		}
		c.JSON(200, gin.H{"message": "Hello, World!"})
		return

	}

	c.JSON(429, gin.H{"error": "Too Many Requests"})
	// Tempo decorrido desde o ultimo consumo de token
}
