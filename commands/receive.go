package commands

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"path/filepath"

	"github.com/callensm/byte/utils"
	"github.com/spf13/cobra"
)

var dir string
var port uint32
var autoApprove bool
var receiveCmd = &cobra.Command{
	Use:   "receive",
	Short: "Download files sent over socket connection",
	Run:   receiveFunc,
}

func init() {
	receiveCmd.Flags().StringVarP(&dir, "dir", "d", "", "Directory to download received files to")
	receiveCmd.Flags().Uint32VarP(&port, "port", "p", 4500, "Port to listen for socket connections")
	receiveCmd.Flags().BoolVarP(&autoApprove, "auto-approve", "a", false, "Whether or not to require file structure approval")
	receiveCmd.MarkFlagRequired("dir")
	rootCmd.AddCommand(receiveCmd)
}

func receiveFunc(cmd *cobra.Command, args []string) {
	logger.Clear()

	// Create the port string to listen on
	// and the absolute path of the argued destination directory
	// for the files to be rewritten into
	portAddr := fmt.Sprintf(":%d", port)
	path, err := filepath.Abs(dir)
	utils.Catch(err)

	// Ensure that the argued path is an existing directory
	if !utils.IsDir(path) {
		logger.Error(fmt.Sprintf("The path %s is either not a directory or does not exist.", path))
	}

	// Create a TCP server to listen on the destined port
	server, err := net.Listen("tcp", portAddr)
	defer server.Close()
	utils.Catch(err)
	utils.CreateSpinner(22, "green", fmt.Sprintf("Listening for connections on %s", portAddr), "listening")

	// When a connection comes in, accept it if there aren't any errors
	conn, err := server.Accept()
	defer conn.Close()
	utils.Catch(err)
	utils.RemoveSpinner("listening", fmt.Sprintf("Connection received from %s", conn.RemoteAddr()), true)

	// Create a buffered reader writer for the connection and read the sent
	// JSON structure describing the file system that is about to be
	// received for download
	client := utils.NewClient(&conn)
	utils.CreateSpinner(22, "blue", "Waiting to receive file structure...", "get_struct")
	structure := client.Fetch('\x00')
	fileTree := utils.NewTreeFromJSON(bytes.TrimRight(structure, "\x00"))
	utils.RemoveSpinner("get_struct", "Received file structure being sent!", true)

	// Get the user's approval for the received file structure
	// exit gracefully if denied and continue if approved
	// after sending the other user the approval status
	appr := []byte("y")
	if !autoApprove {
		logger.Tree(fileTree.Display())
		res := logger.Prompt("Do you want to approve this tree (y/n): ")
		appr = []byte(res)
	}

	client.Post(appr)
	if string(appr) == "n" {
		os.Exit(0)
	}

	// Use the generated file tree sent over the connection to create
	// directories that do not exist for files to be written in
	utils.CreateSpinner(22, "blue", "Creating file directory structure...", "make_dirs")
	utils.CreateDirectories(fileTree, dir)
	utils.RemoveSpinner("make_dirs", "File directories created!", true)
}
