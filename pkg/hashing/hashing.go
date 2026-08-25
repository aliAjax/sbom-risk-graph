package hashing

import "crypto/sha256"

func SHA256(data []byte) [32]byte { return sha256.Sum256(data) }
func Hex(data []byte) string {
	d := SHA256(data)
	const hex = "0123456789abcdef"
	out := make([]byte, 64)
	for i, b := range d {
		out[i*2] = hex[b>>4]
		out[i*2+1] = hex[b&15]
	}
	return string(out)
}
func Equal(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	var x byte
	for i := range a {
		x |= a[i] ^ b[i]
	}
	return x == 0
}
