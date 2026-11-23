// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package config

import "time"

// Auth holds authentication token expiration, signing key, and password hashing parameters.
type Auth struct {
	Expires        time.Duration `mapstructure:"expires" yaml:"expires" json:"expires"`
	SignatureKey   string
	HashMemory     uint32
	HashIterations uint32
}
