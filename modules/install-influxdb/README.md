# InfluxDB Install Script

This folder contains a script for installing InfluxDB and its dependencies. Use this script to create an
InfluxDB [Amazon Machine Image (AMI)](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/AMIs.html) that can be
deployed in [AWS](https://aws.amazon.com/) across an Auto Scaling Group using the [influxdb-cluster
module](https://github.com/gruntwork-io/terraform-aws-influx/tree/master/modules/influxdb-cluster).

This script has been tested on the following operating systems:

* Ubuntu 16.04
* Amazon Linux 2

There is a good chance it will work on other flavors of Debian, CentOS, and RHEL as well.

## Quick start

This module depends on [bash-commons](https://github.com/gruntwork-io/bash-commons), so you must install that project
first as documented in its README.

The easiest way to use this module is with the [Gruntwork Installer](https://github.com/gruntwork-io/gruntwork-installer):

```bash
gruntwork-install \
  --module-name "install-influxdb" \
  --repo "https://github.com/gruntwork-io/terraform-aws-influx" \
  --tag "<VERSION>"
```  

Checkout the [releases](https://github.com/gruntwork-io/terraform-aws-influx/releases) to find the latest version.

The `install-influxdb` script will install both InfluxDB meta and data binaries as well as their dependencies.
The [run-influxdb](https://github.com/gruntwork-io/terraform-aws-influx/tree/master/modules/run-influxdb/bin)
script determines whether to startup either as a meta node or a data node.

We recommend running the `install-influxdb` script as part of a [Packer](https://www.packer.io/) template to 
create an InfluxDB [Amazon Machine Image (AMI)](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/AMIs.html).
You can then deploy the AMI across an Auto Scaling Group using the [influxdb-cluster 
module](https://github.com/gruntwork-io/terraform-aws-influx/tree/master/modules/influxdb-cluster) (see the 
[examples folder](https://github.com/gruntwork-io/terraform-aws-influx/tree/master/examples) for fully-working sample code).

## Command line Arguments

Run `install-influxdb --help` to see all available arguments.

```
Usage: install-influxdb [options]

This script can be used to install InfluxDB Enterprise and its dependencies. This script has been tested with Ubuntu 18.04 and Amazon Linux 2.

Options:

  --version           The version of InfluxDB Enterprise to install. Default: 1.6.2.
  --meta-config-file  Path to a templated meta node configuration file. Default: /tmp/config/influxdb-meta.conf
  --data-config-file  Path to a templated data node configuration file. Default: /tmp/config/influxdb.conf

Example:

  install-influxdb \
    --version 1.6.2 \
    --meta-config-file /tmp/config/influxdb-meta.conf \
    --data-config-file /tmp/config/influxdb.conf
```

## How it works

The `install-influxdb` script does the following:

1. Installs the InfluxDB Meta binaries
1. Installs the InfluxDB Data binaries
1. Replaces default config files with specified templated config files.
