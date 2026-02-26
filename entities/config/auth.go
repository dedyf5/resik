// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package config

type Auth struct {
	Expires        uint64
	SignatureKey   string
	HashMemory     uint32
	HashIterations uint32
}
