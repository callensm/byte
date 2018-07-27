package commands

import (
	"fmt"
	"net"
	"path/filepath"
	"strconv"

	"github.com/callensm/byte/utils"
	"github.com/spf13/cobra"
)

var dir string
var port uint32

func init() {
	receiveCmd.Flags().StringVarP(&dir, "dir", "d", "", "Directory to download received files to")
	receiveCmd.Flags().Uint32VarP(&port, "port", "p", 4500, "Port to listen for socket connections")
	receiveCmd.MarkFlagRequired("dir")
	rootCmd.AddCommand(receiveCmd)
}

var receiveCmd = &cobra.Command{
	Use:   "receive",
	Short: "Download files sent over socket connection",
	Run:   receiveFunc,
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
	utils.RemoveSpinner("listening", fmt.Sprintf("Connection received from %s", conn.RemoteAddr()))

	// Read in the initial message from the client socket
	// indicating how many files to expect to be written to the
	// socket connection for transfer
	numOfFilesBuffer := make([]byte, utils.FileCountBufferSize)
	conn.Read(numOfFilesBuffer)
	count, _ := strconv.ParseInt(string(numOfFilesBuffer), 10, 64)

	// Loop `count` times to accept all expected inbound files
	// to be rewritten to the destination path
	logger.Directory(int(count), path, false)
	for x := 0; x < int(count); x++ {
		utils.ReceiveFile(conn, dir)
	}
}
