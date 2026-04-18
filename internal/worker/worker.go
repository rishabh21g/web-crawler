package worker

import (
	"context"
	"sync"
	"time"

	"github.com/charmbracelet/log"
	"github.com/rishabh21g/web-crawler/internal/fetcher"
	"github.com/rishabh21g/web-crawler/internal/models"
	"github.com/rishabh21g/web-crawler/internal/parser"
	crawler "github.com/rishabh21g/web-crawler/internal/scheduler"
)

func Worker(
	id int,
	tasks chan models.URLTask,
	scheduler *crawler.Scheduler,
	wg *sync.WaitGroup,
) {
	for task := range tasks {

		log.Info("WORKER %d: FETCHING %s (DEPTH: %d)", id, task.URL, task.Depth)

		// 1. Fetch
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		html, err := fetcher.FetchHTML(ctx, task)
		cancel()

		if err != nil {
			log.Error("WORKER %d: ERROR FETCHING %s: %v", id, task.URL, err)
			wg.Done()
			continue
		}

		// 2. Parse
		newTasks, parseErr := parser.ParseHTML(html, task.URL, task.Depth)
		if parseErr != nil {
			log.Error("WORKER %d: ERROR PARSING %s: %v", id, task.URL, parseErr)
			wg.Done()
			continue
		}

		// 3. Schedule new tasks
		for _, newTask := range newTasks {
			if scheduler.CheckAndMark(newTask) {
				wg.Add(1)
				select {
				case tasks <- newTask:
				default:
					log.Debugf("WORKER %d : CHANNEL FULL DROPPING %s", id, task.URL)
					wg.Done()
				}
			}
		}

		// 4. Mark current task done
		log.Infof("CURRENT LENGTH OF QUEUE ~ %d: ", len(tasks))
		wg.Done()
	}
}
