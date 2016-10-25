package lwcache

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSetAndGet(t *testing.T) {
	assert := assert.New(t)
	c := New("test-1")

	key := "test-1"
	expect := "test-1-value"

	c.Set(key, expect, 500*time.Millisecond)
	actual, ok := c.Get(key)
	assert.True(ok)
	assert.Equal(expect, actual)

	// after expire
	time.Sleep(550 * time.Millisecond)
	actual, ok = c.Get(key)
	assert.False(ok)
	assert.Equal(nil, actual)
}

func TestSetExpire(t *testing.T) {
	assert := assert.New(t)
	c := New("test-2")

	key := "test-2"
	expect := "test-1-value"

	c.Set(key, expect, 500*time.Millisecond)
	actual, ok := c.Get(key)
	assert.True(ok)
	assert.Equal(expect, actual)

	c.SetExpire(key, 1*time.Second)
	// after first expiration
	time.Sleep(550 * time.Millisecond)
	actual, ok = c.Get(key)
	assert.True(ok)
	assert.Equal(expect, actual)
}

func TestSetRefresher(t *testing.T) {
	assert := assert.New(t)
	c := New("test-4")

	key := "test-4"
	expect := 0
	refresher := func(c Cache, key interface{}, currentValue interface{}) (interface{}, error) {
		num, ok := currentValue.(int)
		if ok {
			return num + 1, nil
		}
		return 0, errors.New("refresh failed")
	}

	c.Set(key, expect, 10*time.Second)
	c.SetRefresher(refresher)
	c.StartRefresher(key, 1*time.Second)

	for i := 0; i < 5; i++ {
		actual, ok := c.Get(key)
		assert.True(ok)
		assert.Equal(expect, actual, fmt.Sprintf("Test No. %d", i+1))
		time.Sleep(1050 * time.Millisecond)
		expect++
	}

	c.StopRefresher(key)
	for i := 0; i < 2; i++ {
		time.Sleep(1050 * time.Millisecond)
		actual, ok := c.Get(key)
		assert.True(ok)
		assert.Equal(expect, actual) // expect not changed
	}
}

func TestSetRefresher_OnExpired(t *testing.T) {
	assert := assert.New(t)
	c := New("test-5")

	key := "test-5"
	expect := 0
	refresher := func(c Cache, key interface{}, currentValue interface{}) (interface{}, error) {
		num, ok := currentValue.(int)
		if ok {
			return num + 1, nil
		}
		return 0, errors.New("refresh failed")
	}

	c.Set(key, expect, 2*time.Second)
	c.SetRefresher(refresher)
	c.StartRefresher(key, 1*time.Second)

	actual, ok := c.Get(key)
	assert.True(ok)
	assert.Equal(expect, actual)
	time.Sleep(1500 * time.Millisecond)

	// refresh
	actual, ok = c.Get(key)
	assert.True(ok)
	assert.Equal(expect+1, actual)
	time.Sleep(1500 * time.Millisecond)

	// expire
	actual, ok = c.Get(key)
	assert.False(ok)
	assert.Equal(nil, actual)
}

func TestSetRefresher_OnRefreshError(t *testing.T) {
	assert := assert.New(t)
	c := New("test-3")

	key := "test-3"
	expect := 0
	refresher := func(c Cache, key interface{}, currentValue interface{}) (interface{}, error) {
		return 0, errors.New("refresh failed")
	}

	c.Set(key, expect, 10*time.Second)
	c.SetRefresher(refresher)
	c.StartRefresher(key, 1*time.Second)

	for i := 0; i < 5; i++ {
		actual, ok := c.Get(key)
		assert.True(ok)
		assert.Equal(expect, actual, fmt.Sprintf("Test No. %d", i+1))
		time.Sleep(1050 * time.Millisecond)
	}
}

func TestStartRefresher_OnNoRefresher(t *testing.T) {
	c := New("test-6")

	key := "test-6"
	expect := 0

	c.Set(key, expect, 10*time.Second)
	c.StartRefresher(key, 500*time.Millisecond)

	time.Sleep(1 * time.Second)
	// check no panic
}
