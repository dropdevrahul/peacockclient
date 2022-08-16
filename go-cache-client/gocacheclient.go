package gocacheclient

import (
	"fmt"
	"net"
)

const COMMAND_LENGTH int = 11            // in bytes
const KEY_LENGTH int = 64                // in bytes
const MAX_PAYLOAD_LENGTH int = 10 * 1024 // max 10 KB
const NULL_BYTE byte = 0
const GET_COMMAND []byte = []byte("Get        ")
const HEADER_CONTENT_LENGTH_NAME = "CONTENT-LENGTH:"

type Client struct {
	host string
	port string
	conn net.Conn
}

func (c Client) Connect() error {
	conn, err := net.Dial("tcp", c.host+":"+c.port)
	if err != nil {
		return err
	}
	c.conn = conn
}

func (c Client) GetContentLenHeader(n int) (string, error) {
	res, err := fmt.Sprintf("%s%d\n", HEADER_CONTENT_LENGTH_NAME, n)
}

func (c Client) Set(key string, value string) error {

	padKey := fmt.Sprintf("%-*s ", KEY_LENGTH, key)
	if len(value) >= MAX_PAYLOAD_LENGTH {
		return error.Error("value too large")
	}
	valBytes := []byte(value)
	valBytes += NULL_BYTE
	cmd := GET_COMMAND + padKey + valBytes
	h := c.GetContentLenHeader(len(cmd))
	cmd = h + cmd
	_, err := c.conn.Write(cmd)
	return err
}
