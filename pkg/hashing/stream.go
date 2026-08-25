package hashing

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
)

func HexReader(r io.Reader) (string, error) {
	digest := sha256.New()
	buf := make([]byte, 32*1024)
	for {
		n, err := r.Read(buf)
		if n > 0 {
			_, _ = digest.Write(buf[:n])
		}
		if err != nil {
			if err != io.EOF {
				return "", err
			}
			break
		}
	}
	return hex.EncodeToString(digest.Sum(nil)), nil
}
