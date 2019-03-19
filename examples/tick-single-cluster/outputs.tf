output "cluster_asg_name" {
  value = "${module.tick.asg_name}"
}

output "lb_dns_name" {
  value = "${module.load_balancer.alb_dns_name}"
}

output "influxdb_port" {
  value = "${module.influxdb_security_group_rules.api_port}"
}

output "chronograf_port" {
  value = "${module.chronograf_security_group_rules.http_port}"
}

output "kapacitor_port" {
  value = "${module.kapacitor_security_group_rules.http_port}"
}
