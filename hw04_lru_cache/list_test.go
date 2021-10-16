package hw04lrucache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func checkListEQ(t *testing.T, expected []interface{}, actual List) {
	t.Helper()

	elems := make([]interface{}, 0, actual.Len())
	for i := actual.Front(); i != nil; i = i.Next {
		elems = append(elems, i.Value)
	}
	require.Equal(t, len(expected), len(elems))
	require.Equal(t, expected, elems)

	elems = make([]interface{}, 0, actual.Len())
	for i := actual.Back(); i != nil; i = i.Prev {
		elems = append([]interface{}{i.Value}, elems...)
	}

	require.Equal(t, len(expected), len(elems))
	require.Equal(t, expected, elems)
}

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())

		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
		require.Equal(t, 2, l.Len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		checkListEQ(t, []interface{}{80, 60, 40, 10, 30, 50, 70}, l)

		require.Equal(t, 7, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]
		l.MoveToFront(l.Back())  // [50, 70, 80, 60, 40, 10, 30]

		checkListEQ(t, []interface{}{50, 70, 80, 60, 40, 10, 30}, l)

		l.MoveToFront(l.Front().Next.Next.Next) // [60, 50, 70, 80, 40, 10, 30]

		checkListEQ(t, []interface{}{60, 50, 70, 80, 40, 10, 30}, l)

		l.Remove(l.Front().Next.Next)
		checkListEQ(t, []interface{}{60, 50, 80, 40, 10, 30}, l)

		for l.Len() != 0 {
			if l.Len()%2 == 0 {
				l.Remove(l.Front())
			} else {
				l.Remove(l.Back())
			}
		}

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})
}
