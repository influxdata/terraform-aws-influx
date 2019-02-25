# Chronograf Install Script

Chronograf is a complete web-based interface for the entire InfluxData platform.
This folder contains a script for installing Chronograf and its dependencies.

This script has been tested on the following operating systems:

* Ubuntu 18.04
* Amazon Linux 2

There is a good chance it will work on other flavors of Debian, CentOS, and RHEL as well.

## Quick start

This module depends on [bash-commons](https://github.com/gruntwork-io/bash-commons), so you must install that project
first as documented in its README.

The easiest way to use this module is with the [Gruntwork Installer](https://github.com/gruntwork-io/gruntwork-installer):

```bash
gruntwork-install \
  --module-name "install-chronograf" \
  --repo "https://github.com/gruntwork-io/terraform-aws-influx" \
  --tag "<VERSION>"
```

Checkout the [releases](https://github.com/gruntwork-io/terraform-aws-influx/releases) to find the latest version.

The `install-chronograf` script will install the Chronograf binary as well as its dependencies.

We recommend running the `install-chronograf` script as part of a [Packer](https://www.packer.io/) template to 
create a Chronograf [Amazon Machine Image (AMI)](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/AMIs.html) and then deploy the AMI on an EC2 instance. (See the 
[examples folder](https://github.com/gruntwork-io/terraform-aws-influx/tree/master/examples) for fully-working sample code).


## Command line Arguments

Run `install-chronograf --help` to see all available arguments.

```
Usage: install-chronograf [options]

This script can be used to install Chronograf and its dependencies. This script has been tested with Ubuntu 18.04 and Amazon Linux 2.

Options:

  --version       The version of Chronograf to install. Default: 1.7.8.
  --config-file   Path to a custom configuration file. Default: /tmp/config/chronograf.conf

Example:

  install-chronograf \
    --version 1.7.8 \
    --config-file /tmp/config/chronograf.conf
```

## How it works

The `install-chronograf` script does the following:

1. Installs the Chronograf binary
1. Replaces the default Chronograf config file with your custom config file.
