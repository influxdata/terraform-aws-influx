#!/usr/bin/env bash

set -e

# Send the log output from this script to user-data.log, syslog, and the console
# From: https://alestic.com/2010/12/ec2-user-data-output/
exec > >(tee /user-data/mock-user-data.log|logger -t user-data -s 2>/dev/console) 2>&1

source "/opt/influxdb-commons/mount-volume.sh"

function mount_volumes {
  local -r volume_device_name="$1"
  local -r volume_mount_point="$2"
  local -r volume_owner="$3"

  mount_volume "$volume_device_name" "$volume_mount_point" "$volume_owner"
}

function run_kapacitor {
  local -r hostname="$1"
  local -r influxdb_url="$2"
  local -r volume_device_name="$3"
  local -r volume_mount_point="$4"
  local -r volume_owner="$5"

  mount_volumes "$volume_device_name" "$volume_mount_point" "$volume_owner"

  "/opt/kapacitor/bin/run-kapacitor" \
    --auto-fill "<__HOST_NAME__>=$hostname" \
    --auto-fill "<__STORAGE_DIR__>=$volume_mount_point" \
    --auto-fill "<__INFLUXDB_URL__>=$influxdb_url"
}

run_kapacitor \
  "${hostname}" \
  "${influxdb_url}" \
  "${volume_device_name}" \
  "${volume_mount_point}" \
  "${volume_owner}"
