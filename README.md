# wait4it
A simple go application to test whether a port is ready to accept a connection, or check 
MySQL or PostgresQL server is ready or not, Also you can do http call and check the response code and text in response.  
It also supports **timeout** so it can wait for a particular time and then fail.

## Supported Modules
* TCP port
* MySQL
* PostgresQL
* Http

## Install
You can download the latest [release](https://github.com/ph4r5h4d/wait4it/releases), or you can build it yourself.
To build just run `go build -o wait4it`

## Command Line Args
The following command-line flags are supported

* h (host to check, default is 127.0.0.1)
* p (port to check on the host) 
* t (timeout in seconds, time to wait before considering the operation as failed. default is 30)
* u (username for the services that needs username)
* P (password for the services that needs password)
* n (currently this param is used to identify database name for MySQL)  
* ssl (whether to enable or disable ssl-mode for postgres)  
* http-status (for Http check, which status code to expect)  
* http-text (for Http check, find substring inside the response)  

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
#### Exit codes
* 0: connection established successfully
* 1: timed out
* 2: mostly means a validation error or something wrong with the input data

#### Http check
* for the Http check if you do not define status code it will check for 200 status code
* if `http-text` is not defined then wait4it will just check the status code

#### Postgres check
* if `ssl` is not defined then it's `disable` by default