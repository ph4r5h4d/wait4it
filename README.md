# wait4it
 
![DOCKER RELEASE](https://img.shields.io/github/actions/workflow/status/ph4r5h4d/wait4it/main.yml?label=Docker%20release&style=for-the-badge) ![TESTS](https://img.shields.io/github/actions/workflow/status/ph4r5h4d/wait4it/integrations-app.yaml?label=Tests&style=for-the-badge) ![DOCKERTESTS](https://img.shields.io/github/actions/workflow/status/ph4r5h4d/wait4it/integrations-docker.yaml?label=Docker%20Tests&style=for-the-badge) ![SECURITY SCAN](https://img.shields.io/github/actions/workflow/status/ph4r5h4d/wait4it/trivy-scan.yml?label=Security%20Scan&style=for-the-badge) [![Docker Pull](https://img.shields.io/docker/pulls/ph4r5h4d/wait4it?style=for-the-badge)](https://hub.docker.com/r/ph4r5h4d/wait4it)  ![GO Version](https://img.shields.io/github/go-mod/go-version/ph4r5h4d/wait4it?style=for-the-badge) ![TAG](https://img.shields.io/github/v/tag/ph4r5h4d/wait4it?style=for-the-badge) ![LICENSE](https://img.shields.io/github/license/ph4r5h4d/wait4it?style=for-the-badge) ![GitHub Downloads (all assets, all releases)](https://img.shields.io/github/downloads/ph4r5h4d/wait4it/total?style=for-the-badge&label=Github%20Downloads)

A simple go application to test whether a port is ready to accept a connection or check 
MySQL, PostgreSQL, MongoDB or Redis server is ready or not, Also you can do Http call and check 
the response code and text in response.  
It also supports **timeout** so it can wait for a particular time and then fail.  

## Supported Services
* [TCP port](https://wait4it.dev/docs/tcp/)
* [MySQL](https://wait4it.dev/docs/mysql/)
* [PostgresQL](https://wait4it.dev/docs/postgresql/)
* [Http](https://wait4it.dev/docs/http/)
* [MongoDB](https://wait4it.dev/docs/mongodb/)
* [Redis](https://wait4it.dev/docs/redis/)
* [RabbitMQ](https://wait4it.dev/docs/rabbitmq/)
* [Memcached](https://wait4it.dev/docs/memcached/)
* [ElasticSearch](https://wait4it.dev/docs/elasticsearch/)
* [Aerospike](https://wait4it.dev/docs/aerospike/)
* [Kafka](https://wait4it.dev/docs/kafka/)

## Install
You can download the latest [release](https://github.com/ph4r5h4d/wait4it/releases), or you can build it yourself.
To build just run `go build`.
For detailed installation instructions, visit the [installation doc](https://wait4it.dev/docs/installation/).

## Documentation
Visit the [website](https://wait4it.dev) for detailed documentation.

## Security
This project includes automated security scanning using [Trivy](https://trivy.dev/) to ensure the Docker images are free from known vulnerabilities. Security scans are performed:
- On every release
- Weekly on a scheduled basis
- Can be triggered manually

Scan results are available in the repository's Security tab and as downloadable artifacts from the workflow runs.

## Powered by
[![GoLand logo.](https://resources.jetbrains.com/storage/products/company/brand/logos/GoLand.svg)](https://jb.gg/OpenSourceSupport)
