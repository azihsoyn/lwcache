package main

import (
	"fmt"
	"time"

	"github.com/azihsoyn/cache"
)

func main() {
	c := cache.New("sample")
	fmt.Println("now : ", time.Now())
	c.Set(1, "Apple", 1*time.Second)
	c.Set(2, "Banana", 2*time.Second)
	c.Set(3, "Cake", 3*time.Second)
	for j := 0; j < 4; j++ {
		for i := 1; i <= 3; i++ {
			v, ok := c.Get(i)
			fmt.Printf("value : %v, ok : %t, now : %v\n", v, ok, time.Now())
		}
		time.Sleep(1 * time.Second)
	}
}
