FROM golang:1.23-alpine
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
RUN go build -o app .
EXPOSE 8090
CMD ["sh", "-c", "exec /app/app"]