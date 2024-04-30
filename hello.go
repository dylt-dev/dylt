package dylt

import "fmt"
import "rsc.io/quote"

func SayHi () {
	fmt.Println("Hiiiiiii Everybody")
}

func main() {
	fmt.Println("Hi")
	fmt.Println(quote.Go())
}

