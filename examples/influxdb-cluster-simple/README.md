# InfluxDB Single Cluster Example

The root folder of this repo shows an example of Terraform code that uses the
[influxdb-cluster](https://github.com/gruntwork-io/terraform-aws-influx/tree/master/modules/influxdb-cluster) 
module to deploy a [InfluxDB Enterprise](https://www.influxdata.com/time-series-platform/influxdb/) cluster in [AWS](https://aws.amazon.com/). The cluster 
consists of one Auto Scaling Group (ASG) that runs InfluxDB:

![InfluxDB single-cluster architecture](https://github.com/gruntwork-io/terraform-aws-influx/blob/master/_docs/influxdb-single-cluster-architecture.png?raw=true)

This example also deploys a Load Balancer in front of the InfluxDB cluster using the [load-balancer
module](https://github.com/gruntwork-io/terraform-aws-influx/tree/master/modules/load-balancer).

You will need to create an [Amazon Machine Image (AMI)](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/AMIs.html) 
that has InfluxDB installed, which you can do using the [influxdb-ami 
example](https://github.com/gruntwork-io/terraform-aws-influx/tree/master/examples/influxdb-ami). 

To see an example of InfluxDB deployed across separate clusters, see the [tick-multi-cluster example](https://github.com/gruntwork-io/terraform-aws-influx/tree/master/examples/tick-multi-cluster). For more info on how the InfluxDB cluster works, check out the 
[influxdb-cluster](https://github.com/gruntwork-io/terraform-aws-influx/tree/master/modules/influxdb-cluster) documentation.

## Quick start

To deploy a InfluxDB Cluster:

1. `git clone` this repo to your computer.
1. Optional: build a custom InfluxDB AMI. See the
   [influxdb-ami example](https://github.com/gruntwork-io/terraform-aws-influx/tree/master/examples/influxdb-ami)
   documentation for instructions. Make sure to note down the ID of the AMI.
1. Install [Terraform](https://www.terraform.io/).
1. Open the `variables.tf` file in the root of this repo, set the environment variables specified at the top of the
   file, and fill in any other variables that don't have a default. If you built a custom AMI, put its ID into the
   `ami_id` variable. If you didn't, this example will use public AMIs that Gruntwork has published, which are fine for
   testing/learning, but not recommended for production use.
1. Run `terraform init` in the root folder of this repo.
1. Run `terraform apply` in the root folder of this repo.

## Connecting to the cluster

Check out [How do you connect to the InfluxDB 
cluster](https://github.com/gruntwork-io/terraform-aws-influx/tree/master/modules/influxdb-cluster#how-do-you-connect-to-the-influxdb-cluster)
documentation.
