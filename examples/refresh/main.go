package main

import (
	"log"
	"time"

	"github.com/azihsoyn/lwcache"
)

func main() {
	c := lwcache.New("refresh sample")
	key := "now"
	c.Set(key, time.Now().Format("2006-01-02 15:04:05"), 5*time.Second)
	c.SetRefresher(myRefresher)
	c.StartRefresher(key, 1*time.Second)
	c.StartRefresher(key, 1*time.Second) // this is no effect
	for i := 0; i < 10; i++ {
		v, ok := c.Get(key)
		log.Printf("current value(before stop) : %v, ok : %t\n", v, ok)
		time.Sleep(1 * time.Second)
	}

	c.StopRefresher(key)
	for i := 0; i < 10; i++ {
		v, ok := c.Get(key)
		log.Printf("current value(after stop) : %v, ok : %t\n", v, ok)
		time.Sleep(1 * time.Second)
	}
}

// heavy process
func myRefresher(c lwcache.Cache, key interface{}, currentValue interface{}) (interface{}, error) {
	log.Println("called refresh", key)
	time.Sleep(3 * time.Second)
	c.SetExpire(key, 5*time.Second) // on refresh extend expiration
	return time.Now().Format("2006-01-02 15:04:05"), nil
}
