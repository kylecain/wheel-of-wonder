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
