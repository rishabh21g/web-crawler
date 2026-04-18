package main

import (
	"sync"

	"github.com/rishabh21g/web-crawler/internal/models"
	crawler "github.com/rishabh21g/web-crawler/internal/scheduler"
	"github.com/rishabh21g/web-crawler/internal/worker"
)

func main() {
	const WORKERS = 10

	var wg sync.WaitGroup
	tasks := make(chan models.URLTask, 100)

	startURL := "https://books.toscrape.com"

	cfg := models.Config{
		Domain:   "books.toscrape.com",
		MaxDepth: 2,
	}

	scheduler := crawler.NewScheduler(cfg)

	for i := 1; i <= WORKERS; i++ {
		go worker.Worker(i, tasks, scheduler, &wg)
	}

	seed := models.URLTask{URL: startURL, Depth: 0}

	if scheduler.CheckAndMark(seed) {
		wg.Add(1)
		tasks <- seed
	}

	wg.Wait()
	close(tasks)
}
