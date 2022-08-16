package main

import (
	"github.com/dropdevrahul/gocache-go-client/gocacheclient"
)

func main() {
	client := gocacheclient.Client{
		host: "localhost",
		port: "8888",
	}
	client.Connect()
	client.Set("A", "234")
}
