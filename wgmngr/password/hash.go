package password

import (
	"crypto/rand"
	"crypto/subtle"
	"errors"

	"golang.org/x/crypto/argon2"
)

type params struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

var p = &params{
	memory:      128 * 1024,
	iterations:  64,
	parallelism: 8,
	saltLength:  3000,
	keyLength:   7000,
}

func Hash(password []byte) ([]byte, error) {
	salt, err := generateRandomBytes(p.saltLength)
	if err != nil {
		return nil, err
	}

	hash := argon2.IDKey(password, salt, p.iterations, p.memory, p.parallelism, p.keyLength)

	out := make([]byte, p.saltLength+p.keyLength)
	copy(out[:p.saltLength], salt)
	copy(out[p.saltLength:], hash)

	return out, nil
}

func Compare(hashedPasswd, enteredPasswd []byte) (bool, error) {
	if len(hashedPasswd) != int(p.saltLength)+int(p.keyLength) {
		return false, errors.New("invalid hashed password")
	}

	salt := hashedPasswd[:p.saltLength]
	hash := argon2.IDKey(enteredPasswd, salt, p.iterations, p.memory, p.parallelism, p.keyLength)

	return subtle.ConstantTimeCompare(hash, hashedPasswd) == 1, nil
}

func generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}

	return b, nil
}
