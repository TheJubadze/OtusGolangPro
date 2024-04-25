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
		c := NewCache(5)

		c.Set("aaa", 100)
		require.Equal(t, "[{aaa 100}]", c.String())

		c.Set("bbb", 200)
		require.Equal(t, "[{bbb 200}, {aaa 100}]", c.String())

		c.Set("aaa", 300)
		require.Equal(t, "[{aaa 300}, {bbb 200}]", c.String())

		c.Set("ccc", 400)
		require.Equal(t, "[{ccc 400}, {aaa 300}, {bbb 200}]", c.String())

		c.Set("ddd", 500)
		require.Equal(t, "[{ddd 500}, {ccc 400}, {aaa 300}, {bbb 200}]", c.String())

		c.Set("eee", 600)
		require.Equal(t, "[{eee 600}, {ddd 500}, {ccc 400}, {aaa 300}, {bbb 200}]", c.String())

		c.Set("fff", 700)
		require.Equal(t, "[{fff 700}, {eee 600}, {ddd 500}, {ccc 400}, {aaa 300}]", c.String())

		c.Set("ggg", 800)
		require.Equal(t, "[{ggg 800}, {fff 700}, {eee 600}, {ddd 500}, {ccc 400}]", c.String())

		c.Set("hhh", 900)
		require.Equal(t, "[{hhh 900}, {ggg 800}, {fff 700}, {eee 600}, {ddd 500}]", c.String())

		val, ok := c.Get("aaa")
		require.False(t, ok)
		require.Nil(t, val)

		c.Clear()
		val, ok = c.Get("hhh")
		require.False(t, ok)
		require.Nil(t, val)
		require.Equal(t, "[]", c.String())
	})

	t.Run("lru removing", func(t *testing.T) {
		c := NewCache(3)

		c.Set("aaa", 100)
		c.Set("bbb", 200)
		c.Set("aaa", 300)
		c.Set("ccc", 400)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		c.Set("ddd", 500)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})
}

func TestCacheMultithreading(t *testing.T) {
	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(3)

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

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Clear()
		}
	}()

	wg.Wait()
}
