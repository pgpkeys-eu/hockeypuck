#!/bin/bash

set -eu

docker-compose up -d import-keys
docker-compose exec import-keys /bin/bash -c "cd /import; rsync -avr rsync://rsync.cyberbits.eu/sks/dump ."
