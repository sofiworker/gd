package cron

import (
	"github.com/go-co-op/gocron/v2"
	"time"
)

type GoCronClient struct {
	scheduler gocron.Scheduler
}

func NewGoCron() (Cron, error) {
	s, err := gocron.NewScheduler()
	if err != nil {
		return nil, err
	}

	return &GoCronClient{
		scheduler: s,
	}, nil
}

func (g *GoCronClient) AddJob(t *Job) (string, error) {
	job, err := g.scheduler.NewJob(gocron.DurationJob(
		10*time.Second,
	), gocron.NewTask(
		func(a string, b int) {
			// do things
		},
		"hello",
		1,
	))
	if err != nil {
		return "", err
	}
	return job.ID().String(), nil
}

func (g *GoCronClient) DelJob(id string) error {
	//TODO implement me
	panic("implement me")
}

func (g *GoCronClient) StopJob(id string) error {
	//TODO implement me
	panic("implement me")
}

func (g *GoCronClient) Start() error {
	g.scheduler.Start()
	return nil
}
