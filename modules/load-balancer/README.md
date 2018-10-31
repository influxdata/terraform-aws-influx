# Load Balancer

This folder contains a [Terraform](https://www.terraform.io/) module that can be used to deploy an [Application Load 
Balancer (ALB)](https://docs.aws.amazon.com/elasticloadbalancing/latest/application/introduction.html) in front of 
your InfluxDB data nodes cluster to:

1. Perform health checks on the servers in the cluster.
1. Distribute traffic from the Influx CLI and programming libraries across InfluxDB Data nodes.

Note that this module solely deploys the Load Balancer, as you may want to share one load balancer across multiple
applications. To deploy Target Groups, health checks, and routing rules, use the 
[load-balancer-target-group](https://github.com/gruntwork-io/terraform-aws-influx/tree/master/modules/load-balancer-target-group)
module.

See the [examples folder](https://github.com/gruntwork-io/terraform-aws-influx/tree/master/examples) for fully 
working sample code.

## How do you use this module?

This folder defines a [Terraform module](https://www.terraform.io/docs/modules/usage.html), which you can use in your
code by adding a `module` configuration and setting its `source` parameter to URL of this folder:

```hcl
module "load_balancer" {
  # TODO: replace <VERSION> with the latest version from the releases page: https://github.com/gruntwork-io/terraform-aws-influx/releases
  source = "github.com/gruntwork/terraform-aws-influx//modules/load-balancer?ref=<VERSION>"
  
  name       = "influxdb-load-balancer"
  vpc_id     = "vpc-abcd1234"
  subnet_ids = ["subnet-abcd1234", "subnet-efgh5678"]

  http_listener_ports = [8086]
  
  allow_inbound_from_cidr_blocks = ["0.0.0.0/0"]

  # ... See variables.tf for the other parameters you must define for this module
}
```

The above code will create a Load Balancer.

Note the following:

* `source`: Use this parameter in the `module` to specify the URL of the load-balancer module. The double slash (`//`) 
  is intentional and required. Terraform uses it to specify subfolders within a Git repo (see [module 
  sources](https://www.terraform.io/docs/modules/sources.html)). The `ref` parameter specifies a specific Git tag in 
  this repo. That way, instead of using the latest version of this module from the `master` branch, which 
  will change every time you run Terraform, you're using a fixed version of the repo.

* `http_listener_ports`: Which ports the load balancer should listen on for HTTP requests.

* `https_listener_ports_and_certs`: Whic ports the load balancer should listen on for HTTPS requests and which TLS
  certs to use with those ports.

* `allow_inbound_from_cidr_blocks`: Use this variable to specify which IP address ranges can connect to the Load
  Balancer. You can also use `allow_inbound_from_security_groups` to allow specific security groups to connect.

## How is the ALB configured?

The ALB in this module is configured as follows:

1. **Listeners**: The Load Balancer will create a listener for each port specified in `http_listener_ports`.
   
1. **DNS**: You can use the `route53_records` variable to create one more more DNS A Records in [Route 
   53](https://aws.amazon.com/route53/) that point to the Load Balancer. This allows you to use custom domain names to access the Load Balancer.
