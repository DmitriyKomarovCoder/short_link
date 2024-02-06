FROM golang:1.21.0-alpine AS builder

WORKDIR /app

COPY . .

RUN apk --no-cache add git \
    && go mod download \
    && go clean --modcache \
    && CGO_ENABLED=0 GOOS=linux go build -mod=readonly -o server ./cmd/http/httpApi.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/ .

EXPOSE 8080 8081

CMD ["./server"]
