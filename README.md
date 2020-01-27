# wait4it
A simple go application to test whether a port is ready to accept connection or not.
It also supports **timeout** so it can wait for a certain time and then fail.

## Install
You can download the latest release or you can build it yourself.
To build just run `go build -o wait4it`

## Usage
The following command line flags are supported

* h (host to check, default is 127.0.0.1)
* p (port to check on the host) 
* t (timeout in seconds, time to wait before considering the opration as failed. default is 30)

Sample
`./wait4it -h=127.0.0.1 -p=8080 -t=60 `