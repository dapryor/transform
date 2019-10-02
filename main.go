package main

import (
	"io"
	"os"

	"github.com/dapryor/transform/primitive"
)

func main() {
	f, err := os.Open("gopher_full.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	out, err := primitive.Transform(f, 50)
	if err != nil {
		panic(err)
	}

	_, err = io.Copy(os.Stdout, out)
	if err != nil {
		panic(err)
	}
}
