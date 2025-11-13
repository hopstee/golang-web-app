#!/bin/sh
# wait-for.sh

HOST=$1
PORT=$2
TIMEOUT=${3:-30}

echo "Waiting for $HOST:$PORT for up to $TIMEOUT seconds..."
i=0
while ! nc -vz "$HOST" "$PORT" >/dev/null 2>&1; do
    i=$((i+1))
    if [ "$i" -ge "$TIMEOUT" ]; then
        echo "Timeout waiting for $HOST:$PORT"
        exit 1
    fi
    sleep 1
done
echo "$HOST:$PORT is available!"
shift 3
exec "$@"