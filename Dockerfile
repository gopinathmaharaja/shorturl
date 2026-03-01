FROM golang:1.20-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o shorturl \
    cmd/main.go

# Runtime stage
FROM alpine:3.18

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/shorturl .

EXPOSE 3000

CMD ["./shorturl"]
