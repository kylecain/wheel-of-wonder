#!/bin/sh

PUID=${PUID:-99}
PGID=${PGID:-100}
UMASK=${UMASK:-022}

EXISTING_GROUP=$(getent group "$PGID" | cut -d: -f1)
if [ -n "$EXISTING_GROUP" ]; then
    GROUP_NAME="$EXISTING_GROUP"
else
    GROUP_NAME="pgid_$PGID"
    addgroup -g "$PGID" "$GROUP_NAME"
fi

EXISTING_USER=$(getent passwd "$PUID" | cut -d: -f1)
if [ -n "$EXISTING_USER" ]; then
    USER_NAME="$EXISTING_USER"
else
    USER_NAME="puid_$PUID"
    adduser -D -u "$PUID" -G "$GROUP_NAME" "$USER_NAME"
fi

chown -R "$PUID:$PGID" /app
umask "$UMASK"
exec su "$USER_NAME" -c "$@"
