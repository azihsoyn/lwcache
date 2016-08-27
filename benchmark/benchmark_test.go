package cache_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/azihsoyn/lwcache"
)

func BenchmarkNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = cache.New(fmt.Sprintf("%d", i))
	}
}

func BenchmarkGet(b *testing.B) {
	c := cache.New("bench1")
	c.Set(0, 0, 10*time.Second)
	for i := 0; i < b.N; i++ {
		c.Get(i)
	}
}

func BenchmarkSet(b *testing.B) {
	c := cache.New("bench2")
	for i := 0; i < b.N; i++ {
		c.Set(i, i, 10*time.Second)
	}
}

func BenchmarkSetAndGet(b *testing.B) {
	c := cache.New("bench3")
	for i := 0; i < b.N; i++ {
		c.Set(i, i, 10*time.Second)
		c.Get(i)
	}
}
