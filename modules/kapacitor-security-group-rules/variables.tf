# ---------------------------------------------------------------------------------------------------------------------
# REQUIRED PARAMETERS
# You must provide a value for each of these parameters.
# ---------------------------------------------------------------------------------------------------------------------

variable "security_group_id" {
  description = "The ID of the Security Group to which all the rules should be attached."
}

variable "http_port" {
  description = "The HTTP port that Kapacitor will listen on for connections"
}

# ---------------------------------------------------------------------------------------------------------------------
# OPTIONAL PARAMETERS
# These parameters have reasonable defaults.
# ---------------------------------------------------------------------------------------------------------------------

variable "http_port_cidr_blocks" {
  description = "The list of IP address ranges in CIDR notation from which to allow connections to the http_port."
  type        = "list"
  default     = []
}

variable "http_port_security_groups" {
  description = "The list of Security Group IDs from which to allow connections to the http_port. If you update this variable, make sure to update var.http_port_security_groups_num too!"
  type        = "list"
  default     = []
}

variable "http_port_security_groups_num" {
  description = "The number of security group IDs in var.http_port_security_groups. We should be able to compute this automatically, but due to a Terraform limitation, if there are any dynamic resources in var.http_port_security_groups, then we won't be able to: https://github.com/hashicorp/terraform/pull/11482"
  default     = 0
}
