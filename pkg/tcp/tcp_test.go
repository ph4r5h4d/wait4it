package tcp

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"time"
)

func Test_check_Check(t *testing.T) {
	t.Run("free port", func(t *testing.T) {
		var err error

		tcp := &check{
			addr: "localhost",
		}

		tcp.port, err = GetFreePort()
		if err != nil {
			panic(fmt.Sprintf("no free port available for testing: %s", err.Error()))
		}

		fmt.Printf("TCP testing with addr %s port %d\n", tcp.addr, tcp.port)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		got, got1, err := tcp.Check(ctx)
		if (err != nil) != true {
			t.Errorf("Check() error = %v, wantErr %v", err, true)
			cancel()
			return
		}
		cancel()

		if got != false {
			t.Errorf("Check() got = %v, want %v", got, false)
		}

		if got1 != false {
			t.Errorf("Check() got1 = %v, want %v", got1, false)
		}
	})

	t.Run("active port", func(t *testing.T) {
		var s *httptest.Server
		var err error
		tcp := &check{}

		// run a test server
		s = httptest.NewServer(
			http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			}))

		tcp.addr, tcp.port, err = getHostAndPort(s.URL)
		if err != nil {
			panic(err)
		}

		fmt.Printf("TCP testing with addr %s port %d\n", tcp.addr, tcp.port)

		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		got, got1, err := tcp.Check(ctx)
		if (err != nil) != false {
			t.Errorf("Check() error = %v, wantErr %v", err, false)
			cancel()
			return
		}
		cancel()

		if got != true {
			t.Errorf("Check() got = %v, want %v", got, true)
		}

		if got1 != false {
			t.Errorf("Check() got1 = %v, want %v", got1, false)
		}

		s.Close()
	})
}

func getHostAndPort(str string) (string, int, error) {
	u, err := url.Parse(str)
	if err != nil {
		return "", 0, err
	}

	host, port, err := net.SplitHostPort(u.Host)
	if err != nil {
		return "", 0, err
	}

	portInt, err := strconv.ParseInt(port, 10, 32)
	if err != nil {
		return "", 0, err
	}
	return host, int(portInt), nil

}
func Test_check_validate(t *testing.T) {
	type fields struct {
		addr string
		port int
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid sample ",
			fields: fields{
				port: 8080,
			},
			wantErr: false,
		},
		{
			name: "minimum",
			fields: fields{
				port: 1,
			},
			wantErr: false,
		},
		{
			name: "maximum",
			fields: fields{
				port: 65535,
			},
			wantErr: false,
		},
		{
			name: "out of range",
			fields: fields{
				port: 66000,
			},
			wantErr: true,
		},
		{
			name: "greater than max",
			fields: fields{
				port: maxPort + 1,
			},
			wantErr: true,
		},
		{
			name: "less than min",
			fields: fields{
				port: minPort - 1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tcp := &check{
				addr: tt.fields.addr,
				port: tt.fields.port,
			}
			if err := tcp.validate(); (err != nil) != tt.wantErr {
				t.Errorf("validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// GetFreePort asks the kernel for a free open port that is ready to use.
func GetFreePort() (port int, err error) {
	var a *net.TCPAddr
	if a, err = net.ResolveTCPAddr("tcp", "localhost:0"); err == nil {
		var l *net.TCPListener
		if l, err = net.ListenTCP("tcp", a); err == nil {
			defer l.Close()
			return l.Addr().(*net.TCPAddr).Port, nil
		}
	}
	return
}
