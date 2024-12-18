# SE400

Learning Repository for Benchmarking and comparison between different API Frameworks

# Monitoring Stack

A monitoring setup with Prometheus and Grafana for Spring Boot (Java) and Gin (Go) applications.

## Services

- Flask API (port 5500)
- Gin API (port 8090)
- Java API (port 8091)
- Prometheus (port 9090)
- Grafana (port 3000)
- Node Exporter (port 9100)
- cAdvisor (port 8080)
- postgresql for go (port 15432)
- postgresql for java (port 25432)

## Setup

1. Make sure Docker and Docker Compose are installed
2. Run `docker-compose up -d`
3. Access Grafana at http://localhost:3000 (admin/admin)
4. Create new Dashboard and query metrics to your needs
5. To run the request client, navigate to `client` and run `go build -o client main.go metrics.go requests.go` then run the binary using `./client -target="http://localhost:8090/orders/1" -target="http://localhost:8091/orders/1" -maxClients=200 -scaleInterval=500`. Replace the -target flags with URL to your needs.
## API

1. Get order by id

```
curl --location 'localhost:8091/orders/1'
```

2. Place order

```
curl --location 'localhost:8091/orders' \
--header 'Content-Type: application/json' \
--data '{
    "productId": 1,
    "product": {
        "id": 1
    },
    "quantity": 3,
    "customer_id": 123
}'
```
