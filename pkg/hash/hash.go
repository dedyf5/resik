// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package hash

//go:generate mockgen -source hash.go -package mock -destination ./mock/hash.go
type IHash interface {
	Hash(password string) (string, error)
	Compare(password, encodedHash string) (bool, error)
}
