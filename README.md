# wait4it
 
![LICENSE](https://img.shields.io/github/workflow/status/ph4r5h4d/wait4it/Wait4it&#32;CI?style=for-the-badge) [![Docker Pull](https://img.shields.io/docker/pulls/ph4r5h4d/wait4it?style=for-the-badge)](https://hub.docker.com/r/ph4r5h4d/wait4it)  ![GO Version](https://img.shields.io/github/go-mod/go-version/ph4r5h4d/wait4it?style=for-the-badge) ![TAG](https://img.shields.io/github/v/tag/ph4r5h4d/wait4it?style=for-the-badge) ![LICENSE](https://img.shields.io/github/license/ph4r5h4d/wait4it?style=for-the-badge)
  
A simple go application to test whether a port is ready to accept a connection or check 
MySQL or PostgreSQL server is ready or not, Also you can do Http call and check the response code and text in response.  
It also supports **timeout** so it can wait for a particular time and then fail.

## Supported Modules
* TCP port
* MySQL
* PostgresQL
* Http

## Install
You can download the latest [release](https://github.com/ph4r5h4d/wait4it/releases), or you can build it yourself.
To build just run `go build -o wait4it`

## Configuration
### Environment variables
The following environment variables are supported:

* `W4IT_TYPE` (check type [http, tcp, mysql, postgres])
* `W4IT_TIMEOUT` (timeout in seconds)
* `W4IT_HOST` (the host to check)
* `W4IT_PORT` (port to check on the host)
* `W4IT_USERNAME` (username for the services that needs username)
* `W4IT_PASSWORD` (password for the services that needs password)
* `W4IT_DBNAME` (database name for MySQL or PostgreSQL)
* `W4IT_SSL_MODE` (whether to enable or disable ssl-mode for Postgres [disable, enable])
* `W4IT_HTTP_STATUS_CODE` (for Http check, which status code to expect)
* `W4IT_HTTP_TEXT` (for Http check, find substring inside the response)

### Command Line Args
The following command-line flags are supported

* `-type` (check type [http, tcp, mysql, postgres])
* `-h` (the host to check, default is 127.0.0.1)
* `-p` (port to check on the host) 
* `-t` (timeout in seconds, time to wait before considering the operation as failed. default is 30)
* `-u` (username for the services that needs username)
* `-P` (password for the services that needs password)
* `-n` (database name for MySQL or PostgreSQL)  
* `-ssl` (whether to enable or disable ssl-mode for Postgres [disable, enable])  
* `-http-status` (for Http check, which status code to expect)  
* `-http-text` (for Http check, find substring inside the response)  

### Sample
Check a TCP port  
```bash
./wait4it -h=127.0.0.1 -p=8080 -t=60
``` 

Check a MySQL instance
```bash
./wait4it -h=127.0.0.1 -p=3306 -t=60 -u=root -P=secret -n=app 
```
Check a PostgresQL instance
```bash
./wait4it -type=postgres -h=127.0.0.1 -p=5432 -t=60 -u=postgres -P=secret -ssl=disable
```

Check HTTP response and text
```bash
-type=http -h=https://farshad.nematdoust.com -t=60 -status-code=200 -http-text="Software Engineer" 
```

### Docker
You can run this `wait4it` inside a docker container, and it's possible to run this container as init container inside
K8s and Openshift.  
You can use the following image in Docker Hub:  
`docker pull ph4r5h4d/wait4it`

Inside container you can do the following:  

#### Sample
Check a TCP port  
```bash
docker run ph4r5h4d/wait4it -h=127.0.0.1 -p=8080 -t=60
``` 

Check a MySQL instance
```bash
docker run ph4r5h4d/wait4it -h=127.0.0.1 -p=3306 -t=60 -u=root -P=secret -n=app 
```

Check a PostgresQL instance
```bash
docker run ph4r5h4d/wait4it -type=postgres -h=127.0.0.1 -p=5432 -t=60 -u=postgres -P=secret -ssl=disable
```

Check HTTP response and text
```bash
docker run ph4r5h4d/wait4it -type=http -h=https://farshad.nematdoust.com -t=60 -status-code=200 -http-text="Software Engineer" 
```

## Notes
#### Configuration
* note that environment variables have higher priority than command line arguments. This means if you define both `W4IT_TYPE` and `-type`, the application take 
value of environment variable into account.
#### Exit codes
* 0: connection established successfully
* 1: timed out
* 2: mostly means a validation error or something wrong with the input data

#### Http check
* for the Http check if you do not define status code it will check for 200 status code
* if `http-text` is not defined then wait4it will check the status code

#### Postgres check
* if `ssl` is not defined then it's `disable` by default