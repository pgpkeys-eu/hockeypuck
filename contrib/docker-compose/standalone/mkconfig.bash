#!/bin/bash

# This script will not overwrite any existing config, to protect manual edits.
# To regenerate config files you must remove them first.

# Note that `set -a` causes all variables sourced from `.env` to be implicitly `export`ed.
# This is necessary for envsubst

HERE=$(cd "$(dirname "$0")"; pwd)
set -eua

[ -e "$HERE/.env" ]
. "$HERE/.env"

# Check for migrations
if ! grep -q MIGRATION_3_DONE "$HERE/.env"; then
	cat <<EOF

-----------------------------------------------------------------------
WARNING: Site configuration migration is required before continuing.

Please run 'mksite.bash' to update your site configuration, and
then run this script again.
-----------------------------------------------------------------------

EOF
	exit 1
fi

[ ! -f "$HERE/hockeypuck/etc/hockeypuck.conf" ] &&
	envsubst '$FQDN:$FINGERPRINT:$RELEASE:$POSTGRES_USER:$POSTGRES_PASSWORD' \
	< "$HERE/hockeypuck/etc/hockeypuck.conf.tmpl" > "$HERE/hockeypuck/etc/hockeypuck.conf"
