package commands

import (
	"fmt"

	"github.com/callensm/byte/utils"
	"github.com/spf13/cobra"
)

var root string
var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "Inspect the file tree from the source directory argued",
	Run:   viewFunc,
}

func init() {
	viewCmd.Flags().StringVarP(&root, "root", "r", "", "The directory to begin the file tree")
	viewCmd.MarkFlagRequired("root")
	rootCmd.AddCommand(viewCmd)
}

func viewFunc(cmd *cobra.Command, args []string) {
	// Validate that the argued root is a directory
	if !utils.IsDir(root) {
		logger.Error(fmt.Sprintf("%s is not a valid or existing directory", root))
	}

	// Create the Tree structure instance from the argued root
	tree := utils.NewTree(root)

	// Print the indented/formatted version of the marshalled JSON
	logger.Info(fmt.Sprintf("Sub-directories: %d", tree.CountSubTrees()))
	logger.Info(fmt.Sprintf("Files: %d", tree.CountLeaves()))
	logger.Tree(tree.Display())
}
