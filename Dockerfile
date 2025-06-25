FROM --platform=linux/amd64 docker.io/library/golang:1.24.4-alpine AS builder

RUN apk add --no-cache musl-dev gcc sqlite-dev

WORKDIR /app

COPY go.mod go.sum ./
COPY cmd/wheel-of-wonder/ /app/cmd/wheel-of-wonder/
COPY internal /app/internal

RUN go mod download

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags='-s -w -extldflags "-static"' -o /app/server ./cmd/wheel-of-wonder

FROM docker.io/library/alpine

RUN apk add --no-cache tzdata

WORKDIR /app

COPY --from=builder /app/server ./server
COPY --from=builder /app/internal/db/migrations ./internal/db/migrations

CMD ["/app/server"]