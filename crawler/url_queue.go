package crawler

import "sync"

type URLQueue struct {
	queue map[string]struct{}
	mutex sync.Mutex
}

var (
	queue     *URLQueue
	queueOnce sync.Once
)

func NewURLQueue() *URLQueue {
	queueOnce.Do(func() {
		queue = &URLQueue{
			queue: make(map[string]struct{}),
		}
	})
	return queue
}

func (u *URLQueue) Add(url string) {
	u.mutex.Lock()
	defer u.mutex.Unlock()
	u.queue[url] = struct{}{}
}

func (u *URLQueue) Remove(url string) {
	u.mutex.Lock()
	defer u.mutex.Unlock()
	delete(u.queue, url)
}

func (u *URLQueue) IsEmpty() bool {
	u.mutex.Lock()
	defer u.mutex.Unlock()
	return len(u.queue) == 0
}

func (u *URLQueue) GetLength() int {
	u.mutex.Lock()
	defer u.mutex.Unlock()
	return len(u.queue)
}

func (u *URLQueue) GetQueue() map[string]struct{} {
	u.mutex.Lock()
	defer u.mutex.Unlock()
	return u.queue
}

func (u *URLQueue) Reset() {
	u.mutex.Lock()
	defer u.mutex.Unlock()
	u.queue = make(map[string]struct{})
}
