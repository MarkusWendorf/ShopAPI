package argon

import (
	"crypto/rand"
	"encoding/hex"
	"golang.org/x/crypto/argon2"
	"io"
)

var saltLength = 16
var time, memory, keyLen uint32 = 3, 64 * 1024, 64
var threads uint8 = 2

// Hash calculates a ArgonID argon and prepends it with the salt being used.
// Returns a hex string of salt + argon.
func Hash(password string) (string, error) {

	salt := make([]byte, saltLength)
	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		return "", err
	}

	pw := argon2.IDKey([]byte(password), salt, time, memory, threads, keyLen)

	return hex.EncodeToString(salt) + hex.EncodeToString(pw), nil
}

// CheckPassword compares the provided password in plaintext with its hashed value.
// Returns hexHash == argon(passwordPlain).
func CheckPassword(passwordPlain string, hexHash string) bool {

	saltSeparator := hex.EncodedLen(saltLength)

	hexSalt := []byte(hexHash[:saltSeparator])
	hexPw := hexHash[saltSeparator:]

	salt := make([]byte, saltLength)
	hex.Decode(salt, hexSalt)

	comparePassword := argon2.IDKey([]byte(passwordPlain), salt, time, memory, threads, keyLen)

	return hex.EncodeToString(comparePassword) == hexPw
}
