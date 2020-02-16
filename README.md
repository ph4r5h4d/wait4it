# wait4it
A simple go application to test whether a port is ready to accept a connection, and also you can check whether your 
MySQL server is ready or not.
It also supports **timeout** so it can wait for a particular time and then fail.

## Supported Modules
* TCP port
* MySQL
* PostgresSQL

## Install
You can download the latest release, or you can build it yourself.
To build just run `go build -o wait4it`

## Command Line Args
The following command-line flags are supported

* h (host to check, default is 127.0.0.1)
* p (port to check on the host) 
* t (timeout in seconds, time to wait before considering the operation as failed. default is 30)
* u (username for the services that needs username)
* P (password for the services that needs password)
* n (currently this param is used to identify database name for MySQL)  

### Sample
Check a TCP port  
```bash
./wait4it -h=127.0.0.1 -p=8080 -t=60
``` 

Check a MySQL instance
```bash
./wait4it -h=127.0.0.1 -p=3306 -t=60 -u=root -P=secret -n=app 
```
Check a PostgresSQL instance
```bash
./wait4it -type=postgres -h=127.0.0.1 -p=5432 -t=60 -u=postgres -P=secret -ssl=disable
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

Check a PostgresSQL instance
```bash
docker run ph4r5h4d/wait4it -type=postgres -h=127.0.0.1 -p=5432 -t=60 -u=postgres -P=secret -ssl=disable
```

## Notes
#### Exit codes
* 0: connection established successfully
* 1: timed out
* 2: mostly means a validation error or something wrong with the input data
