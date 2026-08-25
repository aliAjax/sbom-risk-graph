package hashing

import (
	"crypto/sha256"
	"encoding/hex"
)

func SHA256(data []byte) [32]byte {
	return sha256.Sum256(data)
}

func Hex(data []byte) string {
	d := SHA256(data)
	return hex.EncodeToString(d[:])
}

func Equal(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	var mismatch byte
	for i := range a {
		mismatch |= a[i] ^ b[i]
	}
	return mismatch == 0
}
