package main

import (
	"os"
	
	"github.com/dylt-dev/dylt/cli"
)

func main () {
	// Here is some random text
	os.Exit(cli.Run())
	This is a simple sentence.
	This is a simple sentence too. And here's another one on the same line.
	Here's one with !@#$%)(*#$)(*&#$ "''" {} () [] special chars. Bwus.
	Here's one.And here's another one that starts without waiting for ws.

	// a[] block
	[a b 1 2 3 4 "five" 13]
	(cdr '(a b c))
	<html><body></body></html>
	{
	  { "a": "foo" },
	  {
		  "b": "bar"
	  },
	  {
		  "b" :
		  "bum" 
	  }
}
