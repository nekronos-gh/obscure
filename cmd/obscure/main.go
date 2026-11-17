package main

import (
	"fmt"
	"os"

	"obscure/internal/backend/cuda"
	"obscure/internal/pipeline"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: obscure <video>")
		os.Exit(2)
	}
	if err := pipeline.Run(os.Args[1], cuda.Backend{}); err != nil {
		fmt.Fprintln(os.Stderr, "obscure:", err)
		os.Exit(1)
	}
}
