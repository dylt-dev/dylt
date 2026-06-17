package common

import (
	"encoding/json"
	"math"
	"os"
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


// Levenshstein Distance
func Lev(a string, b string) int {
	ctx := NewEcoContext(os.Stdout)
	// m := NewM(len(a), len(b))
	// return func(aa string, bb string) int { return lev(ctx, m, aa, bb) }(a, b)
	return lev3(ctx, a, b)
}


type M [][]int

func NewM(x int, y int) M {
	var m M = make([][]int, x)
	for i := range m {
		m[i] = make([]int, y)
	}

	for i := range m {
		for j := range m[i] {
			m[i][j] = -1
		}
	}
	return m
}

func (m M) Has(i int, j int) bool {
	return m[i][j] != -1
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

func lev(ctx *EcoContext, m M, a string, b string) int {
	// ctx.Signature("lev", a, b)
	ctx.Inc()
	defer ctx.Dec()
	
	var n int
	// ctx.Comment("hi")
	if len(a) == 0 {
		n = len(b)
	} else if len(b) == 0 {
		n = len(a)
	} else {
		i := len(a)-1
		j := len(b)-1
		if m.Has(i, j) {
			ctx.Infof("Hit! %q, %q => %d", a, b, n)
			return m[i][j]
		}

		if a[0] == b[0] {
			n = lev(ctx, m, a[1:], b[1:])
		} else {
			n = 1 + min(lev(ctx, m, a[1:], b), lev(ctx, m, a, b[1:]), lev(ctx, m, a[1:], b[1:]))
		}
		m[i][j] = n

	}
	
	ctx.Infof("%q, %q => %d", a, b, n)
	return n
}


func lev3 (ctx *EcoContext, a string, b string) int {
	var lastCol, currCol, swap *[]int
	slice1 := make([]int, len(a)+1)
	slice2 := make([]int, len(a)+1)
	lastCol = &slice1
	currCol = &slice2

	var u, l, ul int
	for j := range b {
		for i := range a {
			ctx.Infof("i=%d j=%d", i, j)
			if i == j {
				if i == 0 && j == 0 {
					(*currCol)[0] = 0
				} else if a[i] == b[j] {
					(*currCol)[i] = (*lastCol)[i-1]
				} else {
					(*currCol)[i] = (*lastCol)[i-1] + 1
				}
			} else {
				if i == 0 {
					(*currCol)[i] = j
				} else if j == 0 {
					(*currCol)[i] = i
				} else {
					// Initalize vars to max values
					l = math.MaxInt
					u = math.MaxInt
					ul = math.MaxInt

					// l
					if j > 0 {
						l = (*lastCol)[i] + 1
					}

					// u
					if i > 0 {
						u = (*currCol)[i-1] + 1
					}

					// ul
					if i > 0 && j > 0 {
						if a[i] == b[j] {
							ul = (*lastCol)[i-1]
						} else {
							ul = (*lastCol)[i-1]+1
						}
					}
		
					ctx.Infof("l=%d u=%d ul=%d", l, u, ul)
					(*currCol)[i] = min(l, u, ul)
				}
			}

			ctx.Infof("currCol[%d]=%d", i, (*currCol)[i])
		}

		swap = lastCol
		lastCol = currCol
		currCol = swap
	}

	return (*lastCol)[len(a)-1]
}

