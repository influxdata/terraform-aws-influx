# ---------------------------------------------------------------------------------------------------------------------
# REQUIRED PARAMETERS
# You must provide a value for each of these parameters.
# ---------------------------------------------------------------------------------------------------------------------

variable "security_group_id" {
  description = "The ID of the Security Group to which all the rules should be attached."
}

variable "raft_port" {
  description = "The Raft consensus protocol port on which Meta/Data nodes communicate with each other"
}

variable "rest_port" {
  description = "The HTTP API port the Meta/Data nodes listen on for internal communication."
}

variable "tcp_port" {
  description = "The port the Meta/Data nodes use for internal communication via a TCP protocol."
}

# ---------------------------------------------------------------------------------------------------------------------
# OPTIONAL PARAMETERS
# These parameters have reasonable defaults.
# ---------------------------------------------------------------------------------------------------------------------

variable "api_port" {
  description = "The HTTP API port the Data nodes listen on for external communication."
  default     = ""
}

variable "raft_port_cidr_blocks" {
  description = "The list of IP address ranges in CIDR notation from which to allow connections to the raft_port."
  type        = "list"
  default     = []
}

variable "raft_port_security_groups" {
  description = "The list of Security Group IDs from which to allow connections to the raft_port. If you update this variable, make sure to update var.raft_port_security_groups_num too!"
  type        = "list"
  default     = []
}

variable "raft_port_security_groups_num" {
  description = "The number of security group IDs in var.raft_port_security_groups. We should be able to compute this automatically, but due to a Terraform limitation, if there are any dynamic resources in var.raft_port_security_groups, then we won't be able to: https://github.com/hashicorp/terraform/pull/11482"
  default     = 0
}

variable "rest_port_cidr_blocks" {
  description = "The list of IP address ranges in CIDR notation from which to allow connections to the rest_port."
  type        = "list"
  default     = []
}

variable "rest_port_security_groups" {
  description = "The list of Security Group IDs from which to allow connections to the rest_port. If you update this variable, make sure to update var.rest_port_security_groups_num too!"
  type        = "list"
  default     = []
}

variable "rest_port_security_groups_num" {
  description = "The number of security group IDs in var.rest_port_security_groups. We should be able to compute this automatically, but due to a Terraform limitation, if there are any dynamic resources in var.rest_port_security_groups, then we won't be able to: https://github.com/hashicorp/terraform/pull/11482"
  default     = 0
}

variable "tcp_port_cidr_blocks" {
  description = "The list of IP address ranges in CIDR notation from which to allow connections to the tcp_port."
  type        = "list"
  default     = []
}

variable "tcp_port_security_groups" {
  description = "The list of Security Group IDs from which to allow connections to the tcp_port. If you update this variable, make sure to update var.tcp_port_security_groups_num too!"
  type        = "list"
  default     = []
}

variable "tcp_port_security_groups_num" {
  description = "The number of security group IDs in var.tcp_port_security_groups. We should be able to compute this automatically, but due to a Terraform limitation, if there are any dynamic resources in var.tcp_port_security_groups, then we won't be able to: https://github.com/hashicorp/terraform/pull/11482"
  default     = 0
}

variable "api_port_cidr_blocks" {
  description = "The list of IP address ranges in CIDR notation from which to allow connections to the api_port."
  type        = "list"
  default     = []
}

variable "api_port_security_groups" {
  description = "The list of Security Group IDs from which to allow connections to the api_port. If you update this variable, make sure to update var.api_port_security_groups_num too!"
  type        = "list"
  default     = []
}

variable "api_port_security_groups_num" {
  description = "The number of security group IDs in var.api_port_security_groups. We should be able to compute this automatically, but due to a Terraform limitation, if there are any dynamic resources in var.api_port_security_groups, then we won't be able to: https://github.com/hashicorp/terraform/pull/11482"
  default     = 0
}
