package gocacheclient

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
)

const COMMAND_LENGTH int = 11 // in bytes
const KEY_LENGTH int = 64     // in bytes
const MAX_PAYLOAD_LENGTH int = 1468
const NULL_BYTE byte = 0

var GET_COMMAND string = "GET        "
var SET_COMMAND string = "SET        "
var ErrInvalidResponseNoNewLine = errors.New("invalid response from server: No new Line")
var ErrInvalidResponseInvalidStatus = errors.New("invalid response from server: Invalid success code")
var ErrFailedResponse = errors.New("failed response")

const HEADER_CONTENT_LENGTH_NAME = "CONTENT-LENGTH:"

type Client struct {
	Host string
	Port string
}

type Response struct {
	Status int // 1 success rest fail 0 undefined
	Error  string
	Data   []byte
}

func (r *Response) IsStatus() bool {
	if r.Status == 1 {
		return true
	}

	return false
}

func (c Client) Connect() (net.Conn, error) {
	conn, err := net.Dial("tcp", c.Host+":"+c.Port)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (c Client) ParseResponse(b []byte) (r *Response, err error) {
	r = &Response{}
	h, d, ok := strings.Cut(string(b), "\n")
	if !ok {
		return r, ErrInvalidResponseNoNewLine
	}

	success, errMsg, ok := strings.Cut(h, " ")
	successCode, err := strconv.Atoi(success)
	if !ok || err != nil {
		return r, ErrInvalidResponseInvalidStatus
	}

	r.Error = errMsg
	r.Status = successCode
	r.Data = []byte(d)

	return r, nil
}

func (c Client) Set(key string, value string) (r *Response, err error) {
	conn, err := c.Connect()
	if err != nil {
		fmt.Println(err.Error())
	}

	//conn.Read(buff)
	padKey := fmt.Sprintf("%-*s ", KEY_LENGTH, key)
	cmd := SET_COMMAND + padKey + value
	cmd = fmt.Sprintf("%-*s", MAX_PAYLOAD_LENGTH, cmd)

	if len(cmd) > MAX_PAYLOAD_LENGTH {
		return r, errors.New("value too large")
	}

	buff := make([]byte, MAX_PAYLOAD_LENGTH)

	cmdBytes := []byte(cmd)

	_, err = conn.Write(cmdBytes)

	defer conn.Close()
	conn.Read(buff)

	r, err = c.ParseResponse(buff)

	return r, err
}

func (c Client) Get(key string) (*Response, error) {
	buff := make([]byte, 1024)
	conn, err := c.Connect()
	if err != nil {
		fmt.Println(err.Error())

		return nil, err
	}

	padKey := fmt.Sprintf("%-*s ", KEY_LENGTH, key)
	cmd := GET_COMMAND + padKey
	cmdBytes := []byte(cmd)

	cmdBytes = append(cmdBytes, NULL_BYTE)

	_, err = conn.Write(cmdBytes)

	conn.Read(buff)

	r, err := c.ParseResponse(buff)
	if err != nil || !r.IsStatus() {
		fmt.Printf("Status Code: %d Err: %s\n", r.Status, r.Error)
		return nil, ErrFailedResponse
	}

	defer conn.Close()

	return r, err
}
