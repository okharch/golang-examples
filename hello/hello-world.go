package main

import "fmt"

var p = fmt.Println

func main() {
	p("Hello!")
	i := 1
	p(i)
	p("Another world hello!")
}
