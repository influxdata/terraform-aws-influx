# Chronograf Server

This folder contains a [Terraform](https://www.terraform.io/) module to deploy a [Chronograf](
https://www.influxdata.com/time-series-platform/chronograf/) server in [AWS](https://aws.amazon.com/). 
The idea is to create an [Amazon Machine Image (AMI)](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/AMIs.html)
that has the Chronograf binary installed using the [install-chronograf](
https://github.com/gruntwork-io/terraform-aws-influx/tree/master/modules/install-chronograf) module.

## How do you use this module?

This folder defines a [Terraform module](https://www.terraform.io/docs/modules/usage.html), which you can use in your
code by adding a `module` configuration and setting its `source` parameter to URL of this folder:

```hcl
module "chronograf_server" {
  # TODO: replace <VERSION> with the latest version from the releases page: https://github.com/gruntwork-io/terraform-aws-influx/releases
  source = "github.com/gruntwork-io/terraform-aws-influx//modules/chronograf-server?ref=<VERSION>"

  # Specify the ID of the Chronograf AMI. You should build this using the scripts in the install-chronograf module.
  ami_id = "ami-abcd1234"
  
  # Configure and start Chronograf during boot. 
  user_data = <<-EOF
              #!/bin/bash
              sudo systemctl start chronograf
              EOF
  
  # ... See variables.tf for the other parameters you must define for the chronograf-server module
}
```

Note the following parameters:

* `source`: Use this parameter to specify the URL of the chronograf-server module. The double slash (`//`) is 
  intentional and required. Terraform uses it to specify subfolders within a Git repo (see [module 
  sources](https://www.terraform.io/docs/modules/sources.html)). The `ref` parameter specifies a specific Git tag in 
  this repo. That way, instead of using the latest version of this module from the `master` branch, which 
  will change every time you run Terraform, you're using a fixed version of the repo.

* `ami_id`: Use this parameter to specify the ID of an Chronograf [Amazon Machine Image 
  (AMI)](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/AMIs.html) to deploy on each server in the cluster. You
  should install Chronograf the scripts in the 
  [install-chronograf](https://github.com/gruntwork-io/terraform-aws-influx/tree/master/modules/install-chronograf) module.
  
* `user_data`: Use this parameter to specify a [User 
  Data](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/user-data.html#user-data-shell-scripts) script that each
  server will run during boot. This is where you can use the 
  [run-chronograf](https://github.com/gruntwork-io/terraform-aws-influx/tree/master/modules/run-chronograf) 
  script to configure and run Chronograf as either a meta node or a data node. 

You can find the other parameters in [variables.tf](variables.tf).

Check out the [examples folder](https://github.com/gruntwork-io/terraform-aws-influx/tree/master/examples) for 
fully-working sample code.

## How do you connect to the Chronograf server?

Once deployed you can simply access the web UI by visiting the http://<public-ip>:<port>, you can get the `public_ip` by running:

```bash
$ terraform output public_ip
```

## What's included in this module?

This module creates the following:

* [Security Group](#security-group)

### Security Group

The EC2 Instance has a Security Group that allows minimal connectivity:

* All outbound requests
* Inbound SSH access from the CIDR blocks and security groups you specify

The Security Group ID is exported as an output variable which you can use with the 
[chronograf-security-group-rules](https://github.com/gruntwork-io/terraform-aws-influx/tree/master/modules/chronograf-security-group-rules)
module to open up all the ports necessary for Chronograf.

### Dedicated instances

If you wish to use dedicated instances, you can set the `tenancy` parameter to `"dedicated"` in this module. 

### Encryption

This module does not currently support specifying encryption information. The official [documentation](
https://docs.influxdata.com/chronograf/v1.7/administration/managing-security/#configuring-tls-transport-layer-security-and-https) contains a guide for enabling SSL.

### Security groups

This module attaches a security group to the EC2 Instance that allows inbound requests as follows:

* **SSH**: For the SSH port (default: 22), you can use the `allowed_ssh_cidr_blocks` parameter to control the list of   
  [CIDR blocks](https://en.wikipedia.org/wiki/Classless_Inter-Domain_Routing) that will be allowed access. You can use 
  the `allowed_inbound_ssh_security_group_ids` parameter to control the list of source Security Groups that will be 
  allowed access.
  
The ID of the security group is exported as an output variable, which you can use with the 
[chronograf-security-group-rules](https://github.com/gruntwork-io/terraform-aws-influx/tree/master/modules/chronograf-security-group-rules)
module to open up all the ports necessary for Chronograf.

### SSH access

You can associate an [EC2 Key Pair](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-key-pairs.html) with the EC2 Instance by specifying the Key Pair's name in the `ssh_key_name` variable. If you don't
want to associate a Key Pair with these servers, set `ssh_key_name` to an empty string.
