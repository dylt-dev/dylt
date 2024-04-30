package main

import "fmt"
import "rsc.io/quote"

func SayHi () {
	fmt.Println("Hi")
}

func main() {
	fmt.Println("Hi")
	fmt.Println(quote.Go())
}

