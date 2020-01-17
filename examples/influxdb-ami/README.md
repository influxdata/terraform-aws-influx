# InfluxDB AMI

This folder shows an example of how to use the 
[install-influxdb](https://github.com/gruntwork-io/terraform-aws-influx/tree/master/modules/install-influxdb)
modules with [Packer](https://www.packer.io/) to create [Amazon Machine 
Images (AMIs)](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/AMIs.html) that have 
[InfluxDB Enterprise](https://www.influxdata.com/time-series-platform/influxdb/), and its dependencies installed on top of:
 
1. Ubuntu 18.04
1. Amazon Linux 2

## Quick start

To build the InfluxDB AMI:

1. `git clone` this repo to your computer.
1. Install [Packer](https://www.packer.io/).
1. Configure your AWS credentials using one of the [options supported by the AWS 
   SDK](http://docs.aws.amazon.com/sdk-for-java/v1/developer-guide/credentials.html). Usually, the easiest option is to
   set the `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` environment variables.
1. Update the `variables` section of the `influxdb.json` Packer template to specify the AWS region and InfluxDB
   version you wish to use.
1. To build an Ubuntu AMI for InfluxDB Enterprise: `packer build -only=ubuntu-ami influxdb.json`.
1. To build an Amazon Linux AMI for InfluxDB Enterprise: `packer build -only=amazon-linux-ami influxdb.json`.

When the build finishes, it will output the IDs of the new AMIs. To see how to deploy this AMI, check out the 
[influxdb-cluster-simple](https://github.com/gruntwork-io/terraform-aws-influx/tree/master/examples/influxdb-cluster-simple) and
[tick-multi-cluster](https://github.com/gruntwork-io/terraform-aws-influx/tree/master/examples/tick-multi-cluster) examples.

## Creating your own Packer template for production usage

When creating your own Packer template for production usage, you can copy the example in this folder more or less 
exactly, except for one change: we recommend replacing the `file` provisioner with a call to `git clone` in a `shell` 
provisioner. Instead of:

```json
{
  "provisioners": [{
    "type": "file",
    "source": "{{template_dir}}/../../../terraform-aws-influx",
    "destination": "/tmp"
  },{
    "type": "shell",
    "inline": [
      "/tmp/terraform-aws-influx/modules/install-influxdb/install-influxdb --version {{user `influxdb_version`}}"
    ],
    "pause_before": "30s"
  }]
}
```

Your code should look more like this:

```json
{
  "provisioners": [{
    "type": "shell",
    "inline": [
      "git clone --branch <MODULE_VERSION> https://github.com/gruntwork-io/terraform-aws-influx.git /tmp/terraform-aws-influx",
      "/tmp/terraform-aws-influx/modules/install-influxdb/install-influxdb --version {{user `influxdb_version`}}"
    ],
    "pause_before": "30s"
  }]
}
```

You should replace `<MODULE_VERSION>` in the code above with the version of this module that you want to use (see
the [Releases Page](https://github.com/gruntwork-io/terraform-aws-influx/releases) for all available versions). 
That's because for production usage, you should always use a fixed, known version of this Module, downloaded from the 
official Git repo via `git clone`. On the other hand, when you're just experimenting with the Module, it's OK to use a 
local checkout of the Module, uploaded from your own computer via the `file` provisioner.

## Local testing

The Packer template in this example folder can build not only AMIs, but also Docker images for local testing. This is
convenient for testing out the various scripts in the `modules` folder without having to wait for an AMI to build and
a bunch of EC2 Instances to boot up. See the [local-mocks folder](https://github.com/gruntwork-io/terraform-aws-influx/tree/master/examples/local-mocks) for
instructions.