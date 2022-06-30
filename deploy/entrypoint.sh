#!/usr/bin/env bash
#
set -e

if [ ! -z "$TZ" ];then
	cp "/usr/share/zoneinfo/$TZ" /etc/localtime
	echo "$TZ" > /etc/timezone
else
	cp "/usr/share/zoneinfo/Asia/Shanghai" /etc/localtime
	echo "Asia/Shanghai" > /etc/timezone
fi

exec "$@"