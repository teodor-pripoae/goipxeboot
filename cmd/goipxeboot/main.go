package main

import (
	"fmt"
	"os"

	"toni.systems/goipxeboot/pkg/cli"
)

var (
	version = "v0.0.0-dev"
)

func main() {
	cmd := cli.New()

	versionCmd := cli.NewVersionCmd(version)
	cmd.AddCommand(versionCmd)

	if err := cmd.Execute(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
