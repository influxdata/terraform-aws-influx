# Kapacitor Run Script

This folder contains a script for configuring and initializing Kapacitor on an [AWS](https://aws.amazon.com/) server. 
This script has been tested on the following operating systems:

* Ubuntu 18.04
* Amazon Linux 2

There is a good chance it will work on other flavors of Debian, CentOS, and RHEL as well.

## Quick start

This script assumes you installed it, plus all of its dependencies (including Kapacitor itself), using the 
[install-kapacitor module](https://github.com/gruntwork-io/terraform-aws-influx/tree/master/modules/install-kapacitor). 

This will:

1. Fill out the templated configuration file with user supplied values.

1. Start Kapacitor on the machine.

We recommend using the `run-kapacitor` command as part of [User 
Data](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/user-data.html#user-data-shell-scripts), so that it executes
when the EC2 Instance is first booting.

See the [examples folder](https://github.com/gruntwork-io/terraform-aws-influx/tree/master/examples) for 
fully-working sample code.

## Command line Arguments

Run `run-kapacitor --help` to see all available arguments.

```
Usage: run-kapacitor [options]

This script can be used to configure and initialize Kapacitor. This script has been tested with Ubuntu 18.04 and Amazon Linux 2.

Options:

  --auto-fill   Search the Kapacitor config file for KEY and replace it with VALUE. May be repeated.

Example:

  run-kapacitor --auto-fill '<__HOST_NAME__>=0.0.0.0' --auto-fill '<__INFLUXDB_URL__>=http://localhost:8086'
```

## Command line Arguments

Run `run-kapacitor --help` to see all available arguments.

```
Usage: run-kapacitor [options]

This script can be used to configure and initialize Kapacitor. This script has been tested with Ubuntu 18.04 and Amazon Linux 2.

Options:

  --auto-fill   Search the Kapacitor config file for KEY and replace it with VALUE. May be repeated.

Example:

  run-kapacitor --auto-fill '<__HOST_NAME__>0.0.0.0' --auto-fill '<__INFLUXDB_URL__>=http://localhost:8086'
```

## Debugging tips and tricks

Some tips and tricks for debugging issues with your Kapacitor agent:

* Use `systemctl status kapacitor` to see if systemd thinks the Kapacitor process is running.
