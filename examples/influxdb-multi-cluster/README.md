# InfluxDB Multi Cluster Example

This folder shows an example of Terraform code that uses the 
[influxdb-cluster](https://github.com/gruntwork-io/terraform-aws-influx/tree/master/modules/influxdb-cluster) 
module to deploy a [InfluxDB Enterprise](https://www.influxdata.com/time-series-platform/influxdb/) cluster in [AWS](https://aws.amazon.com/). The cluster 
consists of two Auto Scaling Groups (ASGs): one for meta nodes and one for data nodes.

You will need to create an [Amazon Machine Image (AMI)](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/AMIs.html) 
that has InfluxDB installed, which you can do using the [influxdb-ami 
example](https://github.com/gruntwork-io/terraform-aws-influx/tree/master/examples/influxdb-ami)). 

To see an example of InfluxDB deployed in a single cluster, see the [influxdb-single-cluster
example](https://github.com/gruntwork-io/terraform-aws-influx/blob/master/examples/influxdb-single-cluster). For
more info on how the InfluxDB cluster works, check out the 
[influxdb-cluster](https://github.com/gruntwork-io/terraform-aws-influx/tree/master/modules/influxdb-cluster) documentation.

## Quick start

To deploy an InfluxDB Cluster:

1. `git clone` this repo to your computer.
1. Optional: build a custom InfluxDB AMI. See the
   [influxdb-ami example](https://github.com/gruntwork-io/terraform-aws-influx/tree/master/examples/influxdb-ami)
   documentation for instructions. Make sure to note down the ID of the AMI.
1. Install [Terraform](https://www.terraform.io/).
1. Open `variables.tf`, set the environment variables specified at the top of the file, and fill in any other variables that
   don't have a default. If you built a custom AMI, put its ID into the `ami_id` variable. If you didn't, this example
   will use public AMIs that Gruntwork has published, which are fine for testing/learning, but not recommended for
   production use.
1. Run `terraform init`.
1. Run `terraform apply`.

## Connecting to the cluster

Check out [How do you connect to the InfluxDB 
cluster](https://github.com/gruntwork-io/terraform-aws-influxdb/tree/master/modules/influxdb-cluster#how-do-you-connect-to-the-influxdb-cluster)
documentation.
