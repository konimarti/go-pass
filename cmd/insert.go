package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/konimarti/go-pass/internal/config"
	"github.com/konimarti/go-pass/internal/fcts"
	"github.com/konimarti/go-pass/internal/gpg"
	"github.com/spf13/cobra"
)

var (
	multiline = false
	echo      = true
	force     = false
)

func init() {
	rootCmd.AddCommand(insertCmd)
	insertCmd.Flags().BoolVarP(&multiline, "multiline", "m", false, "")
	insertCmd.Flags().BoolVarP(&echo, "echo", "e", true, "")
	insertCmd.Flags().BoolVarP(&force, "force", "f", false, "")
}

var insertCmd = &cobra.Command{
	Use:   "insert [--echo,-e | --multiline,-m] [--force,-f] pass-name",
	Short: "Insert new password.",
	Long: `Insert new password. Optionally, echo the password back to the console
	        during entry. Or, optionally, the entry may be multiline. Prompt before
	        overwriting existing password unless forced.`,

	Aliases: []string{"add"},
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := fcts.CheckSneakyPath(args[0])

		// TODO: we should check for / at end of path before appending extension
		passfile := filepath.Join(config.New().Prefix, fmt.Sprintf("%s.gpg", path))
		passdir, _ := filepath.Split(passfile)

		if !echo {
			panic("noecho is currently not implemented yet")
		}

		if _, err := fcts.IsFile(passfile); err == nil && !force {
			// fmt.Println("An entry already exists for $path. Overwrite it?")
			// TODO: ask if you want to overwrite file

			fmt.Println("file exists; use --force to overwrite")
			return
		}

		err := os.MkdirAll(passdir, 0744)
		if err != nil {
			fmt.Println(err)
			return
		}

		var rcpts []string = config.New().GpgRecipientsId
		if len(rcpts) == 0 {
			// FIXME: assuming gpg-id file is in prefix dir
			gpgidFile := filepath.Join(config.New().Prefix, ".gpg-id")
			if !fcts.FileExists(gpgidFile) {
				fmt.Fprintln(os.Stderr, "Error: You must run: PROG init your-gpg-id before you may use the password store.gpgidFile")
			}

			// TODO: implement signing of gpg-id file

			f, err := os.Open(gpgidFile)
			if err != nil {
				fmt.Println(err)
				return
			}
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				rcpts = append(rcpts, strings.Split(scanner.Text(), "#")[0])
			}
		}

		fmt.Println("Enter contents of", path, "and press Ctrl+D when finished:")
		scanner := bufio.NewScanner(os.Stdin)
		var text []string
		for scanner.Scan() {
			text = append(text, scanner.Text())
			if !multiline {
				break
			}
		}

		plain := strings.NewReader(strings.Join(text, "\n"))
		cipher, err := gpg.Encrypt(rcpts, plain)
		if err != nil {
			fmt.Println(err)
			return
		}

		f, err := os.OpenFile(passfile, os.O_RDWR|os.O_CREATE, 0600)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()

		io.Copy(f, cipher)
	},
}
