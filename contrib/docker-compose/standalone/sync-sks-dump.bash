#!/bin/bash

set -eu

docker-compose exec standalone_import-keys_1 /bin/bash -c "cd /import; rsync -avr rsync://rsync.cyberbits.eu/sks/dump ."
