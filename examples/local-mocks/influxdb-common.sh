#!/bin/bash

# This is a mock version of the influxdb-common script that can run entirely locally without
# depending on external dependencies, such as EC2 Metadata and AWS API calls.
# This allows us to test all the scripts completely locally using Docker.

set -e

function get_node_hostname {
  echo -n "$(hostname)"
}
