package common

import (
	cryptorand "crypto/rand"
	"math/rand/v2"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewChaCha8 (t *testing.T) {
	seedarray := [32]byte{}
	var seed []byte = seedarray[:]
	cryptorand.Read(seed)
	chacha8 := rand.NewChaCha8(seedarray)
	assert.NotNil(t, chacha8)
}

func BenchmarkCryptoReadText (b *testing.B) {
	buf := make([]byte, 0, 100000)
	for b.Loop() {
		for range 1000 {
			_ = cryptorand.Text()
			cryptorand.Read(buf)
		}
	}
}

func BenchmarkChacha8ReadText (b *testing.B) {
	seedarray := [32]byte{}
	var seed []byte = seedarray[:]
	cryptorand.Read(seed)
	chacha8 := rand.NewChaCha8(seedarray)
	for b.Loop() {
		buf := make([]byte, 0, 100000)
		for range 10000000 {
			chacha8.Read(buf)
			_ = buf[0] + buf[len(buf)-1]
		}
	}
}

func BenchmarkNada (b *testing.B) {
	for b.Loop() {
		for range 10000000 {
		}
	}
}

