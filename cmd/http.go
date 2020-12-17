package cmd

import (
	"fmt"

	"github.com/microbuilder/elfquery/elf2sql"
	"github.com/microbuilder/elfquery/httpserver"
	"github.com/spf13/cobra"
)

// httpCmd represents the http command
var httpCmd = &cobra.Command{
	Use:   "http filename",
	Short: "HTTP based file analysis",
	Long: `Starts up an HTTP server instance that can be used to perform
detailed analysis of the specified ELF file.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Populate the database with the ELF data
		e := elf2sql.InitDB(args[0])
		defer elf2sql.CloseDB()
		if e != nil {
			fmt.Printf("unable to initialise the SQLite3 database in memory\n")
			return
		}

		// Start thee HTTP server
		port, _ := cmd.Flags().GetInt16("port")
		httpserver.Start(port)
	},
}

func init() {
	rootCmd.AddCommand(httpCmd)

	// Allow a custom port number
	httpCmd.PersistentFlags().Int16P("port", "p", 1443, "Port number")
}
