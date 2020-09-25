package hw04_lru_cache //nolint:golint,stylecheck

import (
	"testing"

	"github.com/stretchr/testify/require"
)

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

		require.Equal(t, 7, l.Len())
		require.Equal(t, 80, l.Front().value)
		require.Equal(t, 70, l.Back().value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.value.(int))
		}
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)
	})

	t.Run("remove from ends", func(t *testing.T) {
		l := NewList()

		firstItem := l.PushBack(10)  // [10]
		middleItem := l.PushBack(20) // [10, 20]
		lastItem := l.PushBack(30)   // [10, 20, 30]

		l.Remove(firstItem) // 20, 30
		require.Equal(t, middleItem, l.Front())

		l.Remove(lastItem) // 20
		require.Equal(t, middleItem, l.Front())
		require.Equal(t, middleItem, l.Back())
	})

	t.Run("move to front from middle", func(t *testing.T) {
		l := NewList()

		firstItem := l.PushBack(10)   // [10]
		middleItem := l.PushFront(20) // [20, 10]
		l.PushFront(30)               // [30, 20, 10]

		l.MoveToFront(middleItem) // 20, 30, 10
		require.Equal(t, middleItem, l.Front())

		l.MoveToFront(firstItem) // 10, 20, 30
		require.Equal(t, firstItem, l.Front())
	})

	t.Run("operation with alien items", func(t *testing.T) {
		l1 := NewList()
		l2 := NewList()

		l1Item := l1.PushFront(1)
		l2Item := l2.PushFront(2)
		l2.Remove(l1Item)
		l2.MoveToFront(l1Item)
		require.Equal(t, l1Item, l1.Front())
		require.Equal(t, l2Item, l2.Front())
	})
}
