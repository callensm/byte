package utils

import (
	"bufio"
	"bytes"
	"net"
)

// Client is a wrapper for managing the reader
// and writer buffers created from a given network connection
type Client struct {
	Writer *bufio.Writer
	Reader *bufio.Reader
	conn   *net.Conn
}

// NewClient creates a new Client struct instance
// and converts the argued network connection struct
// into a buffered reader and writer
func NewClient(conn *net.Conn) *Client {
	return &Client{
		Writer: bufio.NewWriter(*conn),
		Reader: bufio.NewReader(*conn),
		conn:   conn,
	}
}

// Post writes bytes to a socket connection with
// the structs buffered writer instance
func (c *Client) Post(data []byte) {
	_, err := c.Writer.Write(data)
	Catch(err)

	err = c.Writer.Flush()
	Catch(err)
}

// Fetch reads bytes from the socket connection
// with the structs buffered reader instance
func (c *Client) Fetch(delim byte) []byte {
	b, err := c.Reader.ReadBytes(delim)
	Catch(err)
	return bytes.Trim(b, string(delim))
}
