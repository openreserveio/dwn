package cmd

import (
	"context"
	"github.com/openreserveio/dwn/go/applications/dwn/configuration"
	"github.com/openreserveio/dwn/go/applications/dwn/service/keysvc"
	"github.com/openreserveio/dwn/go/log"
	"github.com/openreserveio/dwn/go/observability"
	"github.com/spf13/cobra"
	"os"
)

// keysvcCmd represents the keysvc command
var keysvcCmd = &cobra.Command{
	Use:   "keysvc",
	Short: "OpenReserve DWN Backend Key Service",
	Long:  `Backend gRPC Service for managing keys, authenticating messages, and signing messages`,
	Run: func(cmd *cobra.Command, args []string) {

		log.Info("Observability")
		observability.InitProviderWithJaegerExporter(context.Background(), "Key Service")

		log.Info("Starting DWN Backend Key Service")
		config, err := configuration.Config()
		if err != nil {
			log.Fatal("Configuration Fatal Error:  %v", err)
			os.Exit(1)
		}
		log.Error("Stopping Key Service:  %v", keysvc.Start(config))

	},
}

func init() {
	rootCmd.AddCommand(keysvcCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// orgapiCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// orgapiCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
