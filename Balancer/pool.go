package Balancer

// Heap
type Pool []*Worker

// Implementing Heap interface

// Less - mean worker with lowest jobs count
func (p Pool) Less(i, j int) bool { return p[i].pending < p[j].pending }

// Workers count in pool
func (p Pool) Len() int { return len(p) }

// Swap
func (p Pool) Swap(i, j int) {
	if i >= 0 && i < len(p) && j >= 0 && j < len(p) {
		p[i], p[j] = p[j], p[i]
		p[i].index, p[j].index = i, j
	}
}

// Push Job
func (p *Pool) Push(x interface{}) {
	n := len(*p)
	worker := x.(*Worker)
	worker.index = n
	*p = append(*p, worker)
}

// Pop Job
func (p *Pool) Pop() interface{} {
	old := *p
	n := len(old)
	item := old[n-1]
	item.index = -1
	*p = old[0 : n-1]
	return item
}
