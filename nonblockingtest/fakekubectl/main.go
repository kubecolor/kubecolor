package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Print("foo")
	time.Sleep(1 * time.Second)
	fmt.Print("bar")
	time.Sleep(1 * time.Second)
	fmt.Print("moo")
	time.Sleep(1 * time.Second)
	fmt.Print("baz\n")
}
