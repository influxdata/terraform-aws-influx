# Local mocks

This folder contains some mocks that make it possible to run and test InfluxDB locally, using Docker and Docker 
Compose. The mocks here replace external dependencies, such as EC2 metadata and AWS API calls, with mocks that work 
locally. This is solely to make testing and iterating on the code faster and easier and should NOT be used in 
production!

## Quick start

First, use the Packer template in the [influxdb-ami 
example](https://github.com/gruntwork-io/terraform-aws-influx/tree/master/examples/influxdb-ami) to build a Docker 
image with InfluxDB installed on Ubuntu:

```
packer build -only=influxdb-docker-ubuntu influxdb.json
```

To run the Docker image, head into one of the `examples/influxdb-xxx/local-test` folders and run:

```
docker-compose up
```

Alternatively, to create an Amazon Linux 2 Docker image:

```
packer build -only=influxdb-docker-amazon-linux influxdb.json
```

And to run it:

```
OS_NAME=amazon-linux docker-compose up
```

Wait 10-15 seconds and then run the `influx` CLI without any arguments to connect! No
credentials are required to login.

