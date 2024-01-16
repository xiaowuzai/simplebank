#!/bin/sh

set -e

echo "run db migration"
echo $DB_SOURCE
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

echo "start app"
exec "$@" # 运行所有参数