#!/bin/sh

PUID=${PUID:-1000}
PGID=${PGID:-1000}

EXISTING_GROUP=$(getent group "$PGID" | cut -d: -f1)
if [ -n "$EXISTING_GROUP" ]; then
    GROUP_NAME="$EXISTING_GROUP"
else
    addgroup -g "$PGID" appgroup
    GROUP_NAME="appgroup"
fi

EXISTING_USER=$(getent passwd "$PUID" | cut -d: -f1)
if [ -n "$EXISTING_USER" ]; then
    USER_NAME="$EXISTING_USER"
else
    adduser -D -u "$PUID" -G "$GROUP_NAME" appuser
    USER_NAME="appuser"
fi

chown -R "$PUID:$PGID" /app

exec su "$USER_NAME" -c "$@"
