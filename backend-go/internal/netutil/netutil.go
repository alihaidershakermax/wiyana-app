package netutil

import (
	"context"
	"errors"
	"net"
	"time"
)

// Ipv4Dialer returns a net.Dialer that resolves hostnames with a custom DNS
// resolver contacting Cloudflare's public resolver (1.1.1.1) over UDP and
// keeping only A (IPv4) records. This is required on hosts whose egress
// network cannot reach external IPv6 addresses (e.g. Render Oregon) while
// database/Redis hostnames still publish AAAA records first.
func Ipv4Dialer() *net.Dialer {
	res := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{Timeout: 5 * time.Second}
			return d.DialContext(ctx, "udp", "1.1.1.1:53")
		},
	}
	return &net.Dialer{
		Timeout:       10 * time.Second,
		KeepAlive:     30 * time.Second,
		Resolver:      res,
		FallbackDelay: 300 * time.Millisecond,
	}
}

// Ipv4Lookup resolves host to its IPv4 addresses only. An error is returned
// when no A record exists for the host.
func Ipv4Lookup(ctx context.Context, host string) ([]string, error) {
	addrs, err := Ipv4Dialer().Resolver.LookupHost(ctx, host)
	if err != nil {
		return nil, err
	}
	v4 := make([]string, 0, len(addrs))
	for _, a := range addrs {
		ip := net.ParseIP(a)
		if ip != nil && ip.To4() != nil {
			v4 = append(v4, a)
		}
	}
	if len(v4) == 0 {
		return nil, errors.New("no IPv4 address found for " + host)
	}
	return v4, nil
}

// IsIPv4Literal reports whether s is a plain IPv4 address literal.
func IsIPv4Literal(s string) bool {
	if s == "" {
		return false
	}
	ip := net.ParseIP(s)
	return ip != nil && ip.To4() != nil
}

// HostPortIPv4 replaces the hostname in "host:port" with its first IPv4
// address when the host is not already an IPv4 literal. Hosts that look like
// IPv6 literals (contain ':') are returned unchanged.
func HostPortIPv4(ctx context.Context, hostPort string) (string, error) {
	host, port, err := net.SplitHostPort(hostPort)
	if err != nil {
		return hostPort, err
	}
	if IsIPv4Literal(host) {
		return hostPort, nil
	}
	v4s, err := Ipv4Lookup(ctx, host)
	if err != nil {
		return hostPort, err
	}
	return net.JoinHostPort(v4s[0], port), nil
}
