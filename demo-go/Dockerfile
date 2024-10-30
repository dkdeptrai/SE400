FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o demo-go ./cmd/demo-go

FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/demo-go .
COPY --from=builder /app/configs/development.yaml ./configs/development.yaml

EXPOSE 8090

CMD ["./demo-go"]
