package main

import (
	"fmt"

	"gopkg.in/azihsoyn/lwcache.v1"
)

type cacheKey struct {
	title string
	page  int
}

const numPerPage = 10

func main() {
	c := lwcache.New("flexible-key")
	titles := []string{
		"title1",
		"title2",
		"title3",
	}
	for _, title := range titles {
		item := make([]string, 0, numPerPage)
		for i := 0; i < 100; i++ {
			item = append(item, fmt.Sprintf("%s-%d", title, i))
			if len(item) == numPerPage {
				key := cacheKey{title, i / numPerPage}
				c.Set(key, item, lwcache.NoExpire)
				item = make([]string, 0, numPerPage)
			}
		}
	}
	for _, title := range titles {
		for i := 0; i < 10; i++ {
			key := cacheKey{title, i}
			v, ok := c.Get(key)
			fmt.Printf("key : %v, value : %3v, ok : %5t\n", key, v, ok)
		}
	}
}
