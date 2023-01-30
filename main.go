package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/dropdevrahul/gocacheclient/gocacheclient"
)

func main() {
	h := flag.String("host", "localhost", "gocache server host")
	p := flag.String("port", "8888", "gocache server port")

	flag.Parse()

	c := gocacheclient.Client{
		Host: *h,
		Port: *p,
	}

	var inputs [3]string
	for {
		fmt.Printf("-> ")
		fmt.Scanln(&inputs[0], &inputs[1], &inputs[2])

		if strings.ToUpper(inputs[0]) == "GET" {
			r, err := c.Get(inputs[1])
			if err != nil {
				fmt.Println(err)
				continue
			}

			fmt.Println(string(r.Data))
		} else if strings.ToUpper(inputs[0]) == "SET" {
			r, err := c.Set(inputs[1], inputs[2])
			if err != nil {
				fmt.Println(err)
				continue
			}

			if !r.IsStatus() {
				fmt.Println("Request Failed " + r.Error)
			}
		} else {
			fmt.Println("Invalid command " + inputs[0])
		}

		inputs = [3]string{}
	}
}
