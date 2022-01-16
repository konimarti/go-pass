package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{Use: "pass-cli"}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
