package commands

import (
	"github.com/callensm/byte/utils"
)

var logger = utils.NewLogger()

// Execute runs the root command for the byte CLI
func Execute() {
	err := rootCmd.Execute()
	utils.Catch(err)
}
