package main

import (
	"fmt"
	"io"
	"neofy/internal"
	"os"
)

func main() {
	if err := run(os.Stdout, os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(w io.Writer, args []string) error {
	fmt.Println("Use:", w, args)
	return internal.RunApp()
}
