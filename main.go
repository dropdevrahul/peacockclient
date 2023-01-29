package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/dropdevrahul/gocacheclient/gocacheclient"
)

func main() {
	flag.String("h", "localhost", "gocache server host")
	flag.String("p", "8888", "gocache server port")

	flag.Parse()

	c := gocacheclient.Client{
		Host: "localhost",
		Port: "8888",
	}

	var inputs [3]string
	for {
		fmt.Printf("-> ")
		fmt.Scanln(&inputs[0], &inputs[1], &inputs[2])

		if strings.ToUpper(inputs[0]) == "GET" {
			val, err := c.Get(inputs[1])
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(val)
		} else if strings.ToUpper(inputs[0]) == "SET" {
			err := c.Set(inputs[1], inputs[2])
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Invalid command " + inputs[0])
		}

		inputs = [3]string{}
	}
}
