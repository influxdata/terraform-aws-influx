output "influxdb_meta_nodes_cluster_asg_name" {
  value = "${module.influxdb_meta_nodes.asg_name}"
}

output "influxdb_data_nodes_cluster_asg_name" {
  value = "${module.influxdb_data_nodes.asg_name}"
}

output "lb_dns_name" {
  value = "${module.load_balancer.alb_dns_name}"
}

output "domain_names" {
  value = "${module.load_balancer.domain_names}"
}
