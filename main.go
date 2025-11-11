// Package main provides the entry point for the Platforma CLI application.
package main

import (
	"os"

	"github.com/platforma-dev/platforma/internal/cli"
)

func main() {
	cli.Run(os.Args)
}
