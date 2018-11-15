# InfluxDB Cluster

This folder contains a [Terraform](https://www.terraform.io/) module to deploy an [InfluxDB Enterprise](
https://www.influxdata.com/time-series-platform/influxdb/) cluster in [AWS](https://aws.amazon.com/) on top of an Auto Scaling Group. 
The idea is to create an [Amazon Machine Image (AMI)](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/AMIs.html)
that has InfluxDB meta and data binaries installed using the [install-influxdb](
https://github.com/gruntwork-io/terraform-aws-influx/tree/master/modules/install-influxdb) module.

## How do you use this module?

This folder defines a [Terraform module](https://www.terraform.io/docs/modules/usage.html), which you can use in your
code by adding a `module` configuration and setting its `source` parameter to URL of this folder:

```hcl
module "influxdb_data_cluster" {
  # TODO: replace <VERSION> with the latest version from the releases page: https://github.com/gruntwork-io/terraform-aws-influx/releases
  source = "github.com/gruntwork-io/terraform-aws-influx//modules/influxdb-cluster?ref=<VERSION>"

  # Specify the ID of the InfluxDB AMI. You should build this using the scripts in the install-influxdb module.
  ami_id = "ami-abcd1234"
  
  # Configure and start InfluxDB during boot. 
  user_data = <<-EOF
              #!/bin/bash
              sudo systemctl start influxdb
              EOF
  
  # ... See variables.tf for the other parameters you must define for the influxdb-cluster module
}
```

Note the following parameters:

* `source`: Use this parameter to specify the URL of the influxdb-cluster module. The double slash (`//`) is 
  intentional and required. Terraform uses it to specify subfolders within a Git repo (see [module 
  sources](https://www.terraform.io/docs/modules/sources.html)). The `ref` parameter specifies a specific Git tag in 
  this repo. That way, instead of using the latest version of this module from the `master` branch, which 
  will change every time you run Terraform, you're using a fixed version of the repo.

* `ami_id`: Use this parameter to specify the ID of an InfluxDB [Amazon Machine Image 
  (AMI)](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/AMIs.html) to deploy on each server in the cluster. You
  should install InfluxDB the scripts in the 
  [install-influxdb](https://github.com/gruntwork-io/terraform-aws-influx/tree/master/modules/install-influxdb) module.
  
* `user_data`: Use this parameter to specify a [User 
  Data](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/user-data.html#user-data-shell-scripts) script that each
  server will run during boot. This is where you can use the 
  [run-influxdb](https://github.com/gruntwork-io/terraform-aws-influx/tree/master/modules/run-influxdb) 
  script to configure and run InfluxDB as either a meta node or a data node. 

You can find the other parameters in [variables.tf](variables.tf).

Check out the [examples folder](https://github.com/gruntwork-io/terraform-aws-influx/tree/master/examples) for 
fully-working sample code.

## Why InfluxDB enterprise?

Unlike the Enterprise edition which distributes the entire InfluxDB deployment across multiple meta and data nodes in a cluster,
the OSS edition is a single binary that can simply be installed on a single instance. This makes a robust module like this one
quite unnecessary for single instance InfluxDB OSS deployments.

## How do you connect to the InfluxDB cluster?

While InfluxDB doesn't use a load balancer for intra-cluster communication, one can be optionally setup to communicate
with a data cluster from the outside world. An [Application Load Balancer](
http://docs.aws.amazon.com/elasticloadbalancing/latest/application/introduction.html) is deployed in front of the cluster
and you can use the [InfluxDB CLI](https://docs.influxdata.com/influxdb/v1.6/tools/shell/)
to connect to the load balancer on port `8086`. Alternatively, you can use the [InfluxDB SDK](
https://docs.influxdata.com/influxdb/v1.6/tools/api_client_libraries/) for your programming language of choice.

## What's included in this module?

This module creates the following:

* [Auto Scaling Group](#auto-scaling-group)
* [EBS Volumes](#ebs-volumes)
* [Security Group](#security-group)
* [IAM Role and Permissions](#iam-role-and-permissions)

### Auto Scaling Group

This module runs InfluxDB on top of an [Auto Scaling Group (ASG)](https://aws.amazon.com/autoscaling/). Typically, you
should run the ASG with multiple Instances spread across multiple [Availability 
Zones](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/using-regions-availability-zones.html). Each of the EC2
Instances should be running an AMI that has InfluxDB installed via the 
[install-influxdb](https://github.com/gruntwork-io/terraform-aws-influx/tree/master/modules/install-influxdb)
module. You pass in the ID of the AMI to run using the `ami_id` input parameter.

The [run-influxdb](https://github.com/gruntwork-io/terraform-aws-influx/tree/master/modules/run-influxdb) script ensures that
a new node added by an ASG scale-up event automatically joins the existing cluster. However, you're responsible for de-registering
any dead nodes (terminated EC2 instance) from the cluster when the ASG scales down. See the [cluster commands 
guide](https://docs.influxdata.com/enterprise_influxdb/v1.5/features/cluster-commands/) for instructions.

### EBS Volumes

This module can optionally create an [EBS volume](https://aws.amazon.com/ebs/) for each EC2 Instance in the ASG. You 
can use these volume to store InfluxDB data and are mandotory for InfluxDB data nodes.

We recommend a single EBS volume for the meta nodes for storing its Raft database. And one for data nodes,
to store the `/data`, `/meta`, `/wal`, and `/hh` direcotries.

### Backup and Replication

Setting up an adequate backup and recovery mechanism fo your cluster is hugely important. InfluxDB when used with some
other components of the [TICK stack](https://www.influxdata.com/time-series-platform/) has a very capable multi-datacenter
and cross-cluster replication system. You can find more information of InfluxDB disaster recovery [here](https://www.influxdata.com/blog/multiple-data-center-replication-influxdb/).

### Security Group

Each EC2 Instance in the ASG has a Security Group that allows minimal connectivity:

* All outbound requests
* Inbound SSH access from the CIDR blocks and security groups you specify

The Security Group ID is exported as an output variable which you can use with the 
[influxdb-security-group-rules](https://github.com/gruntwork-io/terraform-aws-influx/tree/master/modules/influxdb-security-group-rules)
module to open up all the ports necessary for InfluxDB.

### IAM Role and Permissions

Each EC2 Instance in the ASG has an [IAM Role](http://docs.aws.amazon.com/IAM/latest/UserGuide/id_roles.html) attached. 
The IAM Role ARN and ID are exported as output variables if you need to add additional permissions.

### How do you roll out updates?

This module currently doesn't support updating the version of InfluxDB installed across the cluster. You can however follow
this [guide](https://docs.influxdata.com/enterprise_influxdb/v1.6/administration/upgrading/) to manually perform the upgrade.

### Dedicated instances

If you wish to use dedicated instances, you can set the `tenancy` parameter to `"dedicated"` in this module. 

### Encryption

This module does not currently support specifying encryption information. The official [documentation](
https://docs.influxdata.com/influxdb/v1.6/administration/https_setup/) contains a guide for enabling SSL.

### Security groups

This module attaches a security group to each EC2 Instance that allows inbound requests as follows:

* **SSH**: For the SSH port (default: 22), you can use the `allowed_ssh_cidr_blocks` parameter to control the list of   
  [CIDR blocks](https://en.wikipedia.org/wiki/Classless_Inter-Domain_Routing) that will be allowed access. You can use 
  the `allowed_inbound_ssh_security_group_ids` parameter to control the list of source Security Groups that will be 
  allowed access.
  
The ID of the security group is exported as an output variable, which you can use with the 
[influxdb-security-group-rules](https://github.com/gruntwork-io/terraform-aws-influx/tree/master/modules/influxdb-security-group-rules)
module to open up all the ports necessary for InfluxDB.

### SSH access

You can associate an [EC2 Key Pair](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-key-pairs.html) with each
of the EC2 Instances in this cluster by specifying the Key Pair's name in the `ssh_key_name` variable. If you don't
want to associate a Key Pair with these servers, set `ssh_key_name` to an empty string.
