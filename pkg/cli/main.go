package cli

import "github.com/spf13/cobra"

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use: "goisoboot",
	}

	cmd.AddCommand(NewServerCmd())

	return cmd
}
