package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"toni.systems/goisoboot/pkg/config"
	"toni.systems/goisoboot/pkg/ipxe"
	"toni.systems/goisoboot/pkg/tftp"
)

func quit(message string) {
	fmt.Printf("Error: %s\n", message)
	os.Exit(1)
}

func server(cmd *cobra.Command, args []string) {
	configFile, err := cmd.Flags().GetString("config")
	if err != nil {
		quit(err.Error())
	}
	config, err := config.FromFile(configFile)
	if err != nil {
		quit(err.Error())
	}

	tftp := tftp.New()
	server, err := ipxe.New(
		ipxe.WithIP(config.HTTP.IP),
		ipxe.WithPort(config.HTTP.Port),
		ipxe.WithRootDir(config.GetRootDir()),
		ipxe.WithIPXE(config.IPXE),
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

	cmd.Flags().StringP("config", "c", "goisoboot.yaml", "Config file to use")

	return cmd
}
