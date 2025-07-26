# Wheel of Wonder

## Go

```zsh
go run cmd/wheel-of-wonder/main.go 
```

## Podman

```zsh
podman build -t wheel-of-wonder .
podman run -v $(pwd)/data:/app/data --env-file .env wheel-of-wonder
```

## GHCR

```zsh
podman build -t wheel-of-wonder:latest .
echo $CR_PAT | podman login ghcr.io -u kylecain --password-stdin
podman tag wheel-of-wonder:latest ghcr.io/kylecain/wheel-of-wonder:latest
podman push ghcr.io/kylecain/wheel-of-wonder:latest
```
