package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version    string = "v0.0.0"    // Semantic version number
	buildDate  string = "undefined" // Build date and time
	versionCmd        = &cobra.Command{
		Use:   "version",
		Short: "Version and build details",
		Long:  ``,
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Printf("%s [%s]\n", version, buildDate)
			return
		},
	}
)

func init() {
	rootCmd.AddCommand(versionCmd)
}
