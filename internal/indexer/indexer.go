package indexer

import (
	"ccrayz/event-scanner/internal/db"
	"ccrayz/event-scanner/internal/indexer/tasks"
	"log"

	"github.com/robfig/cron"
)

type Indexer struct {
	cron     *cron.Cron
	schedule string
	tasks    []Task
}

type Task interface {
	Do(db *db.AppDB)
}

func NewIndexer(schedule string) *Indexer {
	tasks := []Task{tasks.GetNodeInfo{}}
	c := cron.New()
	indexer := &Indexer{
		cron:     c,
		schedule: schedule,
		tasks:    tasks,
	}

	return indexer
}

func (i *Indexer) Run(db *db.AppDB) {
	log.Printf("Running indexer with schedule %s", i.schedule)
	for _, task := range i.tasks {

		err := i.cron.AddFunc(i.schedule, func() {
			task.Do(db)
		})
		if err != nil {
			log.Fatalf("Failed to add task to cron: %v", err)
		}
	}

	i.cron.Start()
}

func (i *Indexer) Stop() {
	i.cron.Stop()
}

func (i *Indexer) AddTask(task Task) {
	i.tasks = append(i.tasks, task)
}
