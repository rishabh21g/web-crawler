package crawler

import (
	"net/url"
	"sync"

	"github.com/rishabh21g/web-crawler/internal/models"
)

// custom data type to check it is valid for scheduling or not
type Scheduler struct {
	Seen map[string]bool
	cfg  models.Config
	mu   sync.Mutex
}

// constructor function just to make object of scheduler struct
func NewScheduler(cfg models.Config) *Scheduler {
	return &Scheduler{
		Seen: make(map[string]bool),
		cfg:  cfg,
	}
}

func (l *Scheduler) CheckAndMark(task models.URLTask) bool {
	parsedURL, err := url.Parse(task.URL)

	// check we limit our depth or not
	if task.Depth > l.cfg.MaxDepth {
		return false
	}
	// parsing the raw url into the go understanble url
	if err != nil {
		return false
	}

	// checking the host already in our domain or not
	if parsedURL.Host != l.cfg.Domain {
		return false
	}

	//Lock for thread safety
	l.mu.Lock()
	defer l.mu.Unlock()

	if _, ok := l.Seen[task.URL]; ok {
		return false

	}
	// all set then mark it seen and return the true
	l.Seen[task.URL] = true
	return true

}
