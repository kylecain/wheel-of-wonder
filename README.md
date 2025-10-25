# Wheel of Wonder

## Go

```zsh
go run cmd/wheel-of-wonder/main.go 
```

## Environment variables

- BOT_TOKEN (required)
    Discord bot token used to authenticate with the Discord API. Treat this as a secret.

- GUILD_ID (optional)
    Discord guild (server) ID (snowflake). Used for guild-scoped command registration or to restrict bot operations to a single server. If GUILD_ID is not provided, commands will be registered globally.

- APPLICATION_ID (required for interactions/command registration)
    Discord application ID (snowflake) for the bot's application. Required when registering slash commands or verifying interactions.

- PUID (optional, default: 99)
    User ID the process should run as inside a container. When mounting host volumes, ensure the UID exists or file permissions may match.

- PGID (optional, default: 100)
    Group ID the process should run as inside a container.

- UMASK (optional, default: 022)
    File-mode creation mask to control default permissions for new files and directories.

## Podman

Local Image

```zsh
podman build -t wheel-of-wonder:local .
podman run \
    -v $(pwd)/data:/app/data \
    --env-file .env \
    wheel-of-wonder:local
```

Remote Image

```zsh
podman pull ghcr.io/kylecain/wheel-of-wonder:latest
podman run \
    -e BOT_TOKEN \
    -v $(pwd)/data:/app/data \
    ghcr.io/kylecain/wheel-of-wonder:latest
```

## GHCR

```zsh
TAG="$(git rev-parse --short HEAD)"
podman build --arch amd64 -t wheel-of-wonder:${TAG} .
echo "$CR_PAT" | podman login ghcr.io -u kylecain --password-stdin
podman tag wheel-of-wonder:${TAG} ghcr.io/kylecain/wheel-of-wonder:${TAG}
podman push ghcr.io/kylecain/wheel-of-wonder:${TAG}
podman tag ghcr.io/kylecain/wheel-of-wonder:${TAG} ghcr.io/kylecain/wheel-of-wonder:latest
podman push ghcr.io/kylecain/wheel-of-wonder:latest
```

## Discord

### Installation

Install Link:

- None

### OAuth2

2OAuth2 URL Generator
Scopes:

- bot

Bot Permissions:

- Use Slash Commands
- Send Messages
- Manage Events

Integration Type:

- Guild Install

### Bot

Authorization Flow
Public Bot:

- False
