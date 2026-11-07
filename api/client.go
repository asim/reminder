package api

import (
	"fmt"
	"net"
	"net/http"
	"net/netip"
	"strings"

	"github.com/asim/reminder/model"
)

// ClientMiddleware only accepts X-Real-IP from explicitly trusted proxy peers.
// An empty configuration trusts loopback, for nginx on the same host.
func ClientMiddleware(cidrs string) (func(http.Handler) http.Handler, error) {
	if cidrs == "" {
		cidrs = "127.0.0.0/8,::1/128"
	}
	var trusted []netip.Prefix
	for _, raw := range strings.Split(cidrs, ",") {
		p, err := netip.ParsePrefix(strings.TrimSpace(raw))
		if err != nil {
			return nil, fmt.Errorf("invalid TRUSTED_PROXY_CIDRS: %w", err)
		}
		trusted = append(trusted, p)
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			host, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				host = r.RemoteAddr
			}
			ip, err := netip.ParseAddr(host)
			if err != nil {
				http.Error(w, "Invalid client address", 400)
				return
			}
			ip = ip.Unmap()
			for _, p := range trusted {
				if p.Contains(ip) {
					if forwarded, err := netip.ParseAddr(r.Header.Get("X-Real-IP")); err == nil {
						ip = forwarded.Unmap()
					}
					break
				}
			}
			next.ServeHTTP(w, r.WithContext(model.WithClient(r.Context(), ip.String())))
		})
	}, nil
}
