package main

import (
	"fmt"
	"time"

	"github.com/azihsoyn/lwcache"
)

func main() {
	c := lwcache.New("no-expire")
	key := "key"
	c.Set(key, time.Now().Format("2006-01-02 15:04:05"), lwcache.NoExpire)
	for i := 0; i < 10; i++ {
		v, ok := c.Get(key)
		fmt.Printf("current value : %v, ok : %t\n", v, ok)
		time.Sleep(1 * time.Second)
	}
}
