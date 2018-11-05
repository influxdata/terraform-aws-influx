output "influxdb_cluster_asg_name" {
  value = "${module.influxdb.asg_name}"
}

output "lb_dns_name" {
  value = "${module.load_balancer.alb_dns_name}"
}

output "load_balancer_port" {
  value = "${module.influxdb_security_group_rules.api_port}"
}
