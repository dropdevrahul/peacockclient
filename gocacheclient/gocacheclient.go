package gocacheclient

import (
	"errors"
	"fmt"
	"net"
)

const COMMAND_LENGTH int = 11 // in bytes
const KEY_LENGTH int = 64     // in bytes
const MAX_PAYLOAD_LENGTH int = 1468
const NULL_BYTE byte = 0

var GET_COMMAND string = "GET        "
var SET_COMMAND string = "SET        "

const HEADER_CONTENT_LENGTH_NAME = "CONTENT-LENGTH:"

type Client struct {
	Host string
	Port string
}

func (c Client) Connect() (net.Conn, error) {
	conn, err := net.Dial("tcp", c.Host+":"+c.Port)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (c Client) Set(key string, value string) error {
	conn, err := c.Connect()
	if err != nil {
		fmt.Println(err.Error())
	}

	//conn.Read(buff)
	padKey := fmt.Sprintf("%-*s ", KEY_LENGTH, key)
	cmd := SET_COMMAND + padKey + value
	cmd = fmt.Sprintf("%-*s", MAX_PAYLOAD_LENGTH, cmd)

	if len(cmd) > MAX_PAYLOAD_LENGTH {
		return errors.New("value too large")
	}

	buff := make([]byte, MAX_PAYLOAD_LENGTH)

	cmdBytes := []byte(cmd)

	_, err = conn.Write(cmdBytes)

	defer conn.Close()
	conn.Read(buff)

	return err
}

func (c Client) Get(key string) (string, error) {
	conn, err := c.Connect()
	if err != nil {
		fmt.Println(err.Error())

		return "", err
	}

	padKey := fmt.Sprintf("%-*s ", KEY_LENGTH, key)
	cmd := GET_COMMAND + padKey

	buff := make([]byte, 1024)

	cmdBytes := []byte(cmd)
	cmdBytes = append(cmdBytes, NULL_BYTE)

	_, err = conn.Write(cmdBytes)

	conn.Read(buff)
	value := string(buff)

	defer conn.Close()

	return value, err
}
