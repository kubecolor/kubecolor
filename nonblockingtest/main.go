package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
)

func main() {
	cmd := exec.Command("go", "run", "./nonblockingtest/fakekubectl")
	r, w, err := os.Pipe()
	if err != nil {
		panic(err)
	}

	cmd.Stdout = w

	go func() {
		if err := cmd.Run(); err != nil {
			panic(err)
		}
		w.Close()
	}()

	reader := bufio.NewReader(r)

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
