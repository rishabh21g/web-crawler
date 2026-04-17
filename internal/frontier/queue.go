package frontier

import (
	"sync"

	"github.com/rishabh21g/web-crawler/internal/models"
)

type Queue struct {
	elements []models.URLTask
	mu       sync.Mutex
}

func (q *Queue) Enqueue(task models.URLTask) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.elements = append(q.elements, task)

}

func (q *Queue) Size() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.elements)
}

func (q *Queue) Dequeue() (models.URLTask, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.elements) == 0 {
		return models.URLTask{}, false
	}
	task := q.elements[0]
	q.elements = q.elements[1:]
	return task, true
}
