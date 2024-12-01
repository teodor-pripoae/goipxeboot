package cli

import "github.com/spf13/cobra"

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use: "goipxeboot",
	}

	cmd.AddCommand(NewServerCmd())

	return cmd
}
