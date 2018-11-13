#!/bin/bash

set -e

# This method is used to retrieve the hostname of the EC2 instance
function get_node_hostname {
  echo -n "$(curl --location --silent --fail --show-error http://169.254.169.254/latest/meta-data/hostname)"
}
