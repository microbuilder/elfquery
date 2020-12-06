package cmd

import (
	"github.com/microbuilder/goelf/httpserver"
	"github.com/spf13/cobra"
)

// httpCmd represents the http command
var httpCmd = &cobra.Command{
	Use:   "http [filename]",
	Short: "HTTP based file analysis",
	Long: `Starts up an HTTP server instance that can be used to perform
detailed analysis of the specified ELF file.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetInt16("port")
		httpserver.Start(port)
	},
}

func init() {
	rootCmd.AddCommand(httpCmd)

	// Allow a custom port number
	httpCmd.PersistentFlags().Int16P("port", "p", 1443, "Port number")
}
