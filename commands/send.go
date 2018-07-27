package commands

import (
	"fmt"
	"io/ioutil"
	"net"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/callensm/byte/utils"
	"github.com/spf13/cobra"
)

var addr string
var source string

func init() {
	sendCmd.Flags().StringVarP(&addr, "addr", "a", "", "Address to attempt to connect to [IP_ADRR:PORT] or [IP_ADDR] for default port")
	sendCmd.Flags().StringVarP(&source, "src", "s", "", "Path to the source of the file(s) being sent")
	sendCmd.MarkFlagRequired("addr")
	sendCmd.MarkFlagRequired("src")
	rootCmd.AddCommand(sendCmd)
}

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send files over socket connection to be downloaded",
	Run:   sendFunc,
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
	utils.RemoveSpinner("connecting", fmt.Sprintf("Connected to %s", finalAddr))

	if utils.IsFile(path) {
		// Tell the receiver socket to only expect
		// one file if the path pointer to a single file
		// and then send the file through the connection
		conn.Write([]byte("001"))
		utils.Upload(conn, path)
	} else {
		// Get the list of all files in the argued
		// directory and create the string that indicates
		// the number of files in that directory
		fileList, _ := ioutil.ReadDir(path)
		lenStr := strconv.Itoa(len(fileList))
		countMsg := strings.Repeat("0", utils.FileCountBufferSize-len(lenStr)) + lenStr

		// Tell the receiver socket how many files to expect
		// could be from 1-999 files, and then loop through the
		// file list and sychronously send each through the connection
		logger.Directory(len(fileList), path, true)
		conn.Write([]byte(countMsg))
		for _, f := range fileList {
			utils.Upload(conn, filepath.Join(path, f.Name()))
			time.Sleep(100 * time.Millisecond)
		}
	}
}
