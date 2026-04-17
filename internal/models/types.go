package models

import "time"

type URLTask struct {
	URL   string
	Depth int
}

type Page struct {
	URL       string
	Status    string
	Links     []string
	FetchedAt time.Time
}

type Config struct {
	MaxPages int
	MaxDepth int
	Domain   string
}
