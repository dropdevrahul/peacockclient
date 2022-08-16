package gocacheclient

import (
	"errors"
	"fmt"
	"net"
)

const COMMAND_LENGTH int = 11            // in bytes
const KEY_LENGTH int = 64                // in bytes
const MAX_PAYLOAD_LENGTH int = 10 * 1024 // max 10 KB
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

func (c Client) GetContentLenHeader(n int) []byte {
	res := fmt.Sprintf("%s%d\n", HEADER_CONTENT_LENGTH_NAME, n)
	return []byte(res)
}

func (c Client) Set(key string, value string) error {
	conn, err := c.Connect()
	if err != nil {
		fmt.Println(err.Error())
	}
	padKey := fmt.Sprintf("%-*s ", KEY_LENGTH, key)
	if len(value) >= MAX_PAYLOAD_LENGTH {
		return errors.New("value too large")
	}
	cmd := SET_COMMAND + padKey + value
	h := c.GetContentLenHeader(len([]byte(cmd)))
	_, err = conn.Write(h)

	buff := make([]byte, 1024)
	conn.Read(buff)
	fmt.Println(string(buff))

	cmdBytes := []byte(cmd)
	cmdBytes = append(cmdBytes, NULL_BYTE)

	n, err := conn.Write(cmdBytes)
	fmt.Println(n)

	conn.Read(buff)
	fmt.Println(string(buff))
	defer conn.Close()
	//for sentBytes < totlen {
	//	fmt.Println("%d", sentBytes)
	//	n, _ := conn.Write(cmdBytes[sentBytes:totlen])
	//	sentBytes += n
	//}
	return err
}

func (c Client) GET(key string) (string, error) {
	conn, err := c.Connect()
	if err != nil {
		fmt.Println(err.Error())
	}

	padKey := fmt.Sprintf("%-*s ", KEY_LENGTH, key)
	cmd := GET_COMMAND + padKey
	h := c.GetContentLenHeader(len([]byte(cmd)))
	_, err = conn.Write(h)

	buff := make([]byte, 1024)
	conn.Read(buff)

	cmdBytes := []byte(cmd)
	cmdBytes = append(cmdBytes, NULL_BYTE)

	_, err = conn.Write(cmdBytes)

	conn.Read(buff)
	value := string(buff)

	defer conn.Close()

	return value, err
	//for sentBytes < totlen {
	//	fmt.Println("%d", sentBytes)
	//	n, _ := conn.Write(cmdBytes[sentBytes:totlen])
	//	sentBytes += n
	//}
}
