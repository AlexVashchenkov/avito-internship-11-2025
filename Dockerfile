FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git make

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/main.go

FROM alpine:3.22.2

WORKDIR /app

RUN apk add --no-cache postgresql16-client

COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/server .
COPY .env .env
EXPOSE 8080

ENTRYPOINT ["./server"]