package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewVersionCmd(version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version of the CLI",
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Println(version)
		},
	}

	return cmd
}
