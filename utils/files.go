package utils

import (
	"io"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	// FileNameBufferSize is the byte size of the file name write
	FileNameBufferSize = 64
	// FileSizeBufferSize is the byte size of the file size write
	FileSizeBufferSize = 10
	// BufferSize is the byte size of the buffer used to transmit file data
	BufferSize = 1024
)

// IsDir returns whether the argued path string
// is a valid and existing directory
func IsDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.Mode().IsDir()
}

// IsFile returns whether the argued path string
// is a valid and existing file
func IsFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.Mode().IsRegular()
}

// SendFile reads the argued file name
// and writes the data to the socket connection
// to be downloaded at the destination socket
func SendFile(conn net.Conn, path string) {
	defer conn.Close()
	file, err := os.Open(path)
	Catch(err)

	fileInfo, err := file.Stat()
	Catch(err)

	size := fillBuffer(strconv.FormatInt(fileInfo.Size(), 10), FileSizeBufferSize)
	name := fillBuffer(fileInfo.Name(), FileNameBufferSize)

	conn.Write([]byte(size))
	conn.Write([]byte(name))

	fileDataBuffer := make([]byte, BufferSize)
	for {
		_, err := file.Read(fileDataBuffer)
		if err == io.EOF {
			break
		}
		conn.Write(fileDataBuffer)
	}
}

// ReceiveFile reads all file data from the incoming
// socket buffer and rewrites the data into a new local
// file with the same name
func ReceiveFile(conn net.Conn, dir string) {
	defer conn.Close()
	fileNameBuffer := make([]byte, FileNameBufferSize)
	fileSizeBuffer := make([]byte, FileSizeBufferSize)

	conn.Read(fileSizeBuffer)
	fileSize, err := strconv.ParseInt(strings.Trim(string(fileSizeBuffer), ":"), 10, 64)
	Catch(err)

	conn.Read(fileNameBuffer)
	fileName := strings.Trim(string(fileNameBuffer), ":")

	newFile, err := os.Create(filepath.Join(dir, fileName))
	Catch(err)
	defer newFile.Close()

	var received int64
	for {
		if (fileSize - received) < BufferSize {
			io.CopyN(newFile, conn, (fileSize - received))
			conn.Read(make([]byte, (received + BufferSize - fileSize)))
			break
		}
		io.CopyN(newFile, conn, BufferSize)
		received += BufferSize
	}
}

func fillBuffer(data string, length int) string {
	diff := length - len(data)
	data += strings.Repeat(":", diff)
	return data
}
