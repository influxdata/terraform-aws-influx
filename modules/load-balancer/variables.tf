# ---------------------------------------------------------------------------------------------------------------------
# REQUIRED PARAMETERS
# You must provide a value for each of these parameters.
# ---------------------------------------------------------------------------------------------------------------------

variable "name" {
  description = "The name to use for the Load Balancer"
}

variable "http_listener_ports" {
  description = "A list of ports to listen on for HTTP requests."
  type        = "list"
}

variable "allow_inbound_from_cidr_blocks" {
  description = "A list of IP addresses in CIDR notation from which the load balancer will allow incoming HTTP/HTTPS requests."
  type        = "list"
}

variable "vpc_id" {
  description = "The ID of the VPC in which to deploy the Load Balancer"
}

variable "subnet_ids" {
  description = "The subnet IDs into which the Load Balancer should be deployed."
  type        = "list"
}

# ---------------------------------------------------------------------------------------------------------------------
# OPTIONAL PARAMETERS
# These parameters have reasonable defaults.
# ---------------------------------------------------------------------------------------------------------------------

variable "allow_inbound_from_security_groups" {
  description = "A list of Security Group IDs from which the load balancer will allow incoming HTTP/HTTPS requests. Any time you change this value, make sure to update var.allow_inbound_from_security_groups too!"
  type        = "list"
  default     = []
}

variable "allow_inbound_from_security_groups_num" {
  description = "The number of Security Group IDs in var.allow_inbound_from_security_groups. We should be able to compute this automatically, but due to a Terraform limitation, if there are any dynamic resources in var.allow_inbound_from_cidr_blocks, then we won't be able to: https://github.com/hashicorp/terraform/pull/11482"
  default     = 0
}

variable "default_target_group_arn" {
  description = "The ARN of a Target Group where all requests that don't match any Load Balancer Listener Rules will be sent. If you set this to empty string, we will send the requests to a \"black hole\" target group that always returns a 503, so we strongly recommend configuring this to be a target group that can instead return a reasonable 404 page."
  default     = ""
}

variable "internal" {
  description = "Set to true to make this an internal load balancer that is only accessible from within the VPC. Set to false to make it publicly accessible."
  default     = false
}

variable "enable_http2" {
  description = "Set to true to enable HTTP/2 on the load balancer."
  default     = true
}

variable "ip_address_type" {
  description = "The type of IP address to use on the load balancer. Must be one of: ipv4, dualstack."
  default     = "ipv4"
}

variable "tags" {
  description = "Custom tags to apply to the load balancer."
  type        = "map"
  default     = {}
}

variable "idle_timeout" {
  description = "The time in seconds that the connection is allowed to be idle."
  default     = 30
}

variable "route53_records" {
  description = "A list of DNS A records to create in Route 53 that point at this Load Balancer. Each item in the list should be an object with the keys 'domain' (the domain name to create) and 'zone_id' (the Route 53 Hosted Zone ID in which to create the DNS A record)."
  type        = "list"
  default     = []

  # Example:
  #
  # default = [
  #   {
  #     domain  = "foo.acme.com"
  #     zone_id = "Z1234ABCDEFG"
  #   }
  # ]
}
