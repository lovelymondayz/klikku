FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata && adduser -D -u 1000 klikku

WORKDIR /app
COPY --from=builder /app/server .
COPY migrations/ ./migrations/

USER klikku

EXPOSE 8083

CMD ["./server"]
