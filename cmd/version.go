package cmd

import (
	"fmt"
	"runtime/debug"
	"time"

	"github.com/spf13/cobra"
)

// Requires Go 1.18+
var Commit = func() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				return setting.Value
			}
		}
	}
	return ""
}()

var (
	version    string = "0.1.0" // Semantic version number
	versionCmd        = &cobra.Command{
		Use:   "version",
		Short: "Version and build details",
		Long:  ``,
		Run: func(_ *cobra.Command, _ []string) {
			now := time.Now()
			fmt.Printf("%s-%s [%s]\n", version, Commit, now.Format(time.RFC822))
		},
	}
)

func init() {
	rootCmd.AddCommand(versionCmd)
}
