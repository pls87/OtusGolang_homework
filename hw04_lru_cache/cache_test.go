package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

type Expected struct {
	inCache bool
	ok      bool
	len     int
	queue   []interface{}
}

type Step struct {
	action   string
	key      Key
	value    interface{}
	expected Expected
}

var emptyTests = []Step{
	{
		action: "get", key: "aaa", value: nil,
		expected: Expected{inCache: false, ok: false, len: 0, queue: []interface{}{}},
	}, {
		action: "get", key: "bbb", value: nil,
		expected: Expected{inCache: false, ok: false, len: 0, queue: []interface{}{}},
	},
}

var simpleTests = []Step{
	{
		action: "set", key: "aaa", value: 100,
		expected: Expected{inCache: false, ok: true, len: 1, queue: []interface{}{
			&cacheItem{"aaa", 100},
		}},
	}, {
		action: "set", key: "bbb", value: 200,
		expected: Expected{inCache: false, ok: true, len: 2, queue: []interface{}{
			&cacheItem{"bbb", 200},
			&cacheItem{"aaa", 100},
		}},
	}, {
		action: "get", key: "aaa", value: 100,
		expected: Expected{inCache: true, ok: true, len: 2, queue: []interface{}{
			&cacheItem{"aaa", 100},
			&cacheItem{"bbb", 200},
		}},
	}, {
		action: "get", key: "bbb", value: 200,
		expected: Expected{inCache: true, ok: true, len: 2, queue: []interface{}{
			&cacheItem{"bbb", 200},
			&cacheItem{"aaa", 100},
		}},
	}, {
		action: "set", key: "aaa", value: 300,
		expected: Expected{inCache: true, ok: true, len: 2, queue: []interface{}{
			&cacheItem{"aaa", 300},
			&cacheItem{"bbb", 200},
		}},
	}, {
		action: "get", key: "aaa", value: 300,
		expected: Expected{inCache: true, ok: true, len: 2, queue: []interface{}{
			&cacheItem{"aaa", 300},
			&cacheItem{"bbb", 200},
		}},
	}, {
		action: "get", key: "ccc", value: nil,
		expected: Expected{inCache: false, ok: false, len: 2, queue: []interface{}{
			&cacheItem{"aaa", 300},
			&cacheItem{"bbb", 200},
		}},
	},
}

var purgeTests = []Step{
	{
		action: "set", key: "aaa", value: 100,
		expected: Expected{inCache: false, ok: true, len: 1, queue: []interface{}{
			&cacheItem{"aaa", 100},
		}},
	}, {
		action: "set", key: "bbb", value: 200,
		expected: Expected{inCache: false, ok: true, len: 2, queue: []interface{}{
			&cacheItem{"bbb", 200},
			&cacheItem{"aaa", 100},
		}},
	}, {
		action: "set", key: "ccc", value: 300,
		expected: Expected{inCache: false, ok: true, len: 3, queue: []interface{}{
			&cacheItem{"ccc", 300},
			&cacheItem{"bbb", 200},
			&cacheItem{"aaa", 100},
		}},
	}, {
		action: "set", key: "eee", value: 800,
		expected: Expected{inCache: false, ok: true, len: 4, queue: []interface{}{
			&cacheItem{"eee", 800},
			&cacheItem{"ccc", 300},
			&cacheItem{"bbb", 200},
			&cacheItem{"aaa", 100},
		}},
	}, {
		action: "get", key: "bbb", value: 200,
		expected: Expected{inCache: true, ok: true, len: 4, queue: []interface{}{
			&cacheItem{"bbb", 200},
			&cacheItem{"eee", 800},
			&cacheItem{"ccc", 300},
			&cacheItem{"aaa", 100},
		}},
	}, {
		action: "set", key: "ddd", value: 400,
		expected: Expected{inCache: false, ok: true, len: 4, queue: []interface{}{
			&cacheItem{"ddd", 400},
			&cacheItem{"bbb", 200},
			&cacheItem{"eee", 800},
			&cacheItem{"ccc", 300},
		}},
	}, {
		action: "get", key: "aaa", value: nil,
		expected: Expected{inCache: false, ok: false, len: 4, queue: []interface{}{
			&cacheItem{"ddd", 400},
			&cacheItem{"bbb", 200},
			&cacheItem{"eee", 800},
			&cacheItem{"ccc", 300},
		}},
	}, {
		action: "set", key: "ccc", value: 500,
		expected: Expected{inCache: true, ok: true, len: 4, queue: []interface{}{
			&cacheItem{"ccc", 500},
			&cacheItem{"ddd", 400},
			&cacheItem{"bbb", 200},
			&cacheItem{"eee", 800},
		}},
	}, {
		action: "get", key: "bbb", value: 200,
		expected: Expected{inCache: true, ok: true, len: 4, queue: []interface{}{
			&cacheItem{"bbb", 200},
			&cacheItem{"ccc", 500},
			&cacheItem{"ddd", 400},
			&cacheItem{"eee", 800},
		}},
	}, {
		action: "set", key: "fff", value: 200,
		expected: Expected{inCache: false, ok: true, len: 4, queue: []interface{}{
			&cacheItem{"fff", 200},
			&cacheItem{"bbb", 200},
			&cacheItem{"ccc", 500},
			&cacheItem{"ddd", 400},
		}},
	},
}

func runSteps(t *testing.T, c Cache, actions []Step) {
	t.Helper()
	for _, tc := range actions {
		switch tc.action {
		case "set":
			require.Equal(t, tc.expected.inCache, c.Set(tc.key, tc.value))
		case "get":
			val, ok := c.Get(tc.key)
			require.Equal(t, tc.expected.ok, ok)
			require.Equal(t, tc.value, val)
		}
		require.Equal(t, tc.expected.len, c.(*lruCache).queue.Len())
		checkListEQ(t, tc.expected.queue, c.(*lruCache).queue)
	}
}

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		runSteps(t, NewCache(10), emptyTests)
	})

	t.Run("simple", func(t *testing.T) {
		runSteps(t, NewCache(5), simpleTests)
	})

	t.Run("purge logic", func(t *testing.T) {
		runSteps(t, NewCache(4), purgeTests)
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
