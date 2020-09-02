# ----------------------------------------------------------------------------------------------------------------------
# REQUIRE A SPECIFIC TERRAFORM VERSION OR HIGHER
# This module has been updated with 0.12 syntax, which means it is no longer compatible with any versions below 0.12.
# ----------------------------------------------------------------------------------------------------------------------

terraform {
  # Source is required for required_providers in TF 13, but is only compatible with TF 12 starting 0.12.26.
  required_version = ">= 0.12.26"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 2.6"
    }
  }
}

# ---------------------------------------------------------------------------------------------------------------------
# CREATE A TARGET GROUP
# This will perform health checks on the servers and receive requests from the Listerers that match Listener Rules.
# ---------------------------------------------------------------------------------------------------------------------

resource "aws_alb_target_group" "tg" {
  name                 = var.target_group_name
  port                 = var.port
  protocol             = var.protocol
  vpc_id               = var.vpc_id
  deregistration_delay = var.deregistration_delay

  health_check {
    port                = "traffic-port"
    protocol            = var.protocol
    interval            = var.health_check_interval
    path                = var.health_check_path
    timeout             = var.health_check_timeout
    healthy_threshold   = var.health_check_healthy_threshold
    unhealthy_threshold = var.health_check_unhealthy_threshold
    matcher             = var.health_check_matcher
  }
}

# ---------------------------------------------------------------------------------------------------------------------
# CREATE LISTENER RULES
# These rules determine which requests get routed to the Target Group
# ---------------------------------------------------------------------------------------------------------------------

resource "aws_alb_listener_rule" "http_path" {
  count = var.listener_arns_num

  listener_arn = element(var.listener_arns, count.index)
  priority     = var.listener_rule_starting_priority + count.index

  action {
    type             = "forward"
    target_group_arn = aws_alb_target_group.tg.arn
  }

  # For backwards compatibility and to support the input format, we translate the old field and values pattern to the
  # nested subblock pattern in AWS provider v3. Only one of the sub blocks will be included based on the value of the
  # `field` attribute of `routing_condition`.
  dynamic "condition" {
    for_each = var.routing_condition != null ? [var.routing_condition] : []

    content {
      dynamic "host_header" {
        for_each = condition.value.field == "host-header" ? [var.routing_condition] : []
        content {
          values = host_header.value.values
        }
      }

      dynamic "http_request_method" {
        for_each = condition.value.field == "http-request-method" ? [var.routing_condition] : []
        content {
          values = http_request_method.value.values
        }
      }

      dynamic "path_pattern" {
        for_each = condition.value.field == "path-pattern" ? [var.routing_condition] : []
        content {
          values = path_pattern.value.values
        }
      }

      dynamic "source_ip" {
        for_each = condition.value.field == "source-ip" ? [var.routing_condition] : []
        content {
          values = source_ip.value.values
        }
      }
    }
  }
}

# ---------------------------------------------------------------------------------------------------------------------
# ATTACH THE AUTO SCALING GROUP (ASG) TO THE LOAD BALANCER
# As a result, each EC2 Instance in the ASG will register with the Load Balancer, go through health checks, and be
# replaced automatically if it starts failing health checks.
# ---------------------------------------------------------------------------------------------------------------------

resource "aws_autoscaling_attachment" "attach" {
  autoscaling_group_name = var.asg_name
  alb_target_group_arn   = aws_alb_target_group.tg.arn
}
