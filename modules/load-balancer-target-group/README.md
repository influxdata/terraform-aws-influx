# Load Balancer Target Group Module

This module can be used to create a [Target 
Group](https://docs.aws.amazon.com/elasticloadbalancing/latest/application/load-balancer-target-groups.html) and
[Listener Rules](https://docs.aws.amazon.com/elasticloadbalancing/latest/application/listener-update-rules.html) for
a Load Balancer created with the [load-balancer 
module](https://github.com/gruntwork-io/terraform-aws-influx/tree/master/modules/load-balancer). You can use this 
module to configure health checks and routing for InfluxDB data nodes. 

The reason the `load-balancer` and `load-balancer-target-group` modules are separate is that you may wish to create
multiple target groups for a single load balancer.

See the [examples folder](https://github.com/gruntwork-io/terraform-aws-influx/tree/master/examples) for fully
working sample code.

## How do you use this module?

Imagine you've deployed InfluxDB using the [influxdb-cluster
module](https://github.com/gruntwork-io/terraform-aws-influx/tree/master/modules/influxdb-cluster) and a Load Balancer
using the [load-balancer module](https://github.com/gruntwork-io/terraform-aws-influx/tree/master/modules/load-balancer):    

```hcl
module "influxdb_data_nodes" {
  # TODO: replace <VERSION> with the latest version from the releases page: https://github.com/gruntwork-io/terraform-aws-influx/releases
  source = "github.com/gruntwork/terraform-aws-influx//modules/influxdb-cluster?ref=<VERSION>"
  
  cluster_name = "${var.cluster_name}"
  
  # ... (other params omitted) ...
}

module "load_balancer" {
  # TODO: replace <VERSION> with the latest version from the releases page: https://github.com/gruntwork-io/terraform-aws-influx/releases
  source = "github.com/gruntwork/terraform-aws-influx//modules/load-balancer?ref=<VERSION>"
  
  name = "${var.cluster_name}"

  http_listener_ports = [8086]

  # ... (other params omitted) ...
}
``` 

Note the following:

* `http_listener_ports`: This tells the Load Balancer to listen for HTTP requests on port 8091 and 4984.
  
To create Target Groups and Listener Rules for InfluxDB, you need to use the
`load-balancer-target-group` module as follows:

```hcl
module "influxdb_target_group" {
  # TODO: replace <VERSION> with the latest version from the releases page: https://github.com/gruntwork-io/terraform-aws-influx/releases
  source = "github.com/gruntwork/terraform-aws-influx//modules/load-balancer-target-group?ref=<VERSION>"

  target_group_name = "${var.cluster_name}-cb"
  asg_name          = "${module.influxdb_data_nodes.asg_name}"
  port              = 8086
  health_check_path = "/ping"

  listener_arns                   = ["${lookup(module.load_balancer.http_listener_arns, 8086)}"]
  num_listener_arns               = 1
  listener_rule_starting_priority = 100
    
  # ... See variables.tf for the other parameters you must define for this module
}
```

Note the following:

* `asg_name`: Use this param to attach the Target Group to the Auto Scaling Group (ASG) used under the hood in the
  InfluxDB cluster so that each EC2 Instance automatically registers with the Target Group, goes 
  through health checks, and gets replaced if it is failing health checks. 

* `listener_arns`: Specify the ARN of the HTTP listener from the Load Balancer module. The InfluxDB Target Group uses
  InfluxDB's port (8086).
  