package indexer

import (
	"log"

	"ccrayz/event-scanner/internal/indexer/task"

	"github.com/robfig/cron"
	"gorm.io/gorm"
)

type Indexer struct {
	cron     *cron.Cron
	schedule string
	tasks    []Task
}

func NewIndexer(schedule string, db *gorm.DB) *Indexer {
	log.Println(db)
	tasks := []Task{SampleTask{}, task.GetNodeInfo{}}
	c := cron.New()
	indexer := &Indexer{
		cron:     c,
		schedule: schedule,
		tasks:    tasks,
	}

	for _, task := range indexer.tasks {
		task.SetDB(db)
	}

	return indexer
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
	SetDB(db *gorm.DB)
}

type SampleTask struct {
	db *gorm.DB
}

func (t SampleTask) Do() {
	log.Println("Running SampleTask")
}

func (t SampleTask) SetDB(db *gorm.DB) {
	t.db = db
}
