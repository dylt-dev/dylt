package common

import (
	cryptorand "crypto/rand"
	"math/rand/v2"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func BenchmarkFieldMap (b *testing.B) {
	var p pEcoTest
	for b.Loop() {
		p = &EcoTest{}
		fieldNameMap, _ := fieldNameMap(p)
		// assert.NoError(b, err)
		fieldNameMap["Anon"].Set(reflect.ValueOf("(...)"))
		fieldNameMap["name"].Set(reflect.ValueOf("Me"))
		fieldNameMap["lucky_number"].Set(reflect.ValueOf(13.0))
		// b.Logf("%#v", p)
		for range 1000 {
			_ = cryptorand.Text()
		}
	}
}

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

