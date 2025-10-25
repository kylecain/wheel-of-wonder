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
    -e PUID=99 \
    -e PGID=100 \
    -e UMASK=022 \
    -v $(pwd)/data:/app/data \
    --env-file .env \
    wheel-of-wonder:local
```

Remote Image

```zsh
podman pull ghcr.io/kylecain/wheel-of-wonder:latest
podman run \
    -e PUID=99 \
    -e PGID=100 \
    -e UMASK=022 \
    -e BOT_TOKEN \
    -e GUILD_ID \
    -v $(pwd)/data:/app/data \
    ghcr.io/kylecain/wheel-of-wonder:latest
```

## GHCR

```zsh
podman build --arch amd64 -t wheel-of-wonder:latest .
echo $CR_PAT | podman login ghcr.io -u kylecain --password-stdin
podman tag wheel-of-wonder:latest ghcr.io/kylecain/wheel-of-wonder:latest
podman push ghcr.io/kylecain/wheel-of-wonder:latest
```

## Discord

### Installation

Install Link:

* None

### OAuth2

2OAuth2 URL Generator
Scopes:

* bot

Bot Permissions:

* Use Slash Commands
* Send Messages
* Manage Events

Integration Type:

* Guild Install

### Bot

Authorization Flow
Public Bot:

* False
