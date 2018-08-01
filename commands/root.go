package commands

import (
	"github.com/callensm/byte/utils"
	"github.com/spf13/cobra"
)

var logger = utils.NewLogger()

var rootCmd = &cobra.Command{
	Use:   "byte",
	Short: "Byte is a file transfer CLI over secure socket connections",
	Run:   rootFunc,
}

func rootFunc(cmd *cobra.Command, args []string) {}

// Execute runs the root command for the byte CLI
func Execute() {
	err := rootCmd.Execute()
	utils.Catch(err)
}
