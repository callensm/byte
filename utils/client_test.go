package utils

import (
	"net"
	"testing"
)

func TestNewClient(t *testing.T) {
	conn := new(net.Conn)
	c := NewClient(conn)

	if c == nil {
		t.Errorf("NewClient returned a null pointer instead of a *Client instance")
	} else {
		if c.Writer == nil {
			t.Errorf("The client's Writer was not correctly instantiated")
		}

		if c.Reader == nil {
			t.Errorf("The client's Reader was not correctly instantiated")
		}

		if c.conn != conn {
			t.Errorf("The input connection interface does not match the client's")
		}
	}
}

func TestPostAndFetch(t *testing.T) {
	msg := "abc"

	go func() {
		tempConn, err := net.Dial("tcp", ":4000")
		if err != nil {
			t.Fatal(err)
		}
		defer tempConn.Close()

		tempClient := NewClient(&tempConn)
		received := tempClient.Fetch('\n')

		if x := string(received); x != msg {
			t.Errorf("Posted %s but received %s", msg, x)
		}
	}()

	l, err := net.Listen("tcp", ":4000")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()

	conn, err := l.Accept()
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	client := NewClient(&conn)
	client.Post([]byte(msg + "\n"))
}
