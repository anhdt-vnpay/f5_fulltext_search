package cmd

import (
	"github.com/anhdt-vnpay/f5_fulltext_search/server"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var RootCmd = &cobra.Command{
	Use:   "main",
	Short: "",
	Long:  "",
}

var ApiCmd = &cobra.Command{
	Use:   "api",
	Short: "Start api service",
	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetInt("port")
		server.GatewayServer(port)
	},
}

func init() {
	ApiCmd.Flags().Int("port", 8000, "port")

	RootCmd.AddCommand(ApiCmd)
}
