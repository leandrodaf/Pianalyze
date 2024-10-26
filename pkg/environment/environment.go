package environment

import "os"

// IsProduction checks if the current environment is set to production
func IsProduction() bool {
	return os.Getenv("GO_ENV") == "production"
}
