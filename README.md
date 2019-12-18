[![Maintained by Gruntwork.io](https://img.shields.io/badge/maintained%20by-gruntwork.io-%235849a6.svg)](https://gruntwork.io/?ref=repo_aws_influx)

# TICK Stack AWS Module

This repo contains the **official** module for deploying the [TICK stack](https://www.influxdata.com/time-series-platform/) on [AWS](https://aws.amazon.com/)
using [Terraform](https://www.terraform.io/) and [Packer](https://www.packer.io/).

![TICK multi-cluster architecture](https://github.com/gruntwork-io/terraform-aws-influx/blob/master/_docs/tick-multi-cluster-architecture.png?raw=true)

## Quick start

If you want to quickly spin up an InfluxDB cluster, you can run the simple example that is in the root of this repo.
Check out [influxdb-cluster-simple example
documentation](https://github.com/gruntwork-io/terraform-aws-influx/blob/master/examples/influxdb-cluster-simple)
for instructions.

## What's in this repo

This repo has the following folder structure:

- [root](https://github.com/gruntwork-io/terraform-aws-influx/tree/master): The root folder contains an example
  of how to deploy InfluxDB as a single-cluster. See
  [influxdb-cluster-simple](https://github.com/gruntwork-io/terraform-aws-influx/blob/master/examples/influxdb-cluster-simple)
  for the documentation.
- [modules](https://github.com/gruntwork-io/terraform-aws-influx/tree/master/modules): This folder contains the
  main implementation code for this Module, broken down into multiple standalone submodules.
- [examples](https://github.com/gruntwork-io/terraform-aws-influx/tree/master/examples): This folder contains
  examples of how to use the submodules.
- [test](https://github.com/gruntwork-io/terraform-aws-influx/tree/master/test): Automated tests for the submodules
  and examples.

## How to use this repo

The general idea is to:

- ### Telegraf

  1. Use the scripts in the
     [install-telegraf](https://github.com/gruntwork-io/terraform-aws-influx/tree/master/modules/install-telegraf)
     modules to create an AMI with Telegraf installed, this AMI will generally be for the Application server.

  1. Configure each application server to execute the
     [run-telegraf](https://github.com/gruntwork-io/terraform-aws-influx/tree/master/modules/run-telegraf)
     script during boot.

- ### InfluxDB

  1. Use the scripts in the
     [install-influxdb](https://github.com/gruntwork-io/terraform-aws-influx/tree/master/modules/install-influxdb)
     modules to create an AMI with InfluxDB Enterprise installed.

  1. Deploy the AMI across one or more Auto Scaling Groups (ASG) using the [influxdb-cluster
     module](https://github.com/gruntwork-io/terraform-aws-influx/tree/master/modules/influxdb-cluster).

  1. Configure each server in the ASGs to execute the
     [run-influxdb](https://github.com/gruntwork-io/terraform-aws-influx/tree/master/modules/run-influxdb)
     script during boot.

  1. Deploy a load balancer in front of the data node ASG.

- ### Chronograf

  1. Use the scripts in the
     [install-chronograf](https://github.com/gruntwork-io/terraform-aws-influx/tree/master/modules/install-chronograf)
     modules to create an AMI with Chronograf installed.

  1. Deploy the AMI in a single Auto Scaling Group (ASG) using the [chronograf-server
     module](https://github.com/gruntwork-io/terraform-aws-influx/tree/master/modules/chronograf-server).

  1. Configure the server to execute the
     [run-chronograf](https://github.com/gruntwork-io/terraform-aws-influx/tree/master/modules/run-chronograf)
     script during boot.

  1. Deploy a load balancer in front of the ASG.

- ### Kapacitor

  1. Use the scripts in the
     [install-kapacitor](https://github.com/gruntwork-io/terraform-aws-influx/tree/master/modules/install-kapacitor)
     modules to create an AMI with Kapacitor installed.

  1. Deploy the AMI in a single Auto Scaling Group (ASG) using the [kapacitor-server
     module](https://github.com/gruntwork-io/terraform-aws-influx/tree/master/modules/kapacitor-server).

  1. Configure the server to execute the
     [run-kapacitor](https://github.com/gruntwork-io/terraform-aws-influx/tree/master/modules/run-kapacitor)
     script during boot.

  1. Deploy a load balancer in front of the ASG.

See the [examples folder](https://github.com/gruntwork-io/terraform-aws-influx/tree/master/examples) for working
sample code.

## What's a Module?

A Module is a canonical, reusable, best-practices definition for how to run a single piece of infrastructure, such
as a database or server cluster. Each Module is written using a combination of [Terraform](https://www.terraform.io/)
and scripts (mostly bash) and include automated tests, documentation, and examples. It is maintained both by the open
source community and companies that provide commercial support.

Instead of figuring out the details of how to run a piece of infrastructure from scratch, you can reuse
existing code that has been proven in production. And instead of maintaining all that infrastructure code yourself,
you can leverage the work of the Module community to pick up infrastructure improvements through
a version number bump.

## Who maintains this Module?

This Module is maintained by [Gruntwork](http://www.gruntwork.io/). If you're looking for help or commercial
support, send an email to [modules@gruntwork.io](mailto:modules@gruntwork.io?Subject=InfluxDB%20for%20AWS%20Module).
Gruntwork can help with:

- Setup, customization, and support for this Module.
- Modules for other types of infrastructure, such as VPCs, Docker clusters, databases, and continuous integration.
- Modules that meet compliance requirements, such as HIPAA.
- Consulting & Training on AWS, Terraform, and DevOps.

## How do I contribute to this Module?

Contributions are very welcome! Check out the
[Contribution Guidelines](https://github.com/gruntwork-io/terraform-aws-influx/tree/master/CONTRIBUTING.md) for instructions.

## How is this Module versioned?

This Module follows the principles of [Semantic Versioning](http://semver.org/). You can find each new release,
along with the changelog, in the [Releases Page](../../releases).

During initial development, the major version will be 0 (e.g., `0.x.y`), which indicates the code does not yet have a
stable API. Once we hit `1.0.0`, we will make every effort to maintain a backwards compatible API and use the MAJOR,
MINOR, and PATCH versions on each release to indicate any incompatibilities.

## License

This code is released under the Apache 2.0 License. Please see
[LICENSE](https://github.com/gruntwork-io/terraform-aws-influx/tree/master/LICENSE) and
[NOTICE](https://github.com/gruntwork-io/terraform-aws-influx/tree/master/NOTICE) for more details.

Copyright &copy; 2018 Gruntwork, Inc.
