# Chronograf Server Security Group Rules Module

This folder contains a [Terraform](https://www.terraform.io/) module that defines the Security Group rules used by a 
[Chronograf](https://www.influxdata.com/time-series-platform/chronograf/) server to control the traffic that is allowed to go in and out of the server. 
These rules are defined in a separate module so that you can add them to any existing Security Group. 

## Quick start

Let's say you want to deploy Chronograf using the [chronograf-server 
module](https://github.com/gruntwork-io/terraform-aws-influx/blob/master/modules/chronograf-server): 

```hcl
module "chronograf_server" {
  # TODO: replace <VERSION> with the latest version from the releases page: https://github.com/gruntwork-io/terraform-aws-influx/releases
  source = "github.com/gruntwork-io/terraform-aws-influx//modules/chronograf-server?ref=<VERSION>"

  # ... (other params omitted) ...
}
```

You can attach the Security Group rules to this server as follows:

```hcl
module "security_group_rules" {
  # TODO: replace <VERSION> with the latest version from the releases page: https://github.com/gruntwork-io/terraform-aws-influx/releases
  source = "github.com/gruntwork-io/terraform-aws-influx//modules/chronograf-security-group-rules?ref=<VERSION>"

  security_group_id = "${module.chronograf_server.security_group_id}"
  
  http_port                 = 8888
  http_port_cidr_blocks     = ["0.0.0.0/0"]
  http_port_security_groups = ["sg-abcd1234"]
  
  # ... (other params omitted) ...
}
```

Note the following parameters:

* `source`: Use this parameter to specify the URL of this module. The double slash (`//`) is intentional 
  and required. Terraform uses it to specify subfolders within a Git repo (see [module 
  sources](https://www.terraform.io/docs/modules/sources.html)). The `ref` parameter specifies a specific Git tag in 
  this repo. That way, instead of using the latest version of this module from the `master` branch, which 
  will change every time you run Terraform, you're using a fixed version of the repo.

* `security_group_id`: Use this parameter to specify the ID of the security group to which the rules in this module
  should be added.

* `http_port`, `http_port_cidr_blocks`, `http_port_security_groups`: This shows an example of how to configure which 
  ports you're using for various Chronograf functionality and which IP address ranges and 
  Security Groups are allowed to connect to that port.
  
You can find the other parameters in [variables.tf](variables.tf).

Check out the [examples folder](https://github.com/gruntwork-io/terraform-aws-influx/blob/master/examples) for 
working sample code.
