#!/bin/sh
set -e
set -x

echo "DB_HOST=$DB_HOST DB_PORT=$DB_PORT" >&2
echo "REDIS_HOST=$REDIS_HOST REDIS_PORT=$REDIS_PORT" >&2

echo "Waiting for Postgres at $DB_HOST:$DB_PORT" >&2
ls -l /usr/local/bin/wait-for.sh
/usr/local/bin/wait-for.sh "$DB_HOST" "$DB_PORT" 30

echo "Waiting for Redis at $REDIS_HOST:$REDIS_PORT" >&2
/usr/local/bin/wait-for.sh "$REDIS_HOST" "$REDIS_PORT" 30

exec ./bin/app