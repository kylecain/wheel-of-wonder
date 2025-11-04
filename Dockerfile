FROM --platform=linux/amd64 docker.io/library/golang:1.25.3-alpine AS builder

RUN apk add --no-cache musl-dev gcc sqlite-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/wheel-of-wonder/ /app/cmd/wheel-of-wonder/
COPY internal /app/internal

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags='-s -w -extldflags "-static"' -o /app/server ./cmd/wheel-of-wonder

FROM docker.io/library/alpine:3.20

RUN apk add --no-cache tzdata

WORKDIR /app

COPY --from=builder /app/server ./server
COPY --from=builder /app/internal/db/migrations ./internal/db/migrations

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
CMD ["/app/server"]