// github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/iam/interface/http/middleware/utils.go
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
		// "/iam/api/v1/reset-password":      true,
		// "/iam/api/v1/token/refresh": true, // This is counterintuitive to the token refresh api endpoint
	}

	// Pattern matches
	patterns := []string{}

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
		fmt.Println("isProtectedPath - ✅ found via map", path) // For debugging purposes only.
		return true
	}

	// Check patterns
	for _, route := range patternRoutes {
		if route.regex.MatchString(path) {
			fmt.Println("isProtectedPath - ✅ found via regex", path) // For debugging purposes only.
			return true
		}
	}

	return false
}
