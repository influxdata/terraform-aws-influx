# Telegraf AMI

This folder shows an example of how to use the 
[install-telegraf](https://github.com/gruntwork-io/terraform-aws-influx/tree/master/modules/install-telegraf)
modules with [Packer](https://www.packer.io/) to create [Amazon Machine 
Images (AMIs)](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/AMIs.html) that have 
[Telegraf](https://www.influxdata.com/time-series-platform/telegraf/), and its dependencies installed on top of:
 
1. Ubuntu 18.04
1. Amazon Linux 2

## Quick start

To build the Telegraf AMI:

1. `git clone` this repo to your computer.
1. Install [Packer](https://www.packer.io/).
1. Configure your AWS credentials using one of the [options supported by the AWS 
   SDK](http://docs.aws.amazon.com/sdk-for-java/v1/developer-guide/credentials.html). Usually, the easiest option is to
   set the `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` environment variables.
1. Update the `variables` section of the `telegraf.json` Packer template to specify the AWS region and Telegraf
   version you wish to use.
1. To build an Ubuntu AMI for Telegraf: `packer build -only=telegraf-ami-ubuntu telegraf.json`.
1. To build an Amazon Linux AMI for Telegraf: `packer build -only=telegraf-ami-amazon-linux telegraf.json`.

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
      "/tmp/terraform-aws-influx/modules/install-telegraf/install-telegraf --version {{user `version`}}"
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
      "/tmp/terraform-aws-influx/modules/install-telegraf/install-telegraf --version {{user `version`}}"
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
