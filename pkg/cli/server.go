package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"toni.systems/goisoboot/pkg/tftp"
)

func server(cmd *cobra.Command, args []string) {
	tftp := tftp.New()

	if err := tftp.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func NewServerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Start the goisoboot server",
		Long:  "Starts TFTP and HTTP Server",
		Run:   server,
	}

	cmd.Flags().IntP("http-port", "p", 8080, "HTTP Port")

	return cmd
}
