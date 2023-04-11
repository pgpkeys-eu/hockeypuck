#!/bin/bash

# Find keys in the Hockeypuck postgres database by keyword search

set -euo pipefail

# Uncomment and edit one of the below for your postgres installation
# for docker-compose/standalone default configuration
SQLCMD="docker exec -i standalone_postgres_1 psql hkp -U hkp"
# for docker-compose/dev default configuration
#SQLCMD="docker exec -i hockeypuck_postgres_1 psql hkp -U docker"
# for non-docker postgres, e.g.
#SQLCMD="psql hkp -U hkp"

usage() {
    cat <<EOF
Usage: $0 [options] SEARCH

If SEARCH is "-", then search parameters of the appropriate type are read from STDIN, one per line.

Command line options are:

-r  each search parameter is a regex (searches against first userid only)
-t  each search parameter is a SQL tsquery

Otherwise each search parameter is a search-engine style webquery.

EOF
    exit 1
}

s_userid_regex() {
    $SQLCMD -t -c "select reverse(rfingerprint), keywords from keys where doc->'userIDs'->0->>'keywords' ~ '$1';"
}

s_keywords_tsquery() {
    $SQLCMD -t -c "select reverse(rfingerprint), keywords from keys, to_tsquery($1) query where query @@ keywords;"
}

s_keywords_websearch() {
    $SQLCMD -t -c "select reverse(rfingerprint), keywords from keys, websearch_to_tsquery('$1') query where query @@ keywords;"
}

[[ ${1:-} ]] || usage

if [[ $1 == -r ]]; then
    shift
    [[ ${1:-} ]] || usage
    COMMAND=s_userid_regex
elif [[ $1 == -t ]]; then
    shift
    [[ ${1:-} ]] || usage
    COMMAND=s_keywords_tsquery
else
    COMMAND=s_keywords_websearch
fi

if [[ $1 == "-" ]]; then
    while read -r pattern ; do
        [[ $pattern && "${pattern:0:1}" != "#" ]] || continue
        echo "# $pattern"
        $COMMAND "$pattern"
    done
else
    $COMMAND "$1"
fi
