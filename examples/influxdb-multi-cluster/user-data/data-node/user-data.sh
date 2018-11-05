#!/usr/bin/env bash

set -e

source "/opt/influxdb-commons/influxdb-common.sh"
source "/opt/influxdb-commons/mount-volume.sh"

function mount_volumes {
  local -r data_volume_device_name="$1"
  local -r data_volume_mount_point="$2"
  local -r volume_owner="$3"

  echo "Mounting EBS Volume for meta, data, wal and hh directories"
  mount_volume "$data_volume_device_name" "$data_volume_mount_point" "$volume_owner"
}

function run {
  local -r meta_asg_name="$1"
  local -r data_asg_name="$2"
  local -r aws_region="$3"
  local -r license_key="$4"
  local -r shared_secret="$5"
  local -r data_volume_device_name="$6"
  local -r data_volume_mount_point="$7"
  local -r volume_owner="$8"
  local -r hostname=$(get_node_hostname)

  mount_volumes "$data_volume_device_name" "$data_volume_mount_point" "$volume_owner"

  local -r meta_dir="$data_volume_mount_point/var/lib/influxdb/meta"
  local -r data_dir="$data_volume_mount_point/var/lib/influxdb/data"
  local -r wal_dir="$data_volume_mount_point/var/lib/influxdb/wal"
  local -r hh_dir="$data_volume_mount_point/var/lib/influxdb/hh"

  "/opt/influxdb/bin/run-influxdb" \
    --hostname "$hostname" \
    --node-type "data" \
    --meta-asg-name "$meta_asg_name" \
    --data-asg-name "$data_asg_name" \
    --region "$aws_region" \
    --auto-fill "<__HOST_NAME__>=$hostname" \
    --auto-fill "<__LICENSE_KEY__>=$license_key" \
    --auto-fill "<__SHARED_SECRET__>=$shared_secret" \
    --auto-fill "<__META_DIR__>=$meta_dir" \
    --auto-fill "<__DATA_DIR__>=$data_dir" \
    --auto-fill "<__WAL_DIR__>=$wal_dir" \
    --auto-fill "<__HH_DIR__>=$hh_dir"
}

# The variables below are filled in via Terraform interpolation
run \
  "${meta_cluster_asg_name}" \
  "${data_cluster_asg_name}" \
  "${aws_region}" \
  "${license_key}" \
  "${shared_secret}" \
  "${data_volume_device_name}" \
  "${data_volume_mount_point}" \
  "${volume_owner}"
