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

	t.Run("delete an item due to queue size", func(t *testing.T) {
		c := NewCache(3)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		wasInCache = c.Set("ccc", 300)
		require.False(t, wasInCache)

		wasInCache = c.Set("ddd", 400)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("delete an old item due to queue size", func(t *testing.T) {
		c := NewCache(3)

		wasInCache := c.Set("aaa", 100) // [(aaa, 100)]
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200) // [(bbb 200), (aaa, 100)]
		require.False(t, wasInCache)

		wasInCache = c.Set("ccc", 300) // [(ccc 300), (bbb 200), (aaa, 100)]
		require.False(t, wasInCache)

		wasInCache = c.Set("aaa", 400) // [(aaa, 400), (ccc 300), (bbb 200)]
		require.True(t, wasInCache)

		val, ok := c.Get("bbb") // [(bbb 200), (aaa, 400), (ccc 300)]
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("ddd", 600) // [(ddd 600), (bbb 200), (aaa, 400)]
		require.False(t, wasInCache)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("clear cache", func(t *testing.T) {
		c := NewCache(3)
		wasInCache := c.Set("aaa", 100) // [(aaa, 100)]
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200) // [(bbb 200), (aaa, 100)]
		require.False(t, wasInCache)

		wasInCache = c.Set("ccc", 300) // [(ccc 300), (bbb 200), (aaa, 100)]
		require.False(t, wasInCache)

		c.Clear()

		val, ok := c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)

		val, ok = c.Get("aaa")
		require.False(t, ok)
		require.Nil(t, val)

		val, ok = c.Get("bbb")
		require.False(t, ok)
		require.Nil(t, val)
	})
}

func TestCacheMultithreading(t *testing.T) {
	t.Skip() // Remove me if task with asterisk completed.

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
