package eco

import (
	cryptorand "crypto/rand"
	"reflect"
	"testing"
)

func BenchmarkFieldMap (b *testing.B) {
	var p *EcoTest
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
