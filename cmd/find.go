package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/konimarti/go-pass/internal/config"
	"github.com/konimarti/go-pass/internal/fcts"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(findCmd)
}

var findCmd = &cobra.Command{
	Use:     "find pass-names...",
	Short:   "List passwords that match pass-names.",
	Long:    ``,
	Aliases: []string{"search"},
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Search terms:", strings.Join(args, ","))
		fcts.TreeFind(os.Stdout, config.New().Prefix, args)
	},
}
