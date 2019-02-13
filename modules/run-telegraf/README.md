# Telegraf Run Script

This folder contains a script for configuring and initializing Telegraf on an [AWS](https://aws.amazon.com/) server. 
This script has been tested on the following operating systems:

* Ubuntu 18.04
* Amazon Linux 2

There is a good chance it will work on other flavors of Debian, CentOS, and RHEL as well.

## Quick start

This script assumes you installed it, plus all of its dependencies (including Telegraf itself), using the 
[install-telegraf module](https://github.com/gruntwork-io/terraform-aws-influx/tree/master/modules/install-telegraf). 

This will:

1. Fill out the templated configuration file with user supplied values. For example:

    ```conf
    ...
    [[outputs.influxdb]]
    ## The full HTTP or UDP URL for your InfluxDB instance.
    urls = ["<__INFLUXDB_URL__>"]

    ## The target database for metrics; will be created as needed.
    database = ["<__DATABASE_NAME__>"]
    ...
    ```

1. Start Telegraf on the machine.

We recommend using the `run-telegraf` command as part of [User 
Data](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/user-data.html#user-data-shell-scripts), so that it executes
when the EC2 Instance is first booting.

See the [examples folder](https://github.com/gruntwork-io/terraform-aws-influx/tree/master/examples) for 
fully-working sample code.

## Command line Arguments

Run `run-telegraf --help` to see all available arguments.

```
Usage: run-telegraf [options]

This script can be used to configure and initialize Telegraf. This script has been tested with Ubuntu 18.04 and Amazon Linux 2.

Options:

  --auto-fill   Search the Telegraf config file for KEY and replace it with VALUE. May be repeated.

Example:

  run-telegraf --auto-fill '<__INFLUXDB_URL__>=http://localhost:8086' --auto-fill '<__DATABASE_NAME__>=telegraf'
```

## Debugging tips and tricks

Some tips and tricks for debugging issues with your Telegraf agent:

* Use `systemctl status telegraf` to see if systemd thinks the Telegraf process is running.
