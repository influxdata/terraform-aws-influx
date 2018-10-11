# InfluxDB Server Security Group Rules Module

This folder contains a [Terraform](https://www.terraform.io/) module that defines the Security Group rules used by a 
[InfluxDB Enterprise](https://www.influxdata.com/time-series-platform/influxdb/) cluster to control the traffic that is allowed to go in and out of the cluster. 
These rules are defined in a separate module so that you can add them to any existing Security Group. 

## Quick start

Let's say you want to deploy influx using the [influx-cluster 
module](https://github.com/gruntwork-io/terraform-aws-influx/blob/master/modules/influx-cluster): 

```hcl
module "influxdb_meta_cluster" {
  # TODO: replace <VERSION> with the latest version from the releases page: https://github.com/gruntwork-io/terraform-aws-influx/releases
  source = "github.com/gruntwork-io/terraform-aws-influx//modules/influxdb-cluster?ref=<VERSION>"

  # ... (other params omitted) ...
}
```

You can attach the Security Group rules to this cluster as follows:

```hcl
module "security_group_rules" {
  # TODO: replace <VERSION> with the latest version from the releases page: https://github.com/gruntwork-io/terraform-aws-influx/releases
  source = "github.com/gruntwork-io/terraform-aws-influx//modules/influxdb-security-group-rules?ref=<VERSION>"

  security_group_id = "${module.influxdb_meta_cluster.security_group_id}"
  
  rest_port                 = 8091
  rest_port_cidr_blocks     = ["0.0.0.0/0"]
  rest_port_security_groups = ["sg-abcd1234"]
  
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

* `raft_port`, `raft_port_cidr_blocks`, `raft_port_security_groups`: This shows an example of how to configure which 
  ports you're using for various InfluxDB functionality, such as the RAFT consensus port, and which IP address ranges and 
  Security Groups are allowed to connect to that port. Check out the [Clustering 
  documentation](https://docs.influxdata.com/enterprise_influxdb/v1.6/concepts/clustering) to understand
  what ports InfluxDB uses.
  
You can find the other parameters in [variables.tf](variables.tf).

Check out the [examples folder](https://github.com/gruntwork-io/terraform-aws-influx/blob/master/examples) for 
working sample code.
