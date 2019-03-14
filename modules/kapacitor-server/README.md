# Kapacitor Server

This folder contains a [Terraform](https://www.terraform.io/) module to deploy an [Kapacitor Enterprise](
https://www.influxdata.com/time-series-platform/kapacitor/) cluster in [AWS](https://aws.amazon.com/) on top of an Auto Scaling Group. 
The idea is to create an [Amazon Machine Image (AMI)](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/AMIs.html)
that has Kapacitor binaries installed using the [install-kapacitor](
https://github.com/gruntwork-io/terraform-aws-influx/tree/master/modules/install-kapacitor) module.

## How do you use this module?

This folder defines a [Terraform module](https://www.terraform.io/docs/modules/usage.html), which you can use in your
code by adding a `module` configuration and setting its `source` parameter to URL of this folder:

```hcl
module "kapacitor_cluster" {
  # TODO: replace <VERSION> with the latest version from the releases page: https://github.com/gruntwork-io/terraform-aws-influx/releases
  source = "github.com/gruntwork-io/terraform-aws-influx//modules/kapacitor-cluster?ref=<VERSION>"

  # Specify the ID of the Kapacitor AMI. You should build this using the scripts in the install-kapacitor module.
  ami_id = "ami-abcd1234"
  
  # Configure and start Kapacitor during boot. 
  user_data = <<-EOF
              #!/bin/bash
              sudo systemctl start kapacitor
              EOF
  
  # ... See variables.tf for the other parameters you must define for the kapacitor-cluster module
}
```

Note the following parameters:

* `source`: Use this parameter to specify the URL of the kapacitor-cluster module. The double slash (`//`) is 
  intentional and required. Terraform uses it to specify subfolders within a Git repo (see [module 
  sources](https://www.terraform.io/docs/modules/sources.html)). The `ref` parameter specifies a specific Git tag in 
  this repo. That way, instead of using the latest version of this module from the `master` branch, which 
  will change every time you run Terraform, you're using a fixed version of the repo.

* `ami_id`: Use this parameter to specify the ID of an Kapacitor [Amazon Machine Image 
  (AMI)](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/AMIs.html) to deploy on each server in the cluster. You
  should install Kapacitor the scripts in the 
  [install-kapacitor](https://github.com/gruntwork-io/terraform-aws-influx/tree/master/modules/install-kapacitor) module.
  
* `user_data`: Use this parameter to specify a [User 
  Data](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/user-data.html#user-data-shell-scripts) script that each
  server will run during boot. This is where you can use the 
  [run-kapacitor](https://github.com/gruntwork-io/terraform-aws-influx/tree/master/modules/run-kapacitor) 
  script to configure and run Kapacitor. 

You can find the other parameters in [variables.tf](variables.tf).

Check out the [examples folder](https://github.com/gruntwork-io/terraform-aws-influx/tree/master/examples) for 
fully-working sample code.

## What's included in this module?

This module creates the following:

* [Auto Scaling Group](#auto-scaling-group)
* [Security Group](#security-group)
* [IAM Role and Permissions](#iam-role-and-permissions)

### Auto Scaling Group

This module runs a single Kapacitor node on top of an [Auto Scaling Group (ASG)](https://aws.amazon.com/autoscaling/) to allow for auto-recovery. The EC2
Instance should be running an AMI that has Kapacitor installed via the 
[install-kapacitor](https://github.com/gruntwork-io/terraform-aws-influx/tree/master/modules/install-kapacitor)
module. You pass in the ID of the AMI to run using the `ami_id` input parameter.


### Security Group

The EC2 Instance in the ASG has a Security Group that allows minimal connectivity:

* All outbound requests
* Inbound SSH access from the CIDR blocks and security groups you specify

The Security Group ID is exported as an output variable which you can use with the 
[kapacitor-security-group-rules](https://github.com/gruntwork-io/terraform-aws-influx/tree/master/modules/kapacitor-security-group-rules)
module to open up all the ports necessary for Kapacitor.

### IAM Role and Permissions

The EC2 Instance in the ASG has an [IAM Role](http://docs.aws.amazon.com/IAM/latest/UserGuide/id_roles.html) attached. 
The IAM Role ARN and ID are exported as output variables if you need to add additional permissions.

### How do you roll out updates?

This module currently doesn't support updating the version of Kapacitor installed across the cluster. You can however follow
this [guide](https://docs.influxdata.com/kapacitor/v1.5/administration/upgrading/) to manually perform the upgrade.

### Dedicated instances

If you wish to use dedicated instances, you can set the `tenancy` parameter to `"dedicated"` in this module. 

### Encryption

This module does not currently support specifying encryption information. The official [documentation](
https://docs.influxdata.com/enterprise_kapacitor/v1.5/administration/security/#kapacitor-enterprise-over-tls) contains a guide for enabling SSL.

### Security groups

This module attaches a security group to each EC2 Instance that allows inbound requests as follows:

* **SSH**: For the SSH port (default: 22), you can use the `allowed_ssh_cidr_blocks` parameter to control the list of   
  [CIDR blocks](https://en.wikipedia.org/wiki/Classless_Inter-Domain_Routing) that will be allowed access. You can use 
  the `allowed_inbound_ssh_security_group_ids` parameter to control the list of source Security Groups that will be 
  allowed access.
  
The ID of the security group is exported as an output variable, which you can use with the 
[kapacitor-security-group-rules](https://github.com/gruntwork-io/terraform-aws-influx/tree/master/modules/kapacitor-security-group-rules)
module to open up all the ports necessary for Kapacitor.

### SSH access

You can associate an [EC2 Key Pair](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-key-pairs.html) with each
of the EC2 Instances in this cluster by specifying the Key Pair's name in the `ssh_key_name` variable. If you don't
want to associate a Key Pair with these servers, set `ssh_key_name` to an empty string.
