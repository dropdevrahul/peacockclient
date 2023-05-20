package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/dropdevrahul/peacockclient"
)

func main() {
	h := flag.String("host", "localhost", "gocache server host")
	p := flag.String("port", "9999", "gocache server port")

	flag.Parse()

	c := peacockclient.Client{
		Host: *h,
		Port: *p,
	}

	var inputs [3]string
	for {
		fmt.Printf("-> ")
		fmt.Scanln(&inputs[0], &inputs[1], &inputs[2])

		switch strings.ToUpper(inputs[0]) {
		case "GET":
			r, err := c.Get(inputs[1])
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(string(r.Data))
		case "SET":
			r, err := c.Set(inputs[1], inputs[2])
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(string(r.Data))
		case "DEL":
			r, err := c.Del(inputs[1])
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(string(r.Data))
		case "SET_TTL":
			ttl, err := strconv.Atoi(inputs[2])
			if err != nil {
				fmt.Println(err)
			}

			r, err := c.SetTTL(inputs[1], time.Second*time.Duration(ttl))
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(string(r.Data))
		case "GET_TTL":
			ttl, _, err := c.GetTTL(inputs[1])
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(ttl)
		default:
			fmt.Println("Invalid command " + inputs[0])
		}
		inputs = [3]string{}
	}
}
