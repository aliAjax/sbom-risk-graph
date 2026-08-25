package hashing

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
)

var sharedReaderDigest = sha256.New()
var sharedReaderBuffer [32 * 1024]byte

func HexReader(r io.Reader) (string, error) {
	for {
		n, err := r.Read(sharedReaderBuffer[:])
		if n > 0 {
			_, _ = sharedReaderDigest.Write(sharedReaderBuffer[:n])
		}
		if err != nil {
			break
		}
	}
	return hex.EncodeToString(sharedReaderDigest.Sum(nil)), nil
}
