package commands

import (
	"fmt"
	"os"

	"github.com/callensm/byte/utils"
)

var logger = utils.NewLogger()

// Execute runs the root command for the byte CLI
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
