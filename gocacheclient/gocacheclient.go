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
const MAX_PAYLOAD_LENGTH int = 2048

const GET_COMMAND string = "GET        "
const SET_COMMAND string = "SET        "
const DEL_COMMAND string = "DEL        "

var ErrInvalidResponseNoNewLine = errors.New("invalid response from server: No new Line")
var ErrInvalidResponseInvalidStatus = errors.New("invalid response from server: Invalid success code")
var ErrFailedResponse = errors.New("failed response")

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

func (c Client) MakeHeader(l int) string {
	h := fmt.Sprintf("%d\n", l)

	return h
}

func (c Client) Send(cmd string) (*Response, error) {
	conn, err := c.Connect()
	if err != nil {
		fmt.Println(err.Error())

		return nil, err
	}

	defer conn.Close()

	if len(cmd) > MAX_PAYLOAD_LENGTH {
		return nil, errors.New("cmd too large")
	}

	h := c.MakeHeader(len(cmd))
	cmd = h + cmd
	cmdBytes := []byte(cmd)

	n, err := conn.Write(cmdBytes)
	if n != len(cmd) {
		return nil, errors.New("Failed to send request")
	}

	buff := make([]byte, MAX_PAYLOAD_LENGTH)
	conn.Read(buff)

	r, err := c.ParseResponse(buff)

	return r, err
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
	padKey := fmt.Sprintf("%-*s ", KEY_LENGTH, key)
	cmd := SET_COMMAND + padKey + value
	cmd = fmt.Sprintf("%-*s", MAX_PAYLOAD_LENGTH, cmd)

	r, err = c.Send(cmd)

	return r, err
}

func (c Client) Get(key string) (*Response, error) {
	padKey := fmt.Sprintf("%-*s ", KEY_LENGTH, key)
	cmd := GET_COMMAND + padKey

	r, err := c.Send(cmd)

	return r, err
}

func (c Client) Del(key string) (*Response, error) {
	padKey := fmt.Sprintf("%-*s ", KEY_LENGTH, key)
	cmd := DEL_COMMAND + padKey

	r, err := c.Send(cmd)
	return r, err
}
