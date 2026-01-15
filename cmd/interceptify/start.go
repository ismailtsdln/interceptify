package interceptify

import (
	"fmt"
	"os"

	"github.com/ismailtsdln/interceptify/pkg/attack"
	"github.com/ismailtsdln/interceptify/pkg/ca"
	"github.com/ismailtsdln/interceptify/pkg/proxy"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the Interceptify proxy engine",
	Long:  `Start the Interceptify proxy engine to begin intercepting and analyzing network traffic.`,
	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetInt("port")
		address, _ := cmd.Flags().GetString("address")

		fmt.Printf("Starting Interceptify engine on %s:%d...\n", address, port)

		home, _ := os.UserHomeDir()
		caDir := home + "/.interceptify"
		os.MkdirAll(caDir, 0700)

		caCertPath := caDir + "/ca.crt"
		caKeyPath := caDir + "/ca.key"

		caInstance, err := ca.NewCA(caCertPath, caKeyPath)
		if err != nil {
			fmt.Printf("Failed to initialize CA: %v\n", err)
			return
		}

		proxyInstance := proxy.NewProxy(fmt.Sprintf("%s:%d", address, port), caInstance)

		// Register built-in plugins
		proxyInstance.Plugins.Register(&attack.LoggerPlugin{})

		if err := proxyInstance.Start(); err != nil {
			fmt.Printf("Failed to start proxy: %v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.Flags().IntP("port", "p", 8080, "Port to listen on")
	startCmd.Flags().StringP("address", "a", "127.0.0.1", "Address to bind to")
}
