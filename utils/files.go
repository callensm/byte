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
	// FileCountBufferSize is the byte size of the number of files write
	FileCountBufferSize = 3
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
	// Open the file at the argued path and get the information for it
	file, err := os.Open(path)
	Catch(err)
	fileInfo, err := file.Stat()
	Catch(err)

	// Create the buffers for the size and name of the file to be write over TCP
	size := fillBuffer(strconv.FormatInt(fileInfo.Size(), 10), FileSizeBufferSize)
	name := fillBuffer(fileInfo.Name(), FileNameBufferSize)

	// Write each of the file information buffers
	conn.Write([]byte(size))
	conn.Write([]byte(name))

	// Create the file content buffer and continuously read the file into
	// until we have reached the end of the file, sending each chunk
	// over the connection before reading more
	fileDataBuffer := make([]byte, BufferSize)
	for {
		_, err := file.Read(fileDataBuffer)
		if err == io.EOF {
			break
		}
		conn.Write(fileDataBuffer)
	}
	logger.FileSuccess(fileInfo.Name())
}

// ReceiveFile reads all file data from the incoming
// socket buffer and rewrites the data into a new local
// file with the same name
func ReceiveFile(conn net.Conn, dir string) {
	// Create the file name and size buffers to read the socket into
	fileNameBuffer := make([]byte, FileNameBufferSize)
	fileSizeBuffer := make([]byte, FileSizeBufferSize)

	// Read the file size from the socket connection into its buffer and parse it
	conn.Read(fileSizeBuffer)
	fileSize, err := strconv.ParseInt(strings.Trim(string(fileSizeBuffer), ":"), 10, 64)
	Catch(err)

	// Read the file name from the socket connection into its buffer and parse it
	conn.Read(fileNameBuffer)
	fileName := strings.Trim(string(fileNameBuffer), ":")

	// Open or create a new file with the argued file name in the destination directory
	newFile, err := os.Create(filepath.Join(dir, fileName))
	Catch(err)
	defer newFile.Close()

	// Continuously copy the buffered file content chunks from the socket connection
	// into the newly created file until we have pipped all chunks into the file
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
	logger.FileSuccess(fileName)
}

// fillBuffer returns a string that fills excess space
// with the ':' character into it is the desired byte length
func fillBuffer(data string, length int) string {
	diff := length - len(data)
	data += strings.Repeat(":", diff)
	return data
}
