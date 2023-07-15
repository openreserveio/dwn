package cmd

import (
	"github.com/openreserveio/dwn/go/applications/dwn/configuration"
	"github.com/openreserveio/dwn/go/applications/dwn/service/hooksvc"
	"github.com/openreserveio/dwn/go/log"
	"github.com/spf13/cobra"
	"os"
)

// apiCmd represents the auth command
var hooksvcCmd = &cobra.Command{
	Use:   "hooksvc",
	Short: "OpenReserve DWN Backend HookService",
	Long:  `Backend gRPC Service for managing Hook definitions`,
	Run: func(cmd *cobra.Command, args []string) {

		log.Info("Starting DWN Backend HookService")
		config, err := configuration.Config()
		if err != nil {
			log.Fatal("Configuration Fatal Error:  %v", err)
			os.Exit(1)
		}
		log.Error("Stopping Hook Service:  %v", hooksvc.Start(config))

	},
}

func init() {
	rootCmd.AddCommand(hooksvcCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// orgapiCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// orgapiCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
