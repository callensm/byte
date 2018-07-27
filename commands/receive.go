package commands

import (
	"fmt"
	"net"
	"path/filepath"

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
	path, err := filepath.Abs(dir)
	if err != nil {
		logger.Error(err.Error())
	} else if !utils.IsDir(path) {
		logger.Error(fmt.Sprintf("The path %s is either not a directory or does not exist.", path))
	}

	portAddr := fmt.Sprintf(":%d", port)
	ln, _ := net.Listen("tcp", portAddr)
	defer ln.Close()
	logger.Info(fmt.Sprintf("Listening for connections on %s", portAddr))

	conn, _ := ln.Accept()
	defer conn.Close()

	logger.Info(fmt.Sprintf("Connection received from %s", conn.RemoteAddr()))
}
