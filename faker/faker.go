package faker

import (
	gorandold "math/rand"
	"time"
	
	thefaker "github.com/jaswdr/faker"
	"github.com/dylt-dev/dylt/rand"

)	


func Bool () bool {
	var r gorandold.Source
	if rand.GlobalLegacySource != nil {
		r = rand.GlobalLegacySource
	} else {
		r = NewSource()
	}

	return thefaker.NewWithSeed(r).Bool()
}


func Int1000 () int {
	var r gorandold.Source
	if rand.GlobalLegacySource != nil {
		r = rand.GlobalLegacySource
	} else {
		r = NewSource()
	}

	return int(thefaker.NewWithSeed(r).Int16Between(0, 999))
}


func LoremWord () string {
	var r gorandold.Source
	if rand.GlobalLegacySource != nil {
		r = rand.GlobalLegacySource
	} else {
		r = NewSource()
	}

	return thefaker.NewWithSeed(r).Lorem().Word()
	
}

func NewSource () gorandold.Source {
	return gorandold.NewSource(time.Now().UTC().UnixNano())
}