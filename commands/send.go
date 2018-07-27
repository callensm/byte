package commands

import (
	"fmt"
	"io/ioutil"
	"net"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/callensm/byte/utils"
	"github.com/spf13/cobra"
)

var addr string
var source string

func init() {
	sendCmd.Flags().StringVarP(&addr, "addr", "a", "", "Address to attempt to connect to [IP_ADRR:PORT] or [IP_ADDR] for default port")
	sendCmd.Flags().StringVarP(&source, "source", "s", "", "Path to the source of the file(s) being sent")
	sendCmd.MarkFlagRequired("addr")
	sendCmd.MarkFlagRequired("source")
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

	path, err := filepath.Abs(source)
	utils.Catch(err)

	if !utils.IsFile(path) && !utils.IsDir(path) {
		logger.Error(fmt.Sprintf("The path %s does not exist as a file or directory", path))
	}

	finalAddr := strings.Join(address, ":")
	logger.Info(fmt.Sprintf("Attempting to connect to %s", finalAddr))

	conn, err := net.Dial("tcp", finalAddr)
	defer conn.Close()
	utils.Catch(err)
	logger.Info(fmt.Sprintf("Connected to %s", finalAddr))

	if utils.IsFile(path) {
		conn.Write([]byte("001"))
		utils.SendFile(conn, path)
	} else {
		fileList, _ := ioutil.ReadDir(path)
		lenStr := strconv.Itoa(len(fileList))
		countMsg := strings.Repeat("0", utils.FileCountBufferSize-len(lenStr)) + lenStr
		conn.Write([]byte(countMsg))
		for _, f := range fileList {
			utils.SendFile(conn, filepath.Join(path, f.Name()))
		}
	}
}
