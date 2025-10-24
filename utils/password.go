package utils

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

// PasswordHasher provides password hashing and verification functionality
type PasswordHasher struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

// NewPasswordHasher creates a new password hasher with default parameters
func NewPasswordHasher() *PasswordHasher {
	return &PasswordHasher{
		memory:      64 * 1024, // 64 MB
		iterations:  3,
		parallelism: 2,
		saltLength:  16,
		keyLength:   32,
	}
}

// Hash hashes a password using Argon2
func (ph *PasswordHasher) Hash(password string) (string, error) {
	// Generate a random salt
	salt := make([]byte, ph.saltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	// Hash the password
	hash := argon2.IDKey([]byte(password), salt, ph.iterations, ph.memory, ph.parallelism, ph.keyLength)

	// Encode the hash in the format: $argon2id$v=19$m=65536,t=3,p=2$salt$hash
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, ph.memory, ph.iterations, ph.parallelism, b64Salt, b64Hash)

	return encodedHash, nil
}

// Verify verifies a password against a hash
func (ph *PasswordHasher) Verify(password, encodedHash string) (bool, error) {
	// Extract the parameters, salt and derived key from the encoded password hash
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return false, errors.New("invalid hash format")
	}

	// Parse parameters
	var version int
	_, err := fmt.Sscanf(parts[2], "v=%d", &version)
	if err != nil {
		return false, fmt.Errorf("failed to parse version: %w", err)
	}
	if version != argon2.Version {
		return false, errors.New("incompatible version")
	}

	var memory, iterations uint32
	var parallelism uint8
	_, err = fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &iterations, &parallelism)
	if err != nil {
		return false, fmt.Errorf("failed to parse parameters: %w", err)
	}

	// Decode salt
	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, fmt.Errorf("failed to decode salt: %w", err)
	}

	// Decode hash
	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, fmt.Errorf("failed to decode hash: %w", err)
	}

	// Hash the password using the same parameters
	keyLength := uint32(len(hash))
	comparisonHash := argon2.IDKey([]byte(password), salt, iterations, memory, parallelism, keyLength)

	// Compare the hashes in constant time
	if subtle.ConstantTimeCompare(hash, comparisonHash) == 1 {
		return true, nil
	}

	return false, nil
}
