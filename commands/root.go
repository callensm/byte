package commands

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "byte",
	Short: "Byte is a file transfer CLI over secure socket connections",
	Run:   rootFunc,
}

func rootFunc(cmd *cobra.Command, args []string) {}
