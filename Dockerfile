FROM golang:1.21 AS builder
ENV GO111MODULE=on
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o mongo-streams

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/mongo-streams /app/mongo-streams
RUN chmod +x /app/mongo-streams

# Run the binary
ENTRYPOINT ["/app/mongo-streams"]