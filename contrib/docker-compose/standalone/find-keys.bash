#!/bin/bash

# Find keys in the Hockeypuck postgres database by keyword search

set -euo pipefail

if [[ ! ${1:-} ]]; then
    cat <<EOF
Usage: $0 [options] SEARCH_STRING

Where options are:

-r  SEARCH_STRING is a regex (searches against first userid only)
-t  SEARCH_STRING is a formatted tsquery

Otherwise SEARCH_STRING is a search-engine style webquery.

EOF
    exit 1
fi

# Uncomment and edit one of the below for your postgres installation
# for docker-compose/standalone default configuration
SQLCMD="docker exec -i standalone_postgres_1 psql hkp -U hkp"
# for docker-compose/dev default configuration
#SQLCMD="docker exec -i hockeypuck_postgres_1 psql hkp -U docker"
# for non-docker postgres, e.g.
#SQLCMD="psql hkp -U hkp"

if [[ $1 == -r ]]; then
    $SQLCMD -t <<<"select reverse(rfingerprint), keywords from keys where doc->'userIDs'->0->>'keywords' ~ '$2';"
elif [[ $1 == -t ]]; then
    $SQLCMD -t <<<"select reverse(rfingerprint), keywords from keys, to_tsquery($2) query where query @@ keywords;"
else
    $SQLCMD -t <<<"select reverse(rfingerprint), keywords from keys, websearch_to_tsquery('$1') query where query @@ keywords;"
fi
