# Telegraf Install Script

Telegraf is an agent for collecting metrics from a variety of sources
and writing them to InfluxDB or other outputs.
This folder contains a script for installing Telegraf and its dependencies.

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
  --module-name "install-telegraf" \
  --repo "https://github.com/gruntwork-io/terraform-aws-influx" \
  --tag "<VERSION>"
```  

Checkout the [releases](https://github.com/gruntwork-io/terraform-aws-influx/releases) to find the latest version.

The `install-telegraf` script will install the Telegraf binary as well as its dependencies.

We recommend including and executing this script in your Application Server's
 [Packer](https://www.packer.io/) template that creates an [Amazon Machine Image (AMI)](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/AMIs.html), so that Telegraf is present to collect
 application and machine metrics.


## Command line Arguments

Run `install-telegraf --help` to see all available arguments.

```
Usage: install-telegraf [options]

This script can be used to install Telegraf and its dependencies. This script has been tested with Ubuntu 18.04 and Amazon Linux 2.

Options:

  --version       The version of Telegraf to install. Default: 1.9.3.
  --config-file   Path to a templated configuration file. Default: /tmp/config/telegraf.conf

Example:

  install-telegraf \
    --version 1.9.3 \
    --config-file /tmp/config/telegraf.conf
```

## How it works

The `install-telegraf` script does the following:

1. Installs the Telegraf binary
1. Replaces the default Telegraf config file with your custom config file.
