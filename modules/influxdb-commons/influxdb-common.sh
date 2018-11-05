#!/bin/bash

set -e

# This method is used to retrieve the hostname of the EC2 instance
function get_node_hostname {
  local -r hostname="$(curl http://169.254.169.254/latest/meta-data/hostname)"
  echo -n "$hostname"
}
