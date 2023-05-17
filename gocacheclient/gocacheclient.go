package gocacheclient

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/dropdevrahul/gocache/protocol"
)

const COMMAND_LENGTH int = 11 // in bytes
const KEY_LENGTH int = 64     // in bytes
const MAX_PAYLOAD_LENGTH int = 2048

const GET_COMMAND string = "GET"
const SET_COMMAND string = "SET"
const DEL_COMMAND string = "DEL"
const SET_TTL_COMMAND string = "SET_TTL"
const GET_TTL_COMMAND string = "GET_TTL"

var ErrInvalidResponseNoNewLine = errors.New("invalid response from server: No new Line")
var ErrInvalidResponseInvalidStatus = errors.New("invalid response from server: Invalid success code")
var ErrFailedResponse = errors.New("failed response")

type Client struct {
	Host string
	Port string
}

type Response struct {
	Data []byte
}

func (c Client) Connect() (net.Conn, error) {
	conn, err := net.Dial("tcp", c.Host+":"+c.Port)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func (c Client) MakeHeader(l int) protocol.Header {
	h := protocol.Header{
		Len: l,
	}

	return h
}

func (c Client) Send(cmd string, key string, payload string) (*Response, error) {
	cmd = c.PadCmd(cmd)
	key = c.PadKey(key)

	conn, err := c.Connect()
	if err != nil {
		log.Println(err.Error())

		return nil, err
	}

	defer conn.Close()

	if len(cmd) > MAX_PAYLOAD_LENGTH {
		return nil, errors.New("cmd too large")
	}

	body := []byte(cmd + key + payload)
	h := c.MakeHeader(len(body))
	cmdBytes := append(h.ToBytes(), body...)

	n, err := conn.Write(cmdBytes)
	if n != len(cmdBytes) {
		return nil, errors.New("Failed to send request")
	}

	buff := bufio.NewReader(conn)
	r, err := c.ReadResponse(buff)

	return r, err
}

func (c Client) ReadResponse(b *bufio.Reader) (*Response, error) {
	r := &Response{}
	h := protocol.Header{}
	err := protocol.ReadHeaders(b, &h)
	if err != nil {
		return r, err
	}

	d, err := protocol.ReadBody(b, h.Len)
	if err != nil {
		return r, err
	}

	r.Data = d
	return r, nil
}

func (c Client) PadKey(key string) string {
	padKey := fmt.Sprintf("%-*s ", KEY_LENGTH, key)
	return padKey
}

func (c Client) PadCmd(cmd string) string {
	return fmt.Sprintf("%11s", cmd)
}

func (c Client) Set(key string, value string) (r *Response, err error) {
	r, err = c.Send(SET_COMMAND, key, value)
	return r, err
}

func (c Client) Get(key string) (*Response, error) {
	r, err := c.Send(GET_COMMAND, key, "")
	return r, err
}

func (c Client) Del(key string) (*Response, error) {
	r, err := c.Send(DEL_COMMAND, key, "")
	return r, err
}

func (c Client) ttlToString(ttl time.Duration) string {
	ttls := fmt.Sprintf("%d", int(ttl.Seconds()))
	return ttls
}

func (c Client) SetTTL(key string, ttl time.Duration) (*Response, error) {
	payload := c.ttlToString(ttl)
	r, err := c.Send(SET_TTL_COMMAND, key, payload)
	return r, err
}

func (c Client) GetTTL(key string) (time.Duration, *Response, error) {
	var ttl time.Duration
	r, err := c.Send(GET_TTL_COMMAND, key, "")
	if err == nil {
		i, err := strconv.Atoi(string(r.Data))
		if err != nil {
			return 0, r, err
		}
		ttl = time.Second * time.Duration(i)
		return ttl, r, err
	}

	return 0, r, err
}
