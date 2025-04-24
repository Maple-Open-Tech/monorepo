// github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/maplesend/interface/http/middleware/blacklist.go
package middleware

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/Maple-Open-Tech/monorepo/cloud/backend/config/constants"
)

// Note: This middleware must have `IPAddressMiddleware` executed first before running.
func (mid *middleware) EnforceBlacklistMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Open our program's context based on the request and save the
		// slash-seperated array from our URL path.
		ctx := r.Context()

		ipAddress, _ := ctx.Value(constants.SessionIPAddress).(string)
		proxies, _ := ctx.Value(constants.SessionProxies).(string)

		// Case 1 of 3: Check banned IP addresses.
		if mid.blacklist.IsBannedIPAddress(ipAddress) {

			// If the client IP address is banned, check to see if the client
			// is making a call to a URL which is not banned, and if the URL
			// is not banned (has not been banned before) then print it to
			// the console logs for future analysis. Else if the URL is banned
			// then don't bother printing to console. The purpose of this code
			// is to not clog the console log with warnings.
			if !mid.blacklist.IsBannedURL(r.URL.Path) {
				mid.logger.Warn("rejected request by ip",
					zap.Any("url", r.URL.Path),
					zap.String("ip_address", ipAddress),
					zap.String("proxies", proxies),
					zap.Any("middleware", "EnforceBlacklistMiddleware"))
			}
			http.Error(w, "forbidden at this time", http.StatusForbidden)
			return
		}

		// Case 2 of 3: Check banned URL.
		if mid.blacklist.IsBannedURL(r.URL.Path) {

			// If the URL is banned, check to see if the client IP address is
			// banned, and if the client has not been banned before then print
			// to console the new offending client ip. Else do not print if
			// the offending client IP address has been banned before. The
			// purpose of this code is to not clog the console log with warnings.
			if !mid.blacklist.IsBannedIPAddress(ipAddress) {
				mid.logger.Warn("rejected request by url",
					zap.Any("url", r.URL.Path),
					zap.String("ip_address", ipAddress),
					zap.String("proxies", proxies),
					zap.Any("middleware", "EnforceBlacklistMiddleware"))
			}

			// DEVELOPERS NOTE:
			// Simply return a 404, but in our console log we can see the IP
			// address whom made this call.
			http.Error(w, "does not exist at this time", http.StatusNotFound)
			return
		}

		// Case 3 of 3: Check banned IP via flagged content.
		bannedIPAddresses, err := mid.bannedIPAddressListAllValuesUseCase.Execute(ctx)
		if err != nil {
			mid.logger.Error("Failed listing all banned IP addresses",
				zap.Any("error", err),
				zap.Any("middleware", "EnforceBlacklistMiddleware"))
			http.Error(w, "forbidden at this time", http.StatusForbidden)
			return
		}
		for _, bannedIP := range bannedIPAddresses {
			if bannedIP == ipAddress {
				mid.logger.Warn("rejected request by banned IP address for flagged content",
					zap.Any("url", r.URL.Path),
					zap.String("ip_address", ipAddress),
					zap.String("proxies", proxies),
					zap.Any("middleware", "EnforceBlacklistMiddleware"))
				http.Error(w, "forbidden at this time", http.StatusForbidden)
				return
			}

		}

		next(w, r.WithContext(ctx))
	}
}
