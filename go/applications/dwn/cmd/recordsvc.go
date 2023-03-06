package cmd

import (
	"github.com/openreserveio/dwn/go/applications/dwn/configuration"
	"github.com/openreserveio/dwn/go/applications/dwn/service/recordsvc"
	"github.com/openreserveio/dwn/go/log"
	"github.com/openreserveio/dwn/go/observability"
	"github.com/spf13/cobra"
	"os"
)

// recordsvcCmd represents the recordsvcCmd command
var recordsvcCmd = &cobra.Command{
	Use:   "recordsvc",
	Short: "OpenReserve DWN Backend RecordService",
	Long:  `Backend gRPC Service for managing records, schemas, and schema definitions`,
	Run: func(cmd *cobra.Command, args []string) {

		log.Info("Observability")
		ctx := cmd.Context()
		sd, _ := observability.InitProviderWithJaegerExporter(ctx, "Record Service")
		defer sd(ctx)

		log.Info("Starting DWN Backend RecordService")
		config, err := configuration.Config()
		if err != nil {
			log.Fatal("Configuration Fatal Error:  %v", err)
			os.Exit(1)
		}
		log.Error("Stopping Record Service:  %v", recordsvc.Start(ctx, config))

	},
}

func init() {
	rootCmd.AddCommand(recordsvcCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// orgapiCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// orgapiCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
