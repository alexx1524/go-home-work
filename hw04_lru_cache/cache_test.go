package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic", func(t *testing.T) {
		c := NewCache(3)
		c.Set("1", 1)
		c.Set("2", 2)
		c.Set("3", 3)

		c.Set("4", 4)
		value, ok := c.Get("1")
		require.False(t, ok)
		require.Nil(t, value)
	})

	t.Run("purge the oldest item", func(t *testing.T) {
		c := NewCache(3)
		c.Set("1", 1)
		c.Set("2", 2)
		c.Set("3", 3) // [3, 2, 1]

		c.Set("1", 100)
		c.Set("2", 200)
		c.Get("3")
		c.Get("1") // [1, 3, 2]

		c.Set("4", 4) // [4, 1, 3]
		value, ok := c.Get("2")
		require.False(t, ok)
		require.Nil(t, value)
	})

	t.Run("clear cache", func(t *testing.T) {
		capacity := 5
		c := NewCache(capacity)

		for i := 1; i <= capacity; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}

		for i := 1; i <= capacity; i++ {
			value, ok := c.Get(Key(strconv.Itoa(i)))
			require.True(t, ok)
			require.Equal(t, i, value)
		}

		c.Clear()

		for i := 1; i <= capacity; i++ {
			value, ok := c.Get(Key(strconv.Itoa(i)))
			require.False(t, ok)
			require.Nil(t, value)
		}
	})
}

func TestCacheMultithreading(t *testing.T) {
	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}
