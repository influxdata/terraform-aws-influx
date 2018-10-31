output "alb_arn" {
  value = "${aws_alb.lb.arn}"
}

output "alb_name" {
  value = "${aws_alb.lb.name}"
}

output "alb_dns_name" {
  value = "${aws_alb.lb.dns_name}"
}

output "domain_names" {
  value = "${aws_route53_record.load_balancer.*.fqdn}"
}

output "http_listener_arns" {
  value = "${zipmap(var.http_listener_ports, aws_alb_listener.http.*.arn)}"
}

output "security_group_id" {
  value = "${aws_security_group.sg.id}"
}
