package queue_test

import (
	"github.com/Somefive/xcontainer/queue"
	"math/rand"
	"sync"
	"testing"
)

func TestPriorityQueue(t *testing.T) {
	n, m := 10000, 10
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = rand.Int()
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
				pq.Push(rand.Int())
				pq.Pop()
			}
		}()
	}
	wg.Wait()
	top := pq.Pop()
	for i := 0; i < n-1; i++ {
		cur := pq.Pop()
		if cur > top {
			t.Errorf("invalid order: %d(index:%d) > %d(index:%d)", cur, i+1, top, i)
		}
	}
}
