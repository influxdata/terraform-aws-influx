# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY A SINGLE NODE TICK CLUSTER
# ---------------------------------------------------------------------------------------------------------------------

provider "aws" {
  # The AWS region in which all resources will be created
  region = "${var.aws_region}"
}

module "tick" {
  # When using these modules in your own code, you will need to use a Git URL with a ref attribute that pins you
  # to a specific version of the modules, such as the following example:
  # source = "git::git@github.com:gruntwork-io/terraform-aws-influx.git//modules/influxdb-cluster?ref=v0.0.1"
  source = "../../modules/influxdb-cluster"

  cluster_name = "${var.cluster_name}"
  min_size     = 1
  max_size     = 1

  # We use small instance types to keep these examples cheap to run. In a production setting, you'll probably want
  # R4 or M4 instances.
  instance_type = "t2.micro"

  ami_id    = "${var.ami_id}"
  user_data = "${data.template_file.user_data_tick.rendered}"

  vpc_id     = "${data.aws_vpc.default.id}"
  subnet_ids = "${data.aws_subnet_ids.default.ids}"

  ebs_block_devices = [
    {
      device_name = "${var.influxdb_volume_device_name}"
      volume_type = "gp2"
      volume_size = 50
    },
    {
      device_name = "${var.kapacitor_volume_device_name}"
      volume_type = "gp2"
      volume_size = 50
    },
  ]

  # To make testing easier, we allow SSH requests from any IP address here. In a production deployment, we strongly
  # recommend you limit this to the IP address ranges of known, trusted servers inside your VPC.
  allowed_ssh_cidr_blocks = ["0.0.0.0/0"]

  ssh_key_name = "${var.ssh_key_name}"

  # To make it easy to test this example from your computer, we allow the InfluxDB servers to have public IPs. In a
  # production deployment, you'll probably want to keep all the servers in private subnets with only private IPs.
  associate_public_ip_address = true

  # An example of custom tags
  tags = [
    {
      key                 = "Environment"
      value               = "development"
      propagate_at_launch = true
    },
    {
      key                 = "NodeType"
      value               = "both"
      propagate_at_launch = true
    },
  ]
}

# ---------------------------------------------------------------------------------------------------------------------
# CREATE USER DATA SCRIPTS THAT WILL RUN ON EACH INSTANCE IN THE VARIOUS CLUSTERS ON BOOT
# ---------------------------------------------------------------------------------------------------------------------

data "template_file" "user_data_tick" {
  template = "${file("${path.module}/user-data/user-data.sh")}"

  vars {
    # InfluxDB
    cluster_asg_name = "${var.cluster_name}"
    aws_region       = "${var.aws_region}"
    license_key      = "${var.license_key}"
    shared_secret    = "${var.shared_secret}"

    # Pass in the data about the EBS volumes so they can be mounted
    influxdb_volume_device_name = "${var.influxdb_volume_device_name}"
    influxdb_volume_mount_point = "${var.influxdb_volume_mount_point}"
    influxdb_volume_owner       = "${var.influxdb_volume_owner}"

    # Telegraf
    influxdb_url  = "http://localhost:8086"
    database_name = "${var.telegraf_database}"

    # Chronograf
    host = "0.0.0.0"
    port = "8888"

    # Kapacitor
    hostname                     = "localhost"
    kapacitor_volume_device_name = "${var.kapacitor_volume_device_name}"
    kapacitor_volume_mount_point = "${var.kapacitor_volume_mount_point}"
    kapacitor_volume_owner       = "${var.kapacitor_volume_owner}"
  }
}

# ---------------------------------------------------------------------------------------------------------------------
# CONFIGURE THE SECURITY GROUP RULES FOR INFLUXDB
# This controls which ports are exposed and who can connect to them
# ---------------------------------------------------------------------------------------------------------------------

module "influxdb_security_group_rules" {
  # When using these modules in your own code, you will need to use a Git URL with a ref attribute that pins you
  # to a specific version of the modules, such as the following example:
  # source = "git::git@github.com:gruntwork-io/terraform-aws-influx.git//modules/influxdb-security-group-rules?ref=v0.0.1"
  source = "../../modules/influxdb-security-group-rules"

  security_group_id = "${module.tick.security_group_id}"

  raft_port = 8089
  rest_port = 8091
  tcp_port  = 8088
  api_port  = 8086

  # To keep this example simple, we allow these ports to be accessed from any IP. In a production
  # deployment, you may want to lock these down just to trusted servers.
  raft_port_cidr_blocks = ["0.0.0.0/0"]

  rest_port_cidr_blocks = ["0.0.0.0/0"]
  tcp_port_cidr_blocks  = ["0.0.0.0/0"]
  api_port_cidr_blocks  = ["0.0.0.0/0"]
}

module "chronograf_security_group_rules" {
  # When using these modules in your own code, you will need to use a Git URL with a ref attribute that pins you
  # to a specific version of the modules, such as the following example:
  # source = "git::git@github.com:gruntwork-io/terraform-aws-influx.git//modules/chronograf-security-group-rules?ref=v0.0.1"
  source = "../../modules/chronograf-security-group-rules"

  security_group_id = "${module.tick.security_group_id}"

  http_port = 8888

  # To keep this example simple, we allow these ports to be accessed from any IP. In a production
  # deployment, you may want to lock these down just to trusted servers.
  http_port_cidr_blocks = ["0.0.0.0/0"]
}

module "kapacitor_security_group_rules" {
  # When using these modules in your own code, you will need to use a Git URL with a ref attribute that pins you
  # to a specific version of the modules, such as the following example:
  # source = "git::git@github.com:gruntwork-io/terraform-aws-influx.git//modules/kapacitor-security-group-rules?ref=v0.0.1"
  source = "../../modules/kapacitor-security-group-rules"

  security_group_id = "${module.tick.security_group_id}"

  http_port = 9092

  # To keep this example simple, we allow these ports to be accessed from any IP. In a production
  # deployment, you may want to lock these down just to trusted servers.
  http_port_cidr_blocks = ["0.0.0.0/0"]
}

# ---------------------------------------------------------------------------------------------------------------------
# ATTACH IAM POLICIES TO EACH CLUSTER
# These policies allow the clusters to automatically bootstrap themselves
# ---------------------------------------------------------------------------------------------------------------------

module "influxdb_iam_policies" {
  # When using these modules in your own code, you will need to use a Git URL with a ref attribute that pins you
  # to a specific version of the modules, such as the following example:
  # source = "git::git@github.com:gruntwork-io/terraform-aws-influx.git//modules/influxdb-iam-policies?ref=v0.0.1"
  source = "../../modules/influxdb-iam-policies"

  iam_role_id = "${module.tick.iam_role_id}"
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY A LOAD BALANCER FOR THE CLUSTERS
# ---------------------------------------------------------------------------------------------------------------------

module "load_balancer" {
  # When using these modules in your own code, you will need to use a Git URL with a ref attribute that pins you
  # to a specific version of the modules, such as the following example:
  # source = "git::git@github.com:gruntwork-io/terraform-aws-influx.git//modules/load-balancer?ref=v0.0.1"
  source = "../../modules/load-balancer"

  name       = "${var.cluster_name}-lb"
  vpc_id     = "${data.aws_vpc.default.id}"
  subnet_ids = "${data.aws_subnet_ids.default.ids}"

  http_listener_ports = [8086, 8888, 9092]

  # To make testing easier, we allow inbound connections from any IP. In production usage, you may want to only allow
  # connectsion from certain trusted servers, or even use an internal load balancer, so it's only accessible from
  # within the VPC

  allow_inbound_from_cidr_blocks = ["0.0.0.0/0"]
  idle_timeout                   = 3600
}

module "influxdb_target_group" {
  # When using these modules in your own code, you will need to use a Git URL with a ref attribute that pins you
  # to a specific version of the modules, such as the following example:
  # source = "git::git@github.com:gruntwork-io/terraform-aws-influx.git//modules/load-balancer-target-group?ref=v0.0.1"
  source = "../../modules/load-balancer-target-group"

  target_group_name    = "${var.cluster_name}-itg"
  asg_name             = "${module.tick.asg_name}"
  port                 = "${module.influxdb_security_group_rules.api_port}"
  health_check_path    = "/ping"
  health_check_matcher = "204"
  vpc_id               = "${data.aws_vpc.default.id}"

  listener_arns                   = ["${lookup(module.load_balancer.http_listener_arns, 8086)}"]
  listener_arns_num               = 1
  listener_rule_starting_priority = 100
}

module "chronograf_target_group" {
  # When using these modules in your own code, you will need to use a Git URL with a ref attribute that pins you
  # to a specific version of the modules, such as the following example:
  # source = "git::git@github.com:gruntwork-io/terraform-aws-influx.git//modules/load-balancer-target-group?ref=v0.0.1"
  source = "../../modules/load-balancer-target-group"

  target_group_name    = "${var.cluster_name}-ctg"
  asg_name             = "${module.tick.asg_name}"
  port                 = "${module.chronograf_security_group_rules.http_port}"
  health_check_path    = "/"
  health_check_matcher = "200"
  vpc_id               = "${data.aws_vpc.default.id}"

  listener_arns                   = ["${lookup(module.load_balancer.http_listener_arns, 8888)}"]
  listener_arns_num               = 1
  listener_rule_starting_priority = 100
}

module "kapacitor_target_group" {
  # When using these modules in your own code, you will need to use a Git URL with a ref attribute that pins you
  # to a specific version of the modules, such as the following example:
  # source = "git::git@github.com:gruntwork-io/terraform-aws-influx.git//modules/load-balancer-target-group?ref=v0.0.1"
  source = "../../modules/load-balancer-target-group"

  target_group_name    = "${var.cluster_name}-ktg"
  asg_name             = "${module.tick.asg_name}"
  port                 = "${module.kapacitor_security_group_rules.http_port}"
  health_check_path    = "/kapacitor/v1/ping"
  health_check_matcher = "204"
  vpc_id               = "${data.aws_vpc.default.id}"

  listener_arns                   = ["${lookup(module.load_balancer.http_listener_arns, 9092)}"]
  listener_arns_num               = 1
  listener_rule_starting_priority = 100
}

# ---------------------------------------------------------------------------------------------------------------------
# DEPLOY INFLUXDB IN THE DEFAULT VPC AND SUBNETS
# Using the default VPC and subnets makes this example easy to run and test, but it means InfluxDB is accessible from
# the public Internet. For a production deployment, we strongly recommend deploying into a custom VPC with private
# subnets.
# ---------------------------------------------------------------------------------------------------------------------

data "aws_vpc" "default" {
  default = true
}

data "aws_subnet_ids" "default" {
  vpc_id = "${data.aws_vpc.default.id}"
}
