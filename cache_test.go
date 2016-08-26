package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSetAndGet(t *testing.T) {
	assert := assert.New(t)
	c := New("test")

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
	c := New("test")

	key := "test-1"
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
