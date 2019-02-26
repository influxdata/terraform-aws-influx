#!/usr/bin/env bash

set -e

# Send the log output from this script to user-data.log, syslog, and the console
# From: https://alestic.com/2010/12/ec2-user-data-output/
exec > >(tee /user-data/mock-user-data.log|logger -t user-data -s 2>/dev/console) 2>&1

function run_chronograf {
  local -r host="$1"
  local -r port="$2"

  "/opt/chronograf/bin/run-chronograf" \
    --auto-fill "<__HOST__>=$host" \
    --auto-fill "<__PORT__>=$port"
}

run_chronograf \
  "${host}" \
  "${port}"
