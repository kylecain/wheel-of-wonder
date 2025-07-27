# Wheel of Wonder

## Go

```zsh
go run cmd/wheel-of-wonder/main.go 
```

## Podman

Local Image

```zsh
podman build -t wheel-of-wonder:local .
podman run \
    -e PUID=$(id -u) \
    -e PGID=$(id -g) \
    -v $(pwd)/data:/app/data \
    --env-file .env \
    wheel-of-wonder:local
```

Remote Image

```zsh
podman pull ghcr.io/kylecain/wheel-of-wonder:latest
podman run \
    -e PUID=$(id -u) \
    -e PGID=$(id -g) \
    -e BOT_TOKEN \
    -e GUILD_ID \
    -e MIGRATION_URL \
    -e DATABASE_URL \
    -v $(pwd)/data:/app/data \
    ghcr.io/kylecain/wheel-of-wonder:latest
```

## GHCR

```zsh
podman build -t wheel-of-wonder:latest .
echo $CR_PAT | podman login ghcr.io -u kylecain --password-stdin
podman tag wheel-of-wonder:latest ghcr.io/kylecain/wheel-of-wonder:latest
podman push ghcr.io/kylecain/wheel-of-wonder:latest
```
