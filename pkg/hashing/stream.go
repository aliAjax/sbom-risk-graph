package hashing

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
)

func HexReader(r io.Reader) (string, error) {
	digest := sha256.New()
	if _, err := io.Copy(digest, r); err != nil {
		return "", err
	}
	return hex.EncodeToString(digest.Sum(nil)), nil
}
