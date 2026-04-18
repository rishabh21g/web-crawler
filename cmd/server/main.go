package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/rishabh21g/web-crawler/internal/crawler"
	"github.com/rishabh21g/web-crawler/internal/fetcher"
	"github.com/rishabh21g/web-crawler/internal/models"
	"github.com/rishabh21g/web-crawler/internal/parser"
)

func Worker(
	id int,
	tasks chan models.URLTask,
	scheduler *crawler.Scheduler,
	wg *sync.WaitGroup,
) {
	for task := range tasks {

		log.Printf("WORKER %d: FETCHING %s (DEPTH: %d)", id, task.URL, task.Depth)

		// 1. Fetch
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		html, err := fetcher.FetchHTML(ctx, task)
		cancel()

		if err != nil {
			log.Printf("WORKER %d: ERROR FETCHING %s: %v", id, task.URL, err)
			wg.Done()
			continue
		}

		// 2. Parse
		newTasks, parseErr := parser.ParseHTML(html, task.URL, task.Depth)
		if parseErr != nil {
			log.Printf("WORKER %d: ERROR PARSING %s: %v", id, task.URL, parseErr)
			wg.Done()
			continue
		}

		// 3. Schedule new tasks
		for _, newTask := range newTasks {
			if scheduler.CheckAndMark(newTask) {
				wg.Add(1)
				tasks <- newTask // enqueue
			}
		}

		// 4. Mark current task done
		wg.Done()
	}
}
func main() {

	const WORKERS = 10
	var wg sync.WaitGroup
	tasks := make(chan models.URLTask, 100)
	startURL := "https://google.com"
	cfg := models.Config{
		Domain:   "https://google.com",
		MaxDepth: 10,
	}
	scheduler := crawler.NewScheduler(cfg)

	for i := 1; i <= WORKERS; i++ {
		wg.Add(1)
		go Worker(i, tasks, scheduler, &wg)
	}

	tasks <- models.URLTask{URL: startURL, Depth: 0}
	wg.Wait()
	close(tasks)
}
