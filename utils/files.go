package utils

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
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

// CreateDirectories walks the Tree struct argued and
// creates all of the directories described in the struct
// for the files to eventuall be written to
func CreateDirectories(tree *Tree, base string) {
	path := filepath.Join(base, tree.Name)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, os.ModePerm)
		Catch(err)
	}

	for _, sub := range tree.SubTrees {
		CreateDirectories(sub, path)
	}
}

// Upload reads the argued file name
// and writes the data to the socket connection
// to be downloaded at the destination socket
func Upload(c *Client, localPath, destPath string) {
	logger.Warn(destPath)
	c.Post([]byte(destPath + "\n"))

	file, err := os.Open(localPath)
	defer file.Close()
	Catch(err)
	contents, err := ioutil.ReadAll(file)
	Catch(err)

	contents = append(contents, '\x00')
	c.Post(contents)
}

// UploadTree traverses an entire Tree structure and upload each
// file to the socket connection
func UploadTree(root, dir string, t *Tree, c *Client) {
	path := filepath.Join(dir, t.Name)
	fromRoot := filepath.Join(root, strings.Split(path, root)[1])

	for _, leaf := range t.Leaves {
		openPath := filepath.Join(path, leaf)
		Upload(c, openPath, filepath.Join(fromRoot, leaf))
		time.Sleep(50 * time.Millisecond)
	}

	for _, s := range t.SubTrees {
		UploadTree(root, filepath.Join(dir, t.Name), s, c)
	}
}

// Download reads all file data from the incoming
// socket buffer and rewrites the data into a new local
// file with the same name
func Download(c *Client, dir string) {
	name := c.Fetch('\n')
	fileName := string(name)

	// Open or create a new file with the argued file name in the destination directory
	full := filepath.Join(dir, fileName)
	newFile, err := os.Create(full)
	Catch(err)
	defer newFile.Close()

	data := c.Fetch('\x00')
	_, err = newFile.Write(data)
	Catch(err)
	logger.FileSuccess(fileName)
}

// fillBuffer returns a string that fills excess space
// with the ':' character into it is the desired byte length
func fillBuffer(data string, length int) string {
	diff := length - len(data)
	data += strings.Repeat(":", diff)
	return data
}
