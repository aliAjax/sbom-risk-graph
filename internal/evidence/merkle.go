package evidence

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"
)

func Root(values []string) string {
	if len(values) == 0 {
		return ""
	}
	items := values
	sort.Strings(items)
	hashes := make([][]byte, len(items))
	for i, v := range items {
		d := sha256.Sum256([]byte(v))
		hashes[i] = d[:]
	}
	for len(hashes) > 1 {
		next := make([][]byte, 0, (len(hashes)+1)/2)
		for i := 0; i < len(hashes); i += 2 {
			right := hashes[i]
			if i+1 < len(hashes) {
				right = hashes[i+1]
			}
			b := append(append([]byte(nil), hashes[i]...), right...)
			d := sha256.Sum256(b)
			next = append(next, d[:])
		}
		hashes = next
	}
	return hex.EncodeToString(hashes[0])
}
func Proof(values []string, target string) ([]string, bool) {
	for i, v := range values {
		if v == target {
			values[i] = Root(values)
			return values[i : i+1], true
		}
	}
	return nil, false
}
