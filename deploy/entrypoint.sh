#!/usr/bin/env bash

set -e

ln -sf "/usr/share/zoneinfo/${TZ}" /etc/localtime
echo "${TZ}" > /etc/timezone

exec "$@"
