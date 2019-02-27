output "public_ip" {
  value = "${aws_instance.chronograf_server.public_ip}"
}

output "security_group_id" {
  value = "${aws_security_group.chronograf_security_group.id}"
}
