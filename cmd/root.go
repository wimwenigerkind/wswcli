package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var rootCmd = &cobra.Command{
	Use:   "PLACEHOLDER_USE",
	Short: "PLACEHOLDER_SHORT",
	Long:  "PLACEHOLDER_LONG",
}

func SetVersionInfo(v, c, d string) {
	version = v
	commit = c
	date = d
	rootCmd.Version = fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, date)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
