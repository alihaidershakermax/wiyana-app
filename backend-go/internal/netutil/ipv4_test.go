package netutil

import (
	"context"
	"net"
	"testing"
	"time"
)

func TestIpv4Dialer(t *testing.T) {
	// dns.google publishes AAAA first; dialer must reach it over IPv4.
	start := time.Now()
	conn, err := Ipv4Dialer().DialContext(context.Background(), "tcp", "dns.google:443")
	elapsed := time.Since(start)
	if err != nil {
		t.Fatalf("dial failed: %v (took %v)", err, elapsed)
	}
	conn.Close()
	t.Logf("dial OK in %v", elapsed)
}

func TestHostPortIPv4(t *testing.T) {
	addr, err := HostPortIPv4(context.Background(), "dns.google:443")
	if err != nil {
		t.Fatalf("HostPortIPv4 failed: %v", err)
	}
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		t.Fatalf("split failed: %v", err)
	}
	if net.ParseIP(host).To4() == nil {
		t.Fatalf("host is not IPv4: %s", host)
	}
	if port != "443" {
		t.Fatalf("port changed: %s", port)
	}
	t.Logf("resolved to %s", addr)
}
