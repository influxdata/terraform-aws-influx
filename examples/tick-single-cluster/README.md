# TICK Single Cluster Example

This folder shows an example of Terraform code that uses the
[influxdb-cluster](https://github.com/gruntwork-io/terraform-aws-influx/tree/master/modules/influxdb-cluster)
module to deploy a [TICK stack](https://www.influxdata.com/time-series-platform/) cluster in [AWS](https://aws.amazon.com/).

This example also deploys a Load Balancer in front of the TICK cluster using the [load-balancer
module](https://github.com/gruntwork-io/terraform-aws-influx/tree/master/modules/load-balancer).

You will need to create an [Amazon Machine Image (AMI)](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/AMIs.html) 
that has all components of the TICK stack installed, which you can do using the [tick-ami 
example](https://github.com/gruntwork-io/terraform-aws-influx/tree/master/examples/tick-ami)). 

To see an example of TICK deployed across separate clusters, see the [tick-multi-cluster
example](https://github.com/gruntwork-io/terraform-aws-influx/blob/master/examples/tick-multi-cluster).

## What resources does this example deploy?

1. A single _all in one server_ behind an ASG where we run 
    [telegraf](/modules/run-telegraf), [influxdb](/modules/run-influxdb),
    [chronograf](/modules/run-chronograf) and [kapacitor](/modules/run-kapacitor)
1. An [Application Load Balancer](https://github.com/gruntwork-io/terraform-aws-load-balancer)

You will need to create [Amazon Machine Images (AMIs)](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/AMIs.html) 
that have all of the TICK-stack components installed. You can do this using: 
- [TICK all-in-one AMI example](/examples/tick-ami)

## Quick start

To deploy a TICK Cluster:

1. `git clone` this repo to your computer.
1. Optional: build a custom TICK AMI. See the
   [tick-ami example](https://github.com/gruntwork-io/terraform-aws-influx/tree/master/examples/tick-ami)
   documentation for instructions. Make sure to note down the ID of the AMI.
1. Install [Terraform](https://www.terraform.io/).
1. Open the `variables.tf` file in this folder, set the environment variables specified at the top of the
   file, and fill in any other variables that don't have a default. If you built a custom AMI, put its ID into the
   `ami_id` variable. If you didn't, this example will use public AMIs that Gruntwork has published, which are fine for
   testing/learning, but not recommended for production use.
1. Run `terraform init`
1. Run `terraform apply`

## Connecting to the cluster

### InfluxDB

Check out [How do you connect to the InfluxDB 
cluster](https://github.com/gruntwork-io/terraform-aws-influx/tree/master/modules/influxdb-cluster#how-do-you-connect-to-the-influxdb-cluster)
documentation.

### Chronograf

Check out [How do you connect to the Chronograf 
cluster](https://github.com/gruntwork-io/terraform-aws-influx/tree/master/modules/chronograf-server#how-do-you-connect-to-the-chronograf-server)
documentation.
