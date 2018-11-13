# ---------------------------------------------------------------------------------------------------------------------
# ENVIRONMENT VARIABLES
# Define these secrets as environment variables
# ---------------------------------------------------------------------------------------------------------------------

# AWS_ACCESS_KEY_ID
# AWS_SECRET_ACCESS_KEY

# ---------------------------------------------------------------------------------------------------------------------
# REQUIRED PARAMETERS
# You must provide a value for each of these parameters.
# ---------------------------------------------------------------------------------------------------------------------

variable "aws_region" {
  description = "The AWS region in which all resources will be created"
}

variable "ami_id" {
  description = "The ID of the AMI to run in the cluster. This should be an AMI built from the Packer template under examples/influxdb-ami/influxdb.json."
}

variable "license_key" {
  description = "The key of your InfluxDB Enterprise license. This should not be set in plain-text and can be passed in as an env var or from a secrets management tool."
}

variable "shared_secret" {
  description = "A long pass phrase that will be used to sign tokens for intra-cluster communication on data nodes. This should not be set in plain-text and can be passed in as an env var or from a secrets management tool."
}

# ---------------------------------------------------------------------------------------------------------------------
# OPTIONAL PARAMETERS
# These parameters have reasonable defaults.
# ---------------------------------------------------------------------------------------------------------------------

variable "api_port" {
  description = "The HTTP API port the Data nodes listen on for external communication."
  default     = "8086"
}

variable "influxdb_meta_nodes_cluster_name" {
  description = "What to name the InfluxDB meta nodes cluster and all of its associated resources"
  default     = "influxdb-cluster-meta"
}

variable "influxdb_data_nodes_cluster_name" {
  description = "What to name the InfluxDB data nodes cluster and all of its associated resources"
  default     = "influxdb-cluster-data"
}

variable "ssh_key_name" {
  description = "The name of an EC2 Key Pair that can be used to SSH to the EC2 Instances in this cluster. Set to an empty string to not associate a Key Pair."
  default     = ""
}

variable "meta_volume_device_name" {
  description = "The device name to use for the EBS Volume used for the meta directory on InfluxDB meta nodes."
  default     = "/dev/xvdh"
}

variable "meta_volume_mount_point" {
  description = "The mount point (folder path) to use for the EBS Volume used for the meta directory on InfluxDB meta nodes."
  default     = "/influxdb-meta"
}

variable "data_volume_device_name" {
  description = "The device name to use for the EBS Volume used for the meta, data, wal and hh directories on InfluxDB nodes."
  default     = "/dev/xvdh"
}

variable "data_volume_mount_point" {
  description = "The mount point (folder path) to use for the EBS Volume used for the meta, data, wal and hh directories on InfluxDB data nodes."
  default     = "/influxdb-data"
}

variable "volume_owner" {
  description = "The OS user who should be made the owner of mount points."
  default     = "influxdb"
}
