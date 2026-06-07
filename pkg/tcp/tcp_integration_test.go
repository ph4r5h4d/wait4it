package tcp

import (
	"context"
	"net"
	"testing"
)

func TestCheck_SuccessfulConnection(t *testing.T) {
	// Start a TCP listener on a random port
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to listen: %v", err)
	}
	defer ln.Close()

	addr := ln.Addr().(*net.TCPAddr)
	c := &Check{
		Addr: "127.0.0.1",
		Port: addr.Port,
	}
	if err := c.Validate(); err != nil {
		t.Fatalf("Validate() returned unexpected error: %v", err)
	}

	ok, endOnError, err := c.Check(context.Background())
	if !ok {
		t.Errorf("Check() ok = false, want true; err = %v", err)
	}
	if endOnError {
		t.Error("Check() endOnError = true, want false")
	}
	if err != nil {
		t.Errorf("Check() returned unexpected error: %v", err)
	}
}

func TestCheck_ConnectionRefused(t *testing.T) {
	// Listen and immediately close to get a free port, then try to connect
	// The port should be free and connection should be refused
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to listen: %v", err)
	}
	addr := ln.Addr().(*net.TCPAddr)
	port := addr.Port
	ln.Close()

	c := &Check{
		Addr: "127.0.0.1",
		Port: port,
	}
	if err := c.Validate(); err != nil {
		t.Fatalf("Validate() returned unexpected error: %v", err)
	}

	ok, endOnError, err := c.Check(context.Background())
	if ok {
		t.Error("Check() ok = true, want false for connection refused")
	}
	if endOnError {
		t.Error("Check() endOnError = true, want false for connection refused")
	}
	if err == nil {
		t.Error("Check() should return error for connection refused")
	}
}

func TestCheck_ContextCancellation(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to listen: %v", err)
	}
	defer ln.Close()

	addr := ln.Addr().(*net.TCPAddr)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	c := &Check{
		Addr: "127.0.0.1",
		Port: addr.Port,
	}
	if err := c.Validate(); err != nil {
		t.Fatalf("Validate() returned unexpected error: %v", err)
	}

	ok, _, err := c.Check(ctx)
	if ok {
		t.Error("Check() ok = true, want false for cancelled context")
	}
	if err == nil {
		t.Error("Check() should return error for cancelled context")
	}
}

func TestCheck_AcceptsConnection(t *testing.T) {
	// Start a server that accepts connections
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to listen: %v", err)
	}
	defer ln.Close()

	// Accept connections in background
	go func() {
		conn, err := ln.Accept()
		if err == nil {
			_ = conn.Close()
		}
	}()

	addr := ln.Addr().(*net.TCPAddr)
	c := &Check{
		Addr: "127.0.0.1",
		Port: addr.Port,
	}
	if err := c.Validate(); err != nil {
		t.Fatalf("Validate() returned unexpected error: %v", err)
	}

	ok, endOnError, err := c.Check(context.Background())
	if !ok {
		t.Errorf("Check() ok = false, want true; err = %v", err)
	}
	if endOnError {
		t.Error("Check() endOnError = true, want false")
	}
	if err != nil {
		t.Errorf("Check() returned unexpected error: %v", err)
	}
}

func TestCheck_InvalidPort(t *testing.T) {
	c := &Check{
		Addr: "127.0.0.1",
		Port: 99999,
	}
	if err := c.Validate(); err == nil {
		t.Error("Validate() should return error for invalid port")
	}
}

func TestCheck_CustomDialer(t *testing.T) {
	// Test that a custom dialer can be injected
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to listen: %v", err)
	}
	defer ln.Close()

	// Accept connections in background
	go func() {
		conn, err := ln.Accept()
		if err == nil {
			_ = conn.Close()
		}
	}()

	addr := ln.Addr().(*net.TCPAddr)
	c := &Check{
		Addr:   "127.0.0.1",
		Port:   addr.Port,
		dialer: &net.Dialer{},
	}
	if err := c.Validate(); err != nil {
		t.Fatalf("Validate() returned unexpected error: %v", err)
	}

	ok, _, err := c.Check(context.Background())
	if !ok {
		t.Errorf("Check() ok = false, want true with custom dialer; err = %v", err)
	}
}