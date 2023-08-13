package middleware

// Ported from Goji's middleware, source:
// https://github.com/zenazn/goji/tree/master/web/middleware

import (
	"net"
	"net/http"
	"strings"
)

const (
	trueClientIP  = "True-Client-IP"
	xRealIP       = "X-Real-IP"
	xForwardedFor = "X-Forwarded-For"
)

// RealIP is a middleware that sets a http.Request's RemoteAddr to the results
// of parsing either the True-Client-IP, X-Real-IP or the X-Forwarded-For headers
// (in that order).
//
// This middleware should be inserted fairly early in the middleware stack to
// ensure that subsequent layers (e.g., request loggers) which examine the
// RemoteAddr will see the intended value.
//
// You should only use this middleware if you can trust the headers passed to
// you (in particular, the two headers this middleware uses), for example
// because you have placed a reverse proxy like HAProxy or nginx in front of
// chi. If your reverse proxies are configured to pass along arbitrary header
// values from the client, or if you use this middleware without a reverse
// proxy, malicious clients will be able to make you very sad (or, depending on
// how you're using RemoteAddr, vulnerable to an attack of some sort).
func RealIP(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if rip := realIP(r); rip != "" {
			r.RemoteAddr = rip
		}
		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func realIP(r *http.Request) string {
	if tcip := r.Header.Get(trueClientIP); isValidIP(tcip) {
		return tcip
	}

	if xrip := r.Header.Get(xRealIP); isValidIP(xrip) {
		return xrip
	}

	if xff := r.Header.Get(xForwardedFor); xff != "" {
		i := strings.Index(xff, ",")
		if i == -1 {
			i = len(xff)
		}
		return xff[:i]
	}

	return ""
}

func isValidIP(ip string) bool {
	return ip != "" && net.ParseIP(ip) != nil
}
