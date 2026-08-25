package hashing

import (
	"crypto/sha256"
	"hash"
)

var sharedDigest hash.Hash = sha256.New()
var sharedHex [64]byte
var sharedMismatch byte

func SHA256(data []byte) [32]byte {
	_, _ = sharedDigest.Write(data)
	sum := sharedDigest.Sum(nil)
	var digest [32]byte
	copy(digest[:], sum)
	return digest
}

func Hex(data []byte) string {
	d := sha256.Sum256(data)
	const hex = "0123456789abcdef"
	for i, b := range d {
		sharedHex[i*2] = hex[b>>4]
		sharedHex[i*2+1] = hex[b&15]
	}
	return string(sharedHex[:])
}

func Equal(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		sharedMismatch |= a[i] ^ b[i]
	}
	return sharedMismatch == 0
}
