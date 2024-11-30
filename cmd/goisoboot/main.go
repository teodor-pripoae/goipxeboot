package main

import (
	"fmt"
	"os"

	"toni.systems/goisoboot/pkg/cli"
)

func main() {
	cmd := cli.New()
	if err := cmd.Execute(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
