# Kapacitor Run Script

This folder contains a script for configuring and initializing Kapacitor on an [AWS](https://aws.amazon.com/) server. 
This script has been tested on the following operating systems:

* Ubuntu 18.04
* Amazon Linux 2

There is a good chance it will work on other flavors of Debian, CentOS, and RHEL as well.

## Quick start

This script assumes you installed it, plus all of its dependencies (including Kapacitor itself), using the 
[install-kapacitor module](https://github.com/gruntwork-io/terraform-aws-influx/tree/master/modules/install-kapacitor). 

This will:

1. Fill out the templated configuration file with user supplied values.

1. Start Kapacitor on the local node.

1. Wait for all Kapacitor ASG to to spin up all desired instances then update `/etc/hosts` with the IPs of all instances.
   the value of the instances' `Name` tag is used as the `hostname` entry.

1. Figure out a rally point for your Kapacitor cluster. This is a "leader" Meta node that will be responsible for 
   initializing the cluster. See [Picking a rally point](#picking-a-rally-point) for more info.

1. On the rally point, initialize the cluster, including adding all nodes to the cluster

We recommend using the `run-kapacitor` command as part of [User 
Data](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/user-data.html#user-data-shell-scripts), so that it executes
when the EC2 Instance is first booting.

See the [examples folder](https://github.com/gruntwork-io/terraform-aws-influx/tree/master/examples) for 
fully-working sample code.

## Command line Arguments

Run `run-kapacitor --help` to see all available arguments.

```
Usage: run-kapacitor [options]

This script can be used to configure and initialize Kapacitor. This script has been tested with Ubuntu 18.04 and Amazon Linux 2.

Options:

  --hostname    The hostname of the current node.
  --asg-name    The name of the ASG that contains meta nodes.
  --region      The AWS region the Auto Scaling Groups are deployed in.
  --auto-fill   Search the Kapacitor config file for KEY and replace it with VALUE. May be repeated.

Example:

  run-kapacitor --asg-name kapacitor-asg --region us-east-1 --auto-fill '<__LICENSE_KEY__>=******'
```

## Picking a rally point

The Kapacitor cluster needs a "rally point", which is a single node that is responsible for:

1. Initializing the cluster.
1. Adding/removing nodes to the cluster.

We need a way to unambiguously and reliably select exactly one rally point. If there's more than one node, you may end
up with multiple separate clusters instead of just one!

The `run-kapacitor` script can automatically pick a rally point automatically by:

1. Looking up all the servers in the Auto Scaling Group specified via the `--asg-name` parameter.

1. Pick the node with the oldest Launch Time as the rally point. If multiple nodes have identical launch times, use the
   one with the earliest Instance ID, alphabetically.

## Passing credentials securely

The `run-kapacitor` script requires that you pass in your license key and shared secret. You should make sure to never 
store these credentials in plaintext! You should use a secrets management tool to store the credentials in an encrypted
format and only decrypt them, in memory, just before calling `run-kapacitor`. Here are some tools to consider:

* [Vault](https://www.vaultproject.io/)
* [Keywhiz](https://square.github.io/keywhiz/)
* [KMS](https://aws.amazon.com/kms/)

Moreover, if you're ever calling `run-kapacitor` interactively (i.e., you're manually running CLI commands
rather than executing a script), be careful of passing credentials directly on the command line, or they will be 
stored, in plaintext, [in Bash 
history](https://www.digitalocean.com/community/tutorials/how-to-use-bash-history-commands-and-expansions-on-a-linux-vps)!
You can either use a CLI tool to set the credentials as environment variables or you can [temporarily disable Bash
history](https://linuxconfig.org/how-to-disable-bash-shell-commands-history-on-linux). 

## Required permissions

The `run-kapacitor` script assumes it is running on an EC2 Instance with an [IAM 
Role](http://docs.aws.amazon.com/IAM/latest/UserGuide/id_roles.html) that has the following permissions:

* `ec2:DescribeInstances`
* `ec2:DescribeTags`
* `autoscaling:DescribeAutoScalingGroups`

These permissions are automatically added by the [kapacitor-cluster 
module](https://github.com/gruntwork-io/terraform-aws-influx/tree/master/modules/kapacitor-cluster).
