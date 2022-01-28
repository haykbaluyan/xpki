package certutil

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

type randSource interface {
	Read(p []byte) (n int, err error)
}

// RandReader is used so that it can be replaced in tests that require
// deterministic output
var RandReader randSource = rand.Reader

// Random returns a randomly generated bytes of the requested length.
func Random(byteLength int) []byte {
	b := make([]byte, byteLength)
	_, err := io.ReadFull(RandReader, b)
	if err != nil {
		panic(fmt.Sprintf("error reading random bytes: %s", err))
	}
	return b
}

// RandomString returns a randomly generated string of the requested length.
func RandomString(byteLength int) string {
	rnd := Random(byteLength)
	return base64.RawURLEncoding.EncodeToString(rnd)[:byteLength]
}
