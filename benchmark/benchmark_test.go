package lwcache_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/azihsoyn/lwcache"
)

func BenchmarkNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = lwcache.New(fmt.Sprintf("%d", i))
	}
}

func BenchmarkGet(b *testing.B) {
	c := lwcache.New("bench1")
	c.Set(0, 0, 1*time.Microsecond)
	for i := 0; i < b.N; i++ {
		c.Get(i)
	}
}

func BenchmarkGet_SimpleMap(b *testing.B) {
	m := make(map[interface{}]interface{})
	m[0] = 0
	for i := 0; i < b.N; i++ {
		_ = m[i]
	}
}

func BenchmarkSet(b *testing.B) {
	c := lwcache.New("bench2")
	for i := 0; i < b.N; i++ {
		c.Set(i, i, 10*time.Second)
	}
}

func BenchmarkSet_ConcurrentAccess(b *testing.B) {
	c := lwcache.New("bench3")
	ch := make(chan int, 8192)
	for i := 0; i < b.N; i++ {
		select {
		case ch <- i:
			go func(i int) {
				c.Set(i, i, 10*time.Second)
			}(i)
		default:
		}
	}
}

func BenchmarkStartRefresher_ConcurrentAccess(b *testing.B) {
	c := lwcache.New("bench4")
	c.SetRefresher(func(c lwcache.Cache, key, currentValue interface{}) (interface{}, error) {
		return currentValue, nil
	})
	ch := make(chan int, 8192)
	for i := 0; i < b.N; i++ {
		select {
		case ch <- i:
			go func(i int) {
				c.Set(i, i, 10*time.Second)
				c.StartRefresher(i, 1*time.Millisecond)
			}(i)
		default:
		}
	}
}

/* fatal error: concurrent map writes
func BenchmarkSet_WithSimpleMap(b *testing.B) {
	m := make(map[int]int)
	ch := make(chan int, 8192)
	for i := 0; i < b.N; i++ {
		select {
		case ch <- i:
			go func(i int) {
				m[i] = i
			}(i)
		default:
		}
	}
}
*/

/* fatal error: concurrent map writes
func BenchmarkSet_SimpleMap(b *testing.B) {
	m := make(map[int]int)
	for i := 0; i < b.N; i++ {
		go func(i int) {
			m[i] = i
		}(i)
	}
}
*/

func BenchmarkSetAndGet(b *testing.B) {
	c := lwcache.New("bench5")
	for i := 0; i < b.N; i++ {
		c.Set(i, i, 10*time.Second)
		c.Get(i)
	}
}

func BenchmarkSetAndGet_WithSimpleMap(b *testing.B) {
	m := make(map[interface{}]interface{})
	for i := 0; i < b.N; i++ {
		m[i] = i
		_ = m[i]
	}
}
