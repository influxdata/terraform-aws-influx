# Kapacitor Install Script

Kapacitor is a Real-Time Streaming Data Processing Engine.
This folder contains a script for installing Kapacitor and its dependencies.

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
  --module-name "install-kapacitor" \
  --repo "https://github.com/gruntwork-io/terraform-aws-influx" \
  --tag "<VERSION>"
```

Checkout the [releases](https://github.com/gruntwork-io/terraform-aws-influx/releases) to find the latest version.

The `install-kapacitor` script will install the Kapacitor binary as well as its dependencies.

We recommend running the `install-kapacitor` script as part of a [Packer](https://www.packer.io/) template to 
create a Kapacitor [Amazon Machine Image (AMI)](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/AMIs.html).
You can then deploy the AMI across an Auto Scaling Group using the [kapacitor-cluster 
module](https://github.com/gruntwork-io/terraform-aws-influx/tree/master/modules/kapacitor-cluster) (see the 
[examples folder](https://github.com/gruntwork-io/terraform-aws-influx/tree/master/examples) for fully-working sample code).


## Command line Arguments

Run `install-kapacitor --help` to see all available arguments.

```
Usage: install-kapacitor [options]

This script can be used to install Kapacitor and its dependencies. This script has been tested with Ubuntu 18.04 and Amazon Linux 2.

Options:

  --version       The version of Kapacitor to install. Default: 1.5.2.
  --config-file   Path to a custom configuration file. Default: /tmp/config/kapacitor.conf

Example:

  install-kapacitor \
    --version 1.5.2 \
    --config-file /tmp/config/kapacitor.conf
```

## How it works

The `install-kapacitor` script does the following:

1. Installs the Kapacitor binary
1. Replaces the default Kapatictor config file with your custom config file.
