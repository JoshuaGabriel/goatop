package main

import (
	"os"
	"github.com/JoshuaGabriel/goup/cmd"
)

func main() {
	if err := cmd.Main(os.Args); err != nil {
		os.Exit(1)
	}
}
