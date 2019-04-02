# TICK AMI

This folder shows an example of how to use [Packer](https://www.packer.io/) to create [Amazon Machine 
Images (AMIs)](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/AMIs.html) that have all of the TICK components
(InfluxDB, Chronograf, Kapacitor and Telegraf) installed on top of:
 
1. Ubuntu 18.04
1. Amazon Linux 2

## Quick start

To build the TICK AMI:

1. `git clone` this repo to your computer.
1. Install [Packer](https://www.packer.io/).
1. Configure your AWS credentials using one of the [options supported by the AWS 
   SDK](http://docs.aws.amazon.com/sdk-for-java/v1/developer-guide/credentials.html). Usually, the easiest option is to
   set the `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` environment variables.
1. Update the `variables` section of the `tick.json` Packer template to specify the AWS region and Elasticsearch
   version you wish to use.
1. To build the Ubuntu AMI: `packer build -only=tick-ami-ubuntu tick.json`.
1. To build the Amazon Linux 2 AMI: `packer build -only=tick-ami-amazon-linux tick.json`.