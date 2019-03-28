output "telegraf_database" {
  value = "${var.telegraf_database}"
}

output "influxdb_dns" {
  value = "${module.influxdb_load_balancer.alb_dns_name}"
}

output "influxdb_port" {
  value = "${module.influxdb_data_nodes_security_group_rules.api_port}"
}

output "chronograf_dns" {
  value = "${module.chronograf_load_balancer.alb_dns_name}"
}

output "chronograf_port" {
  value = "${module.chronograf_security_group_rules.http_port}"
}

output "kapacitor_dns" {
  value = "${module.kapacitor_load_balancer.alb_dns_name}"
}

output "kapacitor_port" {
  value = "${module.kapacitor_security_group_rules.http_port}"
}
