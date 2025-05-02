// github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/mapleauth/interface/http/middleware/middleware.go
package middleware

import (
	"context"
	"net/http"

	uc_bannedipaddress "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/mapleauth/usecase/bannedipaddress"
	uc_user "github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/mapleauth/usecase/federateduser"
	"github.com/Maple-Open-Tech/monorepo/cloud/backend/pkg/security/jwt"
)

type Middleware interface {
	Attach(fn http.HandlerFunc) http.HandlerFunc
	Shutdown(ctx context.Context)
}

type middleware struct {
	jwt                                 jwt.Provider
	userGetBySessionIDUseCase           uc_user.UserGetBySessionIDUseCase
	bannedIPAddressListAllValuesUseCase uc_bannedipaddress.BannedIPAddressListAllValuesUseCase
}

func NewMiddleware(
	jwtp jwt.Provider,
	uc1 uc_user.UserGetBySessionIDUseCase,
	uc2 uc_bannedipaddress.BannedIPAddressListAllValuesUseCase,
) Middleware {
	return &middleware{
		jwt:                                 jwtp,
		userGetBySessionIDUseCase:           uc1,
		bannedIPAddressListAllValuesUseCase: uc2,
	}
}

// Attach function attaches to HTTP router to apply for every API call.
func (mid *middleware) Attach(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Apply base middleware to all requests
		handler := mid.applyBaseMiddleware(fn)

		// Check if the path requires authentication
		if isProtectedPath(r.URL.Path) {
			// Apply auth middleware for protected paths
			handler = mid.PostJWTProcessorMiddleware(handler)
			handler = mid.JWTProcessorMiddleware(handler)
			// handler = mid.EnforceBlacklistMiddleware(handler)
		}

		handler(w, r)
	}
}

// Attach function attaches to HTTP router to apply for every API call.
func (mid *middleware) applyBaseMiddleware(fn http.HandlerFunc) http.HandlerFunc {
	// Apply middleware in reverse order (bottom up)
	handler := fn
	handler = mid.URLProcessorMiddleware(handler)

	return handler
}

// Shutdown shuts down the middleware.
func (mid *middleware) Shutdown(ctx context.Context) {
	// Log a message to indicate that the HTTP server is shutting down.
}
