package cmd

import (
	"aid-server/configs"
	"aid-server/server"
	"errors"
	"github.com/spf13/cobra"
	"net"
	"net/http"
)

var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "An AID Server CLI application",
	Long:  "You can use this CLI to start the AID web server.",
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the AID web server",
	Long:  "Start the AID web server and handle incoming requests.",
	RunE: func(cmd *cobra.Command, args []string) error {
		ln, err := net.Listen("tcp", net.JoinHostPort(configs.Configs.Host.Host, configs.Configs.Host.Port))
		if err != nil {
			return err
		}
		return serve(ln, server.Serve)
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}

func serve(listener net.Listener, serveFunc func(ln net.Listener) error) error {
	err := serveFunc(listener)
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}

	return err
}

func Execute() error {
	return rootCmd.Execute()
}
