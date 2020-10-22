# ---------------------------------------------------------------------------------------------------------------------
# ATTACH SECURITY GROUP RULE TO ALLOW INCOMING CONNECTIONS ON CHRONGRAF'S LISTENING PORT
# ---------------------------------------------------------------------------------------------------------------------

# ----------------------------------------------------------------------------------------------------------------------
# REQUIRE A SPECIFIC TERRAFORM VERSION OR HIGHER
# This module has been updated with 0.12 syntax, which means it is no longer compatible with any versions below 0.12.
# ----------------------------------------------------------------------------------------------------------------------

terraform {
  # This module is now only being tested with Terraform 0.13.x. However, to make upgrading easier, we are setting
  # 0.12.26 as the minimum version, as that version added support for required_providers with source URLs, making it
  # forwards compatible with 0.13.x code.
  required_version = ">= 0.12.26"
}

# ------------------------------------------------------------------------------
# CREATE SECURITY GROUP RULES OPTIMIZED FOR CHRONOGRAF
# ------------------------------------------------------------------------------

resource "aws_security_group_rule" "http_port_cidr_blocks" {
  count             = length(var.http_port_cidr_blocks) >= 1 ? 1 : 0
  type              = "ingress"
  from_port         = var.http_port
  to_port           = var.http_port
  protocol          = "tcp"
  security_group_id = var.security_group_id
  cidr_blocks       = var.http_port_cidr_blocks
}

resource "aws_security_group_rule" "http_port_security_groups" {
  count                    = var.http_port_security_groups_num == null ? 0 : var.http_port_security_groups_num
  type                     = "ingress"
  from_port                = var.http_port
  to_port                  = var.http_port
  protocol                 = "tcp"
  security_group_id        = var.security_group_id
  source_security_group_id = element(var.http_port_security_groups, count.index)
}
