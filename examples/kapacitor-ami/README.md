# Kapacitor AMI

This folder shows an example of how to use the 
[install-kapacitor](https://github.com/gruntwork-io/terraform-aws-influx/tree/master/modules/install-kapacitor)
modules with [Packer](https://www.packer.io/) to create [Amazon Machine 
Images (AMIs)](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/AMIs.html) that have 
[Kapacitor](https://www.influxdata.com/time-series-platform/kapacitor/), and its dependencies installed on top of:
 
1. Ubuntu 18.04
1. Amazon Linux 2

Kapacitor is usually installed as a collection agent on your application's EC2 instance(s).
This Kapacitor AMI is only useful for starting up Kapacitor to pull remote data or accept data from a remote service,
e.g. connecting to a queue like Kafka/PubSub or using it as a scraper to pull prometheus metrics.

## Quick start

To build the Kapacitor AMI:

1. `git clone` this repo to your computer.
1. Install [Packer](https://www.packer.io/).
1. Configure your AWS credentials using one of the [options supported by the AWS 
   SDK](http://docs.aws.amazon.com/sdk-for-java/v1/developer-guide/credentials.html). Usually, the easiest option is to
   set the `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` environment variables.
1. Update the `variables` section of the `kapacitor.json` Packer template to specify the AWS region and Kapacitor
   version you wish to use.
1. To build an Ubuntu AMI for Kapacitor: `packer build -only=kapacitor-ami-ubuntu kapacitor.json`.
1. To build an Amazon Linux AMI for Kapacitor: `packer build -only=kapacitor-ami-amazon-linux kapacitor.json`.

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
      "/tmp/terraform-aws-influx/modules/install-kapacitor/install-kapacitor --version {{user `version`}}"
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
      "/tmp/terraform-aws-influx/modules/install-kapacitor/install-kapacitor --version {{user `version`}}"
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
