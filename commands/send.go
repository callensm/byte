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
var file string

func init() {
	sendCmd.Flags().StringVarP(&addr, "addr", "a", "", "Address to attempt to connect to [IP_ADRR:PORT] or [IP_ADDR] for default port")
	sendCmd.Flags().StringVarP(&file, "file", "f", "", "Path to the file being sent")
	sendCmd.MarkFlagRequired("addr")
	sendCmd.MarkFlagRequired("file")
	rootCmd.AddCommand(sendCmd)
}

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send files over socket connection to be downloaded",
	Run:   sendFunc,
}

func sendFunc(cmd *cobra.Command, args []string) {
	address := strings.Split(addr, ":")
	if len(address) == 1 {
		address = append(address, "4500")
	}

	path, err := filepath.Abs(file)
	if err != nil {
		logger.Error(err.Error())
	} else if !utils.IsFile(path) {
		logger.Error(fmt.Sprintf("The path %s is either not a file or does not exist", path))
	}

	finalAddr := strings.Join(address, ":")
	logger.Info(fmt.Sprintf("Attempting to connect to %s", finalAddr))

	conn, _ := net.Dial("tcp", finalAddr)
	defer conn.Close()

	logger.Info("Connected to destination socket!")
}
