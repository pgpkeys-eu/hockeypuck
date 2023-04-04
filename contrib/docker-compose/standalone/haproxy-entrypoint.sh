#!/bin/sh

haproxy "$@" &

trap exit TERM
while :; do
  sleep 1800
  ps | awk '{ if($4 == "haproxy") {print $1} } ' | xargs kill -HUP
done
