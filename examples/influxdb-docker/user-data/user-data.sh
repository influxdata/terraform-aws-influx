#!/usr/bin/env bash

set -e

# Send the log output from this script to user-data.log, syslog, and the console
# From: https://alestic.com/2010/12/ec2-user-data-output/
exec > >(tee /user-data/mock-user-data.log|logger -t user-data -s 2>/dev/console) 2>&1

function run {
  local -r asg_name="$1"
  local -r aws_region="$2"
  local -r license_key="$3"
  local -r shared_secret="$4"
  local -r hostname="$(hostname)"

  "/opt/influxdb/bin/run-influxdb" \
    --node-type "meta" \
    --meta-asg-name "$asg_name" \
    --data-asg-name "$asg_name" \
    --region "$aws_region" \
    --auto-fill "<__HOST_NAME__>=$hostname" \
    --auto-fill "<__LICENSE_KEY__>=$license_key"

  "/opt/influxdb/bin/run-influxdb" \
    --node-type "data" \
    --meta-asg-name "$asg_name" \
    --data-asg-name "$asg_name" \
    --region "$aws_region" \
    --auto-fill "<__HOST_NAME__>=$hostname" \
    --auto-fill "<__LICENSE_KEY__>=$license_key" \
    --auto-fill "<__SHARED_SECRET__>=$shared_secret"
}

run \
  "${cluster_asg_name}" \
  "${mock_aws_region}" \
  "${license_key}" \
  "${shared_secret}"
