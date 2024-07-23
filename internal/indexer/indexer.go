package indexer

import (
	"log"

	"ccrayz/event-scanner/internal/indexer/task"

	"github.com/robfig/cron"
)

type Indexer struct {
	cron     *cron.Cron
	schedule string
	tasks    []Task
}

func NewIndexer(schedule string) *Indexer {
	tasks := []Task{SampleTask{}, task.GetNodeInfo{}}
	c := cron.New()
	return &Indexer{
		cron:     c,
		schedule: schedule,
		tasks:    tasks,
	}
}

func (i *Indexer) Run() {
	log.Printf("Running indexer with schedule %s", i.schedule)

	for _, task := range i.tasks {
		err := i.cron.AddFunc(i.schedule, task.Do)
		if err != nil {
			log.Fatalf("Failed to add task to cron: %v", err)
		}
	}

	i.cron.Start()
}

func (i *Indexer) Stop() {
	i.cron.Stop()
}

type Task interface {
	Do()
}

type SampleTask struct{}

func (t SampleTask) Do() {
	log.Println("Running SampleTask")
}
