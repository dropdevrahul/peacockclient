package peacockclient_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/dropdevrahul/peacockclient"
)

func TestMultipleConnection(t *testing.T) {
	var wg sync.WaitGroup

	for i := 0; i < 8; i++ {
		wg.Add(1)
		c := peacockclient.Client{
			Host: "127.0.0.1",
			Port: "9999",
		}
		go connectAndGo(&c, fmt.Sprintf("%d", i+1), fmt.Sprintf("%d", i), &wg)
	}

	wg.Wait()
}

func connectAndGo(c *peacockclient.Client,
	k string, v string, wg *sync.WaitGroup) {
	defer wg.Done()
	r, _ := c.Set(k, v)
	fmt.Println("SET", k, " =>", v, "server =>", string(r.Data))
	r, _ = c.Get(k)
	fmt.Println("GET Key =>", k, " server =>", string(r.Data))
}
