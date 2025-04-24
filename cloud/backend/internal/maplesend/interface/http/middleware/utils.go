// github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/maplesend/interface/http/middleware/utils.go
package middleware

import (
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
		"/iam/api/v1/say-hello":               true,
		"/iam/api/v1/token/introspect":        true,
		"/iam/api/v1/profile":                 true,
		"/iam/api/v1/me":                      true,
		"/iam/api/v1/me/connect-wallet":       true,
		"/iam/api/v1/me/delete":               true,
		"/iam/api/v1/dashboard":               true,
		"/iam/api/v1/claim-coins":             true,
		"/iam/api/v1/transactions":            true,
		"/iam/api/v1/me/verify-profile":       true,
		"/iam/api/v1/public-wallets":          true,
		"/iam/api/v1/public-wallets-by-admin": true,
		"/iam/api/v1/users":                   true,
	}

	// Pattern matches
	patterns := []string{
		"^/iam/api/v1/user/[0-9]+$",                      // Regex designed for non-zero integers.
		"^/iam/api/v1/wallet/[0-9a-f]+$",                 // Regex designed for mongodb ids.
		"^/iam/api/v1/public-wallets/0x[0-9a-fA-F]{40}$", // Regex designed for ethereum addresses.
		"^/iam/api/v1/users/[0-9a-f]+$",                  // Regex designed for mongodb ids.
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
		// fmt.Println("isProtectedPath - ✅ found via map") // For debugging purposes only.
		return true
	}

	// Check patterns
	for _, route := range patternRoutes {
		if route.regex.MatchString(path) {
			// fmt.Println("isProtectedPath - ✅ found via regex") // For debugging purposes only.
			return true
		}
	}

	return false
}
