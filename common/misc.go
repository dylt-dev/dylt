package common

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func GetNumDigits(n int) int {
	if n == 0 {
		return 1
	}
	count := 0
	for n != 0 {
		n /= 10
		count++
	}

	return count
}


func Marshal (a any) []byte {
	var buf []byte

	buf, err := json.Marshal(a)
	if err != nil {
		panic(err)
	}

	return buf
}



func MarshalAndTest(t *testing.T, a any) []byte {
	var buf []byte
	var err error

	_, is := a.(string)
	if is {
		// buf = []byte(s)
		buf, err = json.Marshal(a)
		require.NoError(t, err)
	} else {
		buf, err = json.Marshal(a)
		require.NoError(t, err)
	}

	return buf
}
