package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/konimarti/go-pass/internal/config"
	"github.com/konimarti/go-pass/internal/fcts"
	"github.com/konimarti/go-pass/internal/gpg"
	"github.com/spf13/cobra"
)

var (
	clip   = ""
	qrcode = ""
)

func init() {
	rootCmd.AddCommand(showCmd)
	showCmd.Flags().StringVarP(&clip, "clip", "c", "", "")
	showCmd.Flags().StringVarP(&qrcode, "qrcode", "q", "", "")
}

var showCmd = &cobra.Command{
	Use:     "show [pass-name]",
	Short:   "Show existing password.",
	Long:    ``,
	Aliases: []string{"ls", "list"},
	Args:    cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.New()
		var path string
		if len(args) > 0 {
			path = fcts.CheckSneakyPath(args[0])
		}
		passfile := filepath.Join(cfg.Prefix, fmt.Sprintf("%s.gpg", path))
		passdir := filepath.Join(cfg.Prefix, path)
		if ok, err := fcts.IsFile(passfile); err == nil && ok {
			if clip != "" {
				fmt.Println("clip argument for show not implemented")
			}
			if qrcode != "" {
				fmt.Println("qrcode argument for show not implemented")
			}
			cipher, err := os.Open(passfile)
			if err != nil {
				fmt.Println(err)
				return
			}
			plain, err := gpg.Decrypt(cipher)
			if err != nil {
				fmt.Println(err)
				return
			}
			io.Copy(os.Stdout, plain)

		} else if ok, err := fcts.IsDir(passdir); err == nil && ok {
			var buf bytes.Buffer
			err := fcts.Tree(&buf, passdir)
			if err != nil {
				fmt.Println(err)
				return
			}
			io.Copy(os.Stdout, &buf)
		} else if path == "" {
			fmt.Println("Error: password store is empty. Try \"", filepath.Base(os.Args[0]), " init\".")
		} else {
			fmt.Println("Error:", path, "is not in password store.")
		}
		fmt.Println("")

	},
}
