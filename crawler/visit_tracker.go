package crawler

import "sync"

type VisitTracker struct {
	urls  map[string]struct{} // as go doesnt really have a builtin set like Java, Using empty struct for memory efficiency
	mutex sync.RWMutex
}

var (
	tracker *VisitTracker
	once    sync.Once
)

func NewVisitedURLTracker() *VisitTracker {
	once.Do(func() {
		tracker = &VisitTracker{
			urls: make(map[string]struct{}),
		}
	})
	return tracker
}

func (v *VisitTracker) MarkVisited(url string) {
	v.mutex.Lock()
	defer v.mutex.Unlock()
	v.urls[url] = struct{}{}
}

func (v *VisitTracker) IsVisited(url string) bool {
	v.mutex.RLock()
	defer v.mutex.RUnlock()
	_, ok := v.urls[url]
	return ok
}
