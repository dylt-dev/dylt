package lib

import (
	"fmt"
	"testing"
	"time"
)


func doit (n *int) {
	fmt.Println(*n)
}
func Test (t *testing.T) {
	for i := 0; i < 10; i++ {
		go func() {
			doit(&i)
		}()
	}
	time.Sleep(9)
}
