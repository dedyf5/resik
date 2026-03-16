// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

package hash

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"math"
	"strings"

	"golang.org/x/crypto/argon2"
)

var (
	ErrInvalidHash         = errors.New("invalid hash format")
	ErrHashTooLong         = errors.New("hash length exceeds uint32 limit")
	ErrIncompatibleVersion = errors.New("incompatible argon2 version")
)

type Argon2Config struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

type argon2Hasher struct {
	config *Argon2Config
	source io.Reader
}

func NewArgon2Hasher(config *Argon2Config) IHash {
	if config.Memory < 8*1024 {
		config.Memory = 8 * 1024
	}
	if config.Iterations < 1 {
		config.Iterations = 3
	}
	if config.Parallelism < 1 {
		config.Parallelism = 1
	}
	if config.SaltLength < 1 {
		config.SaltLength = 16
	}
	if config.KeyLength < 1 {
		config.KeyLength = 32
	}

	return &argon2Hasher{
		config: config,
		source: rand.Reader,
	}
}

func (h *argon2Hasher) Hash(password string) (string, error) {
	salt := make([]byte, h.config.SaltLength)
	if _, err := h.source.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, h.config.Iterations, h.config.Memory, h.config.Parallelism, h.config.KeyLength)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, h.config.Memory, h.config.Iterations, h.config.Parallelism, b64Salt, b64Hash)

	return encodedHash, nil
}

func (h *argon2Hasher) Compare(password, encodedHash string) (bool, error) {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return false, ErrInvalidHash
	}

	if parts[1] != "argon2id" {
		return false, ErrInvalidHash
	}

	var version int
	_, err := fmt.Sscanf(parts[2], "v=%d", &version)
	if err != nil {
		return false, err
	}
	if version != argon2.Version {
		return false, ErrIncompatibleVersion
	}

	var memory, iterations uint32
	var parallelism uint8
	_, err = fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &iterations, &parallelism)
	if err != nil {
		return false, err
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}

	decodedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}

	hashLen := len(decodedHash)
	if hashLen > math.MaxUint32 {
		return false, ErrHashTooLong
	}

	keyLength := uint32(hashLen)

	comparisonHash := argon2.IDKey([]byte(password), salt, iterations, memory, parallelism, keyLength)

	if subtle.ConstantTimeCompare(decodedHash, comparisonHash) == 1 {
		return true, nil
	}

	return false, nil
}
