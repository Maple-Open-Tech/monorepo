// github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/maplesend/interface/http/middleware/jwtpost.go
package middleware

import (
	"net"
	"net/http"

	"go.uber.org/zap"
)

func (mid *middleware) EnforceRestrictCountryIPsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Extract IP address from request
		ipStr, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			mid.logger.Warn("failed splitting host port",
				zap.Any("url", r.URL.Path),
				zap.Any("r.RemoteAddr", r.RemoteAddr),
				zap.Any("middleware", "EnforceRestrictCountryIPsMiddleware"))
			http.Error(w, "Invalid IP address", http.StatusBadRequest)
			return
		}

		ip := net.ParseIP(ipStr)
		if ip == nil {
			mid.logger.Warn("failed parsing ip address",
				zap.Any("url", r.URL.Path),
				zap.Any("r.RemoteAddr", ipStr),
				zap.Any("middleware", "EnforceRestrictCountryIPsMiddleware"))
			http.Error(w, "Invalid IP address", http.StatusBadRequest)
			return
		}

		// Perform enforcement of country-wide blocking.
		if mid.IPCountryBlocker.IsBlockedIP(ctx, ip) {
			mid.logger.Warn("rejected request by country ip address",
				zap.Any("url", r.URL.Path),
				zap.Any("ip_address", ip),
				zap.Any("middleware", "EnforceRestrictCountryIPsMiddleware"))
			http.Error(w, "Access denied from your country", http.StatusForbidden)
			return
		}

		next(w, r.WithContext(ctx))
	}
}
