package dns

import (
	"context"
	"errors"
	"net"
	"testing"
)

// mockResolver implements the Resolver interface for testing
type mockResolver struct {
	lookupIPFunc    func(ctx context.Context, network, host string) ([]net.IP, error)
	lookupCNAMEFunc func(ctx context.Context, host string) (string, error)
	lookupMXFunc    func(ctx context.Context, host string) ([]*net.MX, error)
	lookupTXTFunc   func(ctx context.Context, host string) ([]string, error)
	lookupSRVFunc   func(ctx context.Context, service, proto, name string) (string, []*net.SRV, error)
	lookupNSFunc    func(ctx context.Context, host string) ([]*net.NS, error)
	lookupAddrFunc  func(ctx context.Context, addr string) ([]string, error)
}

func (m *mockResolver) LookupIP(ctx context.Context, network, host string) ([]net.IP, error) {
	if m.lookupIPFunc != nil {
		return m.lookupIPFunc(ctx, network, host)
	}
	return nil, errors.New("not implemented")
}

func (m *mockResolver) LookupCNAME(ctx context.Context, host string) (string, error) {
	if m.lookupCNAMEFunc != nil {
		return m.lookupCNAMEFunc(ctx, host)
	}
	return "", errors.New("not implemented")
}

func (m *mockResolver) LookupMX(ctx context.Context, host string) ([]*net.MX, error) {
	if m.lookupMXFunc != nil {
		return m.lookupMXFunc(ctx, host)
	}
	return nil, errors.New("not implemented")
}

func (m *mockResolver) LookupTXT(ctx context.Context, host string) ([]string, error) {
	if m.lookupTXTFunc != nil {
		return m.lookupTXTFunc(ctx, host)
	}
	return nil, errors.New("not implemented")
}

func (m *mockResolver) LookupSRV(ctx context.Context, service, proto, name string) (string, []*net.SRV, error) {
	if m.lookupSRVFunc != nil {
		return m.lookupSRVFunc(ctx, service, proto, name)
	}
	return "", nil, errors.New("not implemented")
}

func (m *mockResolver) LookupNS(ctx context.Context, host string) ([]*net.NS, error) {
	if m.lookupNSFunc != nil {
		return m.lookupNSFunc(ctx, host)
	}
	return nil, errors.New("not implemented")
}

func (m *mockResolver) LookupAddr(ctx context.Context, addr string) ([]string, error) {
	if m.lookupAddrFunc != nil {
		return m.lookupAddrFunc(ctx, addr)
	}
	return nil, errors.New("not implemented")
}

func TestCheck_ARecord_Success(t *testing.T) {
	mock := &mockResolver{
		lookupIPFunc: func(ctx context.Context, network, host string) ([]net.IP, error) {
			return []net.IP{net.ParseIP("1.2.3.4")}, nil
		},
	}

	d := &Check{
		Host:       "example.com",
		RecordType: "A",
		resolver:   mock,
	}

	ok, endOnError, err := d.Check(context.Background())
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

func TestCheck_AAAARecord_Success(t *testing.T) {
	mock := &mockResolver{
		lookupIPFunc: func(ctx context.Context, network, host string) ([]net.IP, error) {
			return []net.IP{net.ParseIP("::1")}, nil
		},
	}

	d := &Check{
		Host:       "example.com",
		RecordType: "AAAA",
		resolver:   mock,
	}

	ok, _, err := d.Check(context.Background())
	if !ok {
		t.Errorf("Check() ok = false, want true; err = %v", err)
	}
}

func TestCheck_CNAMERecord_Success(t *testing.T) {
	mock := &mockResolver{
		lookupCNAMEFunc: func(ctx context.Context, host string) (string, error) {
			return "cdn.example.com.", nil
		},
	}

	d := &Check{
		Host:       "www.example.com",
		RecordType: "CNAME",
		resolver:   mock,
	}

	ok, _, err := d.Check(context.Background())
	if !ok {
		t.Errorf("Check() ok = false, want true; err = %v", err)
	}
}

func TestCheck_MXRecord_Success(t *testing.T) {
	mock := &mockResolver{
		lookupMXFunc: func(ctx context.Context, host string) ([]*net.MX, error) {
			return []*net.MX{{Host: "mail.example.com."}}, nil
		},
	}

	d := &Check{
		Host:       "example.com",
		RecordType: "MX",
		resolver:   mock,
	}

	ok, _, err := d.Check(context.Background())
	if !ok {
		t.Errorf("Check() ok = false, want true; err = %v", err)
	}
}

func TestCheck_TXTRecord_Success(t *testing.T) {
	mock := &mockResolver{
		lookupTXTFunc: func(ctx context.Context, host string) ([]string, error) {
			return []string{"v=spf1 include:example.com ~all"}, nil
		},
	}

	d := &Check{
		Host:       "example.com",
		RecordType: "TXT",
		resolver:   mock,
	}

	ok, _, err := d.Check(context.Background())
	if !ok {
		t.Errorf("Check() ok = false, want true; err = %v", err)
	}
}

func TestCheck_NSRecord_Success(t *testing.T) {
	mock := &mockResolver{
		lookupNSFunc: func(ctx context.Context, host string) ([]*net.NS, error) {
			return []*net.NS{{Host: "ns1.example.com."}}, nil
		},
	}

	d := &Check{
		Host:       "example.com",
		RecordType: "NS",
		resolver:   mock,
	}

	ok, _, err := d.Check(context.Background())
	if !ok {
		t.Errorf("Check() ok = false, want true; err = %v", err)
	}
}

func TestCheck_PTRRecord_Success(t *testing.T) {
	mock := &mockResolver{
		lookupAddrFunc: func(ctx context.Context, addr string) ([]string, error) {
			return []string{"host.example.com."}, nil
		},
	}

	d := &Check{
		Host:       "1.2.3.4",
		RecordType: "PTR",
		resolver:   mock,
	}

	ok, _, err := d.Check(context.Background())
	if !ok {
		t.Errorf("Check() ok = false, want true; err = %v", err)
	}
}

func TestCheck_SRVRecord_Success(t *testing.T) {
	mock := &mockResolver{
		lookupSRVFunc: func(ctx context.Context, service, proto, name string) (string, []*net.SRV, error) {
			return "", []*net.SRV{{Target: "server.example.com."}}, nil
		},
	}

	d := &Check{
		Host:       "_sip._tcp.example.com",
		RecordType: "SRV",
		resolver:   mock,
	}

	ok, _, err := d.Check(context.Background())
	if !ok {
		t.Errorf("Check() ok = false, want true; err = %v", err)
	}
}

func TestCheck_ExpectedValueMatch(t *testing.T) {
	mock := &mockResolver{
		lookupIPFunc: func(ctx context.Context, network, host string) ([]net.IP, error) {
			return []net.IP{net.ParseIP("1.2.3.4")}, nil
		},
	}

	d := &Check{
		Host:       "example.com",
		RecordType: "A",
		Expected:   "1.2.3.4",
		resolver:   mock,
	}

	ok, _, err := d.Check(context.Background())
	if !ok {
		t.Errorf("Check() ok = false, want true; err = %v", err)
	}
}

func TestCheck_ExpectedValueMismatch(t *testing.T) {
	mock := &mockResolver{
		lookupIPFunc: func(ctx context.Context, network, host string) ([]net.IP, error) {
			return []net.IP{net.ParseIP("1.2.3.4")}, nil
		},
	}

	d := &Check{
		Host:       "example.com",
		RecordType: "A",
		Expected:   "5.6.7.8",
		resolver:   mock,
	}

	ok, _, err := d.Check(context.Background())
	if ok {
		t.Error("Check() ok = true, want false for expected value mismatch")
	}
	if err != ErrExpectedNotFound {
		t.Errorf("Check() err = %v, want ErrExpectedNotFound", err)
	}
}

func TestCheck_NoRecords(t *testing.T) {
	mock := &mockResolver{
		lookupIPFunc: func(ctx context.Context, network, host string) ([]net.IP, error) {
			return []net.IP{}, nil
		},
	}

	d := &Check{
		Host:       "example.com",
		RecordType: "A",
		resolver:   mock,
	}

	ok, _, err := d.Check(context.Background())
	if ok {
		t.Error("Check() ok = true, want false for no records")
	}
	if err != ErrNoRecords {
		t.Errorf("Check() err = %v, want ErrNoRecords", err)
	}
}

func TestCheck_ResolverError(t *testing.T) {
	dnsErr := errors.New("DNS lookup failed")
	mock := &mockResolver{
		lookupIPFunc: func(ctx context.Context, network, host string) ([]net.IP, error) {
			return nil, dnsErr
		},
	}

	d := &Check{
		Host:       "example.com",
		RecordType: "A",
		resolver:   mock,
	}

	ok, endOnError, err := d.Check(context.Background())
	if ok {
		t.Error("Check() ok = true, want false for resolver error")
	}
	if endOnError {
		t.Error("Check() endOnError = true, want false for DNS errors")
	}
	if err == nil {
		t.Error("Check() should return error for resolver failure")
	}
}

func TestCheck_ContextCancellation(t *testing.T) {
	mock := &mockResolver{
		lookupIPFunc: func(ctx context.Context, network, host string) ([]net.IP, error) {
			return nil, ctx.Err()
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	d := &Check{
		Host:       "example.com",
		RecordType: "A",
		resolver:   mock,
	}

	ok, _, err := d.Check(ctx)
	if ok {
		t.Error("Check() ok = true, want false for cancelled context")
	}
	if err == nil {
		t.Error("Check() should return error for cancelled context")
	}
}

func TestCheck_DefaultResolverUsed(t *testing.T) {
	// When no custom resolver is set, getResolver() should return net.DefaultResolver
	d := &Check{
		Host:       "example.com",
		RecordType: "A",
	}

	resolver := d.getResolver()
	if resolver != net.DefaultResolver {
		t.Error("getResolver() should return net.DefaultResolver when Server is empty")
	}
}

func TestCheck_CustomServerResolver(t *testing.T) {
	// When a custom server is set, getResolver() should return a custom resolver
	d := &Check{
		Host:       "example.com",
		RecordType: "A",
		Server:     "8.8.8.8:53",
	}

	resolver := d.getResolver()
	if resolver == net.DefaultResolver {
		t.Error("getResolver() should return custom resolver when Server is set")
	}
}

func TestCheck_PartialMatch(t *testing.T) {
	mock := &mockResolver{
		lookupTXTFunc: func(ctx context.Context, host string) ([]string, error) {
			return []string{"v=spf1 include:example.com ~all"}, nil
		},
	}

	d := &Check{
		Host:       "example.com",
		RecordType: "TXT",
		Expected:   "spf1",
		resolver:   mock,
	}

	ok, _, err := d.Check(context.Background())
	if !ok {
		t.Errorf("Check() ok = false, want true for partial match; err = %v", err)
	}
}