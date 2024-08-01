package queue_test

import (
	"github.com/Somefive/xcontainer/queue"
	"github.com/stretchr/testify/require"
	"math/rand"
	"sort"
	"sync"
	"testing"
)

func TestPriorityQueue(t *testing.T) {
	n, m := 10000, 10
	arr := make([]int, n)
	rand.Seed(42)
	for i := 0; i < n; i++ {
		arr[i] = rand.Int() % n
	}
	pq := queue.NewPriorityQueue(arr, func(i int, j int) bool {
		return i < j
	})
	wg := sync.WaitGroup{}
	for i := 0; i < m; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < n/m; j++ {
				pq.Push(rand.Int() % n)
				pq.Pop()
			}
		}()
	}
	wg.Wait()
	expected := pq.Items()
	for i := 0; i < len(expected)>>1; i++ {
		if i*2+1 < len(expected) {
			require.LessOrEqual(t, expected[i], expected[i*2+1])
		}
		if i*2+2 < len(expected) {
			require.LessOrEqual(t, expected[i], expected[i*2+2])
		}
	}
	sort.Slice(expected, func(i, j int) bool { return expected[i] < expected[j] })
	top := pq.Pop()
	require.Equal(t, expected[0], top)
	actual := []int{top}
	for i := 0; i < n-1; i++ {
		cur := pq.Pop()
		actual = append(actual, cur)
		require.LessOrEqual(t, top, cur, "\n%v\n%v", expected, actual)
		top = cur
	}
	require.Equal(t, expected, actual, "\n%v\n%v", expected, actual)
}
