package cmd

import (
	"aid-server/configs"
	"aid-server/server"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"net"
	"net/http"
	"os"
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
		return serve()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}

func serve() error {
	ln, err := net.Listen("tcp", net.JoinHostPort(configs.Configs.Host.Host, configs.Configs.Host.Port))
	if err != nil {
		return err
	}

	err = server.Serve(ln)
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}

	return err
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
