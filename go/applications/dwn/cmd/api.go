package cmd

import (
	"github.com/openreserveio/dwn/go/log"
	"github.com/spf13/cobra"
	"os"
)

// apiCmd represents the auth command
var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "OpenReserve DWN Public APIs",
	Long:  `APIs that confirm to the DIF DWN Specification`,
	Run: func(cmd *cobra.Command, args []string) {

		log.Info("Starting DWN APIs")
		// config, err := configuration.Config()
		if err != nil {
			log.Fatal("Configuration Fatal Error:  %v", err)
			os.Exit(1)
		}

	},
}

func init() {
	rootCmd.AddCommand(apiCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// orgapiCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// orgapiCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
