package security

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

type PasswordService struct{}

func NewPasswordService() *PasswordService {
	return &PasswordService{}
}

const (
	argonTime    = 1
	argonMemory  = 64 * 1024
	argonThreads = 4
	argonKeyLen  = 32
)

func (s *PasswordService) Hash(password string) (string, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		argonTime,
		argonMemory,
		argonThreads,
		argonKeyLen,
	)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encoded := fmt.Sprintf(
		"$argon2id$t=%d$m=%d$p=%d$%s$%s",
		argonTime, argonMemory, argonThreads, b64Salt, b64Hash,
	)

	return encoded, nil

}

func (s *PasswordService) Verify(password string, encodedHash string) bool {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 7 {
		return false
	}

	var t, m, p uint32

	_, err := fmt.Sscanf(parts[2], "t=%d", &t)
	if err != nil {
		return false
	}

	_, err = fmt.Sscanf(parts[3], "m=%d", &m)
	if err != nil {
		return false
	}

	_, err = fmt.Sscanf(parts[4], "p=%d", &p)
	if err != nil {
		return false
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false
	}

	hash, err := base64.RawStdEncoding.DecodeString(parts[6])
	if err != nil {
		return false
	}

	computed := argon2.IDKey(
		[]byte(password),
		salt,
		t,
		m,
		uint8(p),
		uint32(len(hash)),
	)

	return subtle.ConstantTimeCompare(hash, computed) == 1
}
