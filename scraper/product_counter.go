package scraper

import "sync"

type ProductCounter struct {
	mutex sync.Mutex
	count int
}

var (
	counter     *ProductCounter
	counterOnce sync.Once
)

func NewProductCounter() *ProductCounter {
	counterOnce.Do(func() {
		counter = &ProductCounter{
			count: 0,
		}
	})
	return counter
}

func (pc *ProductCounter) Increment() {
	pc.mutex.Lock()
	defer pc.mutex.Unlock()
	pc.count++
}

func (pc *ProductCounter) GetCount() int {
	pc.mutex.Lock()
	defer pc.mutex.Unlock()
	return pc.count
}

func (pc *ProductCounter) Reset() {
	pc.mutex.Lock()
	defer pc.mutex.Unlock()
	pc.count = 0
}
