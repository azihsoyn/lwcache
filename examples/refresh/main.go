package main

import (
	"fmt"
	"time"

	cache "github.com/azihsoyn/lwcache"
)

func main() {
	c := cache.New("refresh sample")
	key := "now"
	c.Set(key, time.Now().Format("2006-01-02 15:04:05"), 5*time.Second)
	c.SetRefresher(key, 1*time.Second, myRefresher)
	for i := 0; i < 10; i++ {
		v, ok := c.Get(key)
		fmt.Printf("current value : %v, ok : %t\n", v, ok)
		time.Sleep(1 * time.Second)
	}
}

func myRefresher(key interface{}) (interface{}, error) {
	return time.Now().Format("2006-01-02 15:04:05"), nil
}
