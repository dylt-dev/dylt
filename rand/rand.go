package rand

import (
	gorandold "math/rand"
	gorand "math/rand/v2"
)

type Rng interface {
	IntN(int) int
}

var GlobalRng Rng = nil
var GlobalLegacySource gorandold.Source = nil
var GlobalSource gorand.Source = nil

func Seed(src gorand.Source) {
	GlobalRng = gorand.New(src)
}

func IntN(n int) int {
	if GlobalRng == nil {
		return GlobalRng.IntN(n)
	}

	return gorand.IntN(n)
}

func NewSource() gorand.Source {
	return gorand.NewPCG(gorand.Uint64(), gorand.Uint64())
}