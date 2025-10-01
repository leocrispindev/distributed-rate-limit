# Distributed Lock and Rate Limiting API

This project, built in **Golang**, implements concurrency control with **distributed locking** and **per-client rate limiting**, using **Hazelcast** for shared state management and **Nginx** as a load balancer.

The rate limiting algorithm used is **Fixed Window Time**, where each client is allowed a fixed number of requests within a given time window. Once the window expires, the quota is automatically reset, providing both simplicity in access control and predictability in usage.

## Architecture

- **Load Balancer (Nginx)** → routes incoming traffic.  
- **Rate Limit Cluster (Go APIs)** → multiple replicas for resilience and scalability.  
- **Hazelcast Cluster** → in-memory storage + distributed locks.  
- **Observability** → Prometheus (metrics) + Grafana (dashboards).  

![Architecture](archtecture.png)

---

## Resilience: Multiple Instances

### API
- Prevents downtime in case of failures (**high availability**).  
- Enables **horizontal scalability**.  
- Supports **rolling updates** with zero downtime.  
- Stateless (state stored in Hazelcast).  

### Hazelcast
- **Data replication**: protects against data loss during failures.  
- **Quorum**: ensures cluster consistency.  
- **Load distribution**: maintains low latency under heavy traffic.  
- **Distributed locks**: prevents race conditions.  

👉 Production recommendation: **minimum of 3 Hazelcast nodes** + **2–3 API replicas**.  

---

## Hazelcast Usage

1. **In-memory storage** → holds client state (counter, window, bucket).  
2. **Distributed lock** → enforces mutual exclusion per client.  
3. **Unique client key** → same key for both lock and data → eliminates race conditions.  

---

## Endpoints

#### Validate Rate Limit Middleware

```bash
curl -X GET http://localhost:9999/example -H "X-Api-Id: <client-id>"
```

#### Response
```json
{
   "message": "Hello World"
}
```

#### Create Client
```bash
curl -X POST http://localhost:9999/bucket -H "Content-Type: application/json"   -d '{"name":"client123"}'
```

#### Response 
```json
{
  "name": "client123",
  "id": "efc8cd40-2574-489f-8be8-79802b5a623a"
}
```

---

### Local Execution (Docker Compose)
```bash
git clone https://github.com/leocrispindev/distributed-rate-limit.git
cd distributed-rate-limit
```
or 
```txt
copy the docker-compose file
```
1.1
```bash
docker-compose up -d

```




### Testing
You can test the application using **K6**:
```bash
k6 run test-rate-limit.js
```
