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
	portAddr := fmt.Sprintf(":%d", port)
	path, err := filepath.Abs(dir)
	utils.Catch(err)

	if !utils.IsDir(path) {
		logger.Error(fmt.Sprintf("The path %s is either not a directory or does not exist.", path))
	}

	server, err := net.Listen("tcp", portAddr)
	defer server.Close()
	utils.Catch(err)

	logger.Info(fmt.Sprintf("Listening for connections on %s", portAddr))

	conn, err := server.Accept()
	defer conn.Close()
	utils.Catch(err)
	logger.Info(fmt.Sprintf("Connection received from %s", conn.RemoteAddr()))

	numOfFilesBuffer := make([]byte, utils.FileCountBufferSize)
	conn.Read(numOfFilesBuffer)
	count, _ := strconv.ParseInt(string(numOfFilesBuffer), 10, 64)

	for x := 0; x < int(count); x++ {
		utils.ReceiveFile(conn, dir)
	}
}
