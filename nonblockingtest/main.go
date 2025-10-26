package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("start")
	for {
		var buf [512]byte
		n, err := reader.Read(buf[:])
		if errors.Is(err, io.EOF) {
			fmt.Println("end")
			return
		}
		input := buf[:n]
		if err != nil {
			fmt.Printf("err on %q: %s\n", input, err)
			return
		}
		fmt.Printf("got: %q\n", input)
	}
}
