package main

import (
	"fmt"
	"time"

	cache "github.com/azihsoyn/lwcache"
)

func main() {
	c := cache.New("sample")
	fmt.Println("now : ", time.Now().Format("2006-01-02 15:04:05"))
	c.Set(1, "Apple", 1*time.Second)
	c.Set(2, "Banana", 2*time.Second)
	c.Set(3, "Cake", 3*time.Second)
	c.SetExpire(1, 5*time.Second)
	for j := 0; j < 4; j++ {
		for i := 1; i <= 3; i++ {
			v, ok := c.Get(i)
			fmt.Printf("value : %6v, ok : %5t, now : %v\n", v, ok, time.Now().Format("2006-01-02 15:04:05"))
		}
		time.Sleep(1 * time.Second)
	}
}
