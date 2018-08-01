package commands

import (
	"fmt"
	"net"
	"path/filepath"
	"strings"

	"github.com/callensm/byte/utils"
	"github.com/spf13/cobra"
)

var addr string
var source string
var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send files over client connection to be downloaded",
	Run:   sendFunc,
}

func init() {
	sendCmd.Flags().StringVarP(&addr, "addr", "a", "", "Address to attempt to connect to [IP_ADRR:PORT] or [IP_ADDR] for default port")
	sendCmd.Flags().StringVarP(&source, "src", "s", "", "Path to the source of the file(s) being sent")
	sendCmd.MarkFlagRequired("addr")
	sendCmd.MarkFlagRequired("src")
	rootCmd.AddCommand(sendCmd)
}

func sendFunc(cmd *cobra.Command, args []string) {
	logger.Clear()

	// Add default port if not in the given flags
	address := strings.Split(addr, ":")
	if len(address) == 1 {
		address = append(address, "4500")
	}

	// Get the absolute path to the source of the content
	path, err := filepath.Abs(source)
	utils.Catch(err)

	// Ensure the path is an existing file or directory and form
	// the final address string to connect to
	if !utils.IsFile(path) && !utils.IsDir(path) {
		logger.Error(fmt.Sprintf("The path %s does not exist as a file or directory", path))
	}

	finalAddr := strings.Join(address, ":")
	utils.CreateSpinner(22, "green", fmt.Sprintf("Attempting to connect to %s", finalAddr), "connecting")

	// Connect through TCP to the destined address:port
	conn, err := net.Dial("tcp", finalAddr)
	defer conn.Close()
	utils.Catch(err)
	utils.RemoveSpinner("connecting", fmt.Sprintf("Connected to %s", finalAddr), true)

	// Create buffered reader and writer from the connection
	client := utils.NewClient(&conn)

	if utils.IsFile(path) {
		// Tell the receiver client to only expect
		// one file if the path pointer to a single file
		// and then send the file through the connection
		// TODO: Rewrite single file transfer to fit new protocol
		conn.Write([]byte("001"))
		utils.Upload(conn, path)
	} else {
		// Create the description file tree for the argued path
		utils.CreateSpinner(22, "blue", "Compiling a descriptive structure for the files to send", "send_struct")
		fileTree := utils.NewTree(path)
		treeEncoding := fileTree.String()

		// Send the encoded file tree JSON data to the receiver
		client.Post([]byte(treeEncoding + "\x00"))
		utils.RemoveSpinner("send_struct", "JSON file structure was sent!", true)

		// Attempt to receive approval to send files
		utils.CreateSpinner(22, "yellow", "Waiting for approval to send files...", "get_approval")
		approved := client.Fetch('\n')

		// Handle the received approval status for the file structure
		apprStatus := string(approved)
		if apprStatus == "n" {
			utils.RemoveSpinner("get_approval", "Your transfer request was denied.", false)
		}
		utils.RemoveSpinner("get_approval", "File transfer approved!", true)
	}
}
