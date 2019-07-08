# ----------------------------------------------------------------------------------------------------------------------
# REQUIRE A SPECIFIC TERRAFORM VERSION OR HIGHER
# This module has been updated with 0.12 syntax, which means it is no longer compatible with any versions below 0.12.
# ----------------------------------------------------------------------------------------------------------------------

terraform {
  required_version = ">= 0.12"
}

resource "aws_security_group_rule" "raft_port_cidr_blocks" {
  count             = signum(length(var.raft_port_cidr_blocks))
  type              = "ingress"
  from_port         = var.raft_port
  to_port           = var.raft_port
  protocol          = "tcp"
  security_group_id = var.security_group_id
  cidr_blocks       = var.raft_port_cidr_blocks
}

resource "aws_security_group_rule" "raft_port_security_groups" {
  count                    = var.raft_port_security_groups_num == null ? 0 : var.raft_port_security_groups_num
  type                     = "ingress"
  from_port                = var.raft_port
  to_port                  = var.raft_port
  protocol                 = "tcp"
  security_group_id        = var.security_group_id
  source_security_group_id = element(var.raft_port_security_groups, count.index)
}

resource "aws_security_group_rule" "rest_port_cidr_blocks" {
  count             = signum(length(var.rest_port_cidr_blocks))
  type              = "ingress"
  from_port         = var.rest_port
  to_port           = var.rest_port
  protocol          = "tcp"
  security_group_id = var.security_group_id
  cidr_blocks       = var.rest_port_cidr_blocks
}

resource "aws_security_group_rule" "rest_port_security_groups" {
  count                    = var.rest_port_security_groups_num == null ? 0 : var.rest_port_security_groups_num
  type                     = "ingress"
  from_port                = var.rest_port
  to_port                  = var.rest_port
  protocol                 = "tcp"
  security_group_id        = var.security_group_id
  source_security_group_id = element(var.rest_port_security_groups, count.index)
}

resource "aws_security_group_rule" "tcp_port_cidr_blocks" {
  count             = signum(length(var.tcp_port_cidr_blocks))
  type              = "ingress"
  from_port         = var.tcp_port
  to_port           = var.tcp_port
  protocol          = "tcp"
  security_group_id = var.security_group_id
  cidr_blocks       = var.tcp_port_cidr_blocks
}

resource "aws_security_group_rule" "tcp_port_security_groups" {
  count                    = var.tcp_port_security_groups_num == null ? 0 : var.tcp_port_security_groups_num
  type                     = "ingress"
  from_port                = var.tcp_port
  to_port                  = var.tcp_port
  protocol                 = "tcp"
  security_group_id        = var.security_group_id
  source_security_group_id = element(var.tcp_port_security_groups, count.index)
}

resource "aws_security_group_rule" "api_port_cidr_blocks" {
  count             = signum(length(var.api_port_cidr_blocks))
  type              = "ingress"
  from_port         = var.api_port
  to_port           = var.api_port
  protocol          = "tcp"
  security_group_id = var.security_group_id
  cidr_blocks       = var.api_port_cidr_blocks
}

resource "aws_security_group_rule" "api_port_security_groups" {
  count                    = var.api_port_security_groups_num == null ? 0 : var.api_port_security_groups_num
  type                     = "ingress"
  from_port                = var.api_port
  to_port                  = var.api_port
  protocol                 = "tcp"
  security_group_id        = var.security_group_id
  source_security_group_id = element(var.api_port_security_groups, count.index)
}
