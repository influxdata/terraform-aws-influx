#!/usr/bin/env bash

set -e

# Send the log output from this script to user-data.log, syslog, and the console
# From: https://alestic.com/2010/12/ec2-user-data-output/
exec > >(tee /user-data/mock-user-data.log|logger -t user-data -s 2>/dev/console) 2>&1

function run_telegraf {
  local -r influxdb_url="$1"
  local -r database_name="$2"

  "/opt/telegraf/bin/run-telegraf" \
    --auto-fill "<__INFLUXDB_URL__>=$influxdb_url" \
    --auto-fill "<__DATABASE_NAME__>=$database_name"
}

run_telegraf \
  "${influxdb_url}" \
  "${database_name}"
