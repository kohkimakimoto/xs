package main

import (
	"fmt"
	"os"

	"github.com/kohkimakimoto/xs/internal"
)

func main() {
	if err := internal.Run(os.Args); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
