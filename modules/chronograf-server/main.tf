# ---------------------------------------------------------------------------------------------------------------------
# CREATE AN EC2 INSTANCE TO RUN CHRONOGRAF
# ---------------------------------------------------------------------------------------------------------------------

resource "aws_instance" "chronograf_server" {
  ami                    = "${var.ami_id}"
  instance_type          = "${var.instance_type}"
  user_data              = "${var.user_data}"
  key_name               = "${var.ssh_key_name}"
  vpc_security_group_ids = ["${aws_security_group.chronograf_security_group.id}"]
  tags                   = "${var.tags}"
}

# ---------------------------------------------------------------------------------------------------------------------
# CREATE A SECURITY GROUP TO CONTROL WHAT REQUESTS CAN GO IN AND OUT OF THE EC2 INSTANCE
# We export the ID of the security group as an output variable so users can attach custom rules.
# ---------------------------------------------------------------------------------------------------------------------

resource "aws_security_group" "chronograf_security_group" {
  name_prefix = "chronograf"
  description = "Security group for the Chronograf server"
  vpc_id      = "${var.vpc_id}"
}

# ---------------------------------------------------------------------------------------------------------------------
# ATTACH DEFAULT SECURITY GROUP RULES TO ALLOW SSH ACCESS AND ALL OUTBOUND CONNECTIONS
# ---------------------------------------------------------------------------------------------------------------------

resource "aws_security_group_rule" "allow_ssh_inbound" {
  count       = "${length(var.allowed_ssh_cidr_blocks) >= 1 ? 1 : 0}"
  type        = "ingress"
  from_port   = "${var.ssh_port}"
  to_port     = "${var.ssh_port}"
  protocol    = "tcp"
  cidr_blocks = ["${var.allowed_ssh_cidr_blocks}"]

  security_group_id = "${aws_security_group.chronograf_security_group.id}"
}

resource "aws_security_group_rule" "allow_ssh_inbound_from_security_group_ids" {
  count                    = "${var.allowed_ssh_security_group_ids_num}"
  type                     = "ingress"
  from_port                = "${var.ssh_port}"
  to_port                  = "${var.ssh_port}"
  protocol                 = "tcp"
  source_security_group_id = "${element(var.allowed_ssh_security_group_ids, count.index)}"

  security_group_id = "${aws_security_group.chronograf_security_group.id}"
}

resource "aws_security_group_rule" "allow_all_outbound" {
  type        = "egress"
  from_port   = 0
  to_port     = 0
  protocol    = "-1"
  cidr_blocks = ["0.0.0.0/0"]

  security_group_id = "${aws_security_group.chronograf_security_group.id}"
}
