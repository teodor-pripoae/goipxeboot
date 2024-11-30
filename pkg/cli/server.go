package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"toni.systems/goisoboot/pkg/ipxe"
	"toni.systems/goisoboot/pkg/tftp"
)

func quit(message string) {
	fmt.Printf("Error: %s\n", message)
	os.Exit(1)
}

func server(cmd *cobra.Command, args []string) {
	serverIP, err := cmd.Flags().GetString("server-ip")
	if err != nil {
		quit(err.Error())
	} else if serverIP == "" {
		quit("server-ip is required")
	}

	httpPort, err := cmd.Flags().GetInt("http-port")
	if err != nil {
		quit(err.Error())
	}

	tftp := tftp.New()
	server, err := ipxe.New(
		ipxe.WithIP(serverIP),
		ipxe.WithPort(httpPort),
	)

	errChan := make(chan error)
	done := make(chan struct{})

	go func() {
		if err := tftp.Run(); err != nil {
			errChan <- err
		}
		done <- struct{}{}
	}()

	go func() {
		if err := server.Run(); err != nil {
			errChan <- err
		}
		done <- struct{}{}
	}()

	select {
	case err := <-errChan:
		quit(err.Error())
	case <-done:
		fmt.Println("Server stopped")
	}
}

func NewServerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Start the goisoboot server",
		Long:  "Starts TFTP and HTTP Server",
		Run:   server,
	}

	cmd.Flags().StringP("server-ip", "i", "", "Server IP")
	cmd.Flags().IntP("http-port", "p", 8080, "HTTP Port")

	return cmd
}
