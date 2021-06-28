#!/bin/bash

HERE=$(cd $(dirname $0); pwd)
set -eux

[ -e "$HERE/.env" ]

env - $("$HERE/.env") envsubst \
	< "$HERE/hockeypuck/etc/hockeypuck.conf.tmpl" > "$HERE/hockeypuck/etc/hockeypuck.conf"
env - $("$HERE/.env") envsubst \
	< "$HERE/nginx/conf.d/hockeypuck.conf.tmpl" > "$HERE/nginx/conf.d/hockeypuck.conf"
