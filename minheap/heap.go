package minheap

import "container/heap"

// Element is a struct that holds an element and its frequency.
type Element struct {
	Value     string
	Frequency int
}

// MinHeap is a min-heap of values by their frequency
type MinHeap []*Element

func NewMinHeap() *MinHeap {
	h := &MinHeap{}
	heap.Init(h)
	return h
}

func (h MinHeap) Len() int {
	return len(h)
}

// Less implements comparison by Frequency
func (h MinHeap) Less(i, j int) bool {
	return h[i].Frequency < h[j].Frequency
}

func (h MinHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *MinHeap) Push(x interface{}) {
	*h = append(*h, x.(*Element))
}

// Pop the last element in the slice (min frequency)
func (h *MinHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// PopTopFrequent pops from the heap the top 'k' elements by frequency
func (h *MinHeap) PopTopFrequent(k int) []*Element {
	topKFrequent := make([]*Element, 0, k)
	for h.Len() > 0 {
		element := heap.Pop(h).(*Element)
		topKFrequent = append(topKFrequent, element)
	}
	return topKFrequent
}
