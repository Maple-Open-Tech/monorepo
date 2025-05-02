// github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/ipe/interface/http/middleware/utils.go
package middleware

import (
	"fmt"
	"regexp"
)

type protectedRoute struct {
	pattern string
	regex   *regexp.Regexp
}

var (
	exactPaths    = make(map[string]bool)
	patternRoutes []protectedRoute
)

func init() {
	// Exact matches
	exactPaths = map[string]bool{
		"/ipe/api/v1/say-hello":               true,
		"/ipe/api/v1/token/introspect":        true,
		"/ipe/api/v1/profile":                 true,
		"/ipe/api/v1/me":                      true,
		"/ipe/api/v1/me/connect-wallet":       true,
		"/ipe/api/v1/me/delete":               true,
		"/ipe/api/v1/dashboard":               true,
		"/ipe/api/v1/claim-coins":             true,
		"/ipe/api/v1/transactions":            true,
		"/ipe/api/v1/me/verify-profile":       true,
		"/ipe/api/v1/public-wallets":          true,
		"/ipe/api/v1/public-wallets-by-admin": true,
		"/ipe/api/v1/users":                   true,
	}

	// Pattern matches
	patterns := []string{
		"^/ipe/api/v1/user/[0-9]+$",                      // Regex designed for non-zero integers.
		"^/ipe/api/v1/wallet/[0-9a-f]+$",                 // Regex designed for mongodb ids.
		"^/ipe/api/v1/public-wallets/0x[0-9a-fA-F]{40}$", // Regex designed for ethereum addresses.
		"^/ipe/api/v1/users/[0-9a-f]+$",                  // Regex designed for mongodb ids.
	}

	// Precompile patterns
	patternRoutes = make([]protectedRoute, len(patterns))
	for i, pattern := range patterns {
		patternRoutes[i] = protectedRoute{
			pattern: pattern,
			regex:   regexp.MustCompile(pattern),
		}
	}
}

func isProtectedPath(path string) bool {
	// fmt.Println("isProtectedPath - path:", path) // For debugging purposes only.

	// Check exact matches first (O(1) lookup)
	if exactPaths[path] {
		fmt.Println("isProtectedPath - ✅ found via map") // For debugging purposes only.
		return true
	}

	// Check patterns
	for _, route := range patternRoutes {
		if route.regex.MatchString(path) {
			fmt.Println("isProtectedPath - ✅ found via regex") // For debugging purposes only.
			return true
		}
	}

	return false
}
