package cron

type Job struct {
}

type Cron interface {
	AddJob(t *Job) (string, error)
	DelJob(id string) error
	StopJob(id string) error
	Start() error
}
