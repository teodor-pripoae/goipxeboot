package cli

import "github.com/spf13/cobra"

func server(cmd *cobra.Command, args []string) error {
	return nil
}

func NewServerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Start the goisoboot server",
		Long:  "Starts TFTP and HTTP Server",
		RunE:  server,
	}

	cmd.Flags().IntP("http-port", "p", 8080, "HTTP Port")

	return cmd
}
