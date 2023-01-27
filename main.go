package main

import (
	"fmt"

	"github.com/dropdevrahul/gocacheclient/gocacheclient"
)

func main() {
	client := gocacheclient.Client{
		Host: "localhost",
		Port: "8888",
	}
	i := 0
	for i < 10 {
		i += 1
		err := client.Set("A", fmt.Sprintf("%d", i))
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		val, err := client.GET("A")
		fmt.Println("Value received", val)
	}

	v, _ := client.GET("ahsdg")
	fmt.Printf("non existent key: %s \n", v)
}
