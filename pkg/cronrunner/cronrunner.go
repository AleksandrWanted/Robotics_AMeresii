package cronrunner

import (
	"context"
	"github.com/AleksandrWanted/AMeresii_SMART_HOME/internal/err_stack"
	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
	"time"
)

type CronRunner struct {
	ctx     context.Context
	c       *cron.Cron
	jobMap  map[string]cron.EntryID
	jobInfo map[string]JobMeta
}

type JobMeta struct {
	ID          string
	Name        string `yaml:"name"`
	Schedule    string `yaml:"schedule"`
	Description string `yaml:"description"`
	Started     time.Time
}

func New(ctx context.Context) *CronRunner {
	return &CronRunner{
		ctx: ctx,
		c: cron.New(
			cron.WithLocation(time.UTC),
			cron.WithSeconds(),
		),
		jobMap:  make(map[string]cron.EntryID, 0),
		jobInfo: make(map[string]JobMeta, 0),
	}
}

func (cr *CronRunner) AddJob(
	schedule string,
	job JobMeta,
	fn func(ctx context.Context),
) error {
	jobID, err := cr.c.AddFunc(schedule, func() {
		fn(cr.ctx)
	})
	job.ID = uuid.New().String()
	job.Started = time.Now()
	if err != nil {
		return err_stack.WithStack(err)
	}
	cr.jobMap[job.ID] = jobID
	cr.jobInfo[job.ID] = job

	return nil
}

func (cr *CronRunner) DeleteJobByID(id uuid.UUID) {
	if entryID, ok := cr.jobMap[id.String()]; ok {
		cr.c.Remove(entryID)
	}
	delete(cr.jobInfo, id.String())
}

func (cr *CronRunner) DeleteAll() {
	for k := range cr.jobMap {
		if entryID, ok := cr.jobMap[k]; ok {
			cr.c.Remove(entryID)
		}
	}
	for k := range cr.jobInfo {
		delete(cr.jobInfo, k)
	}
}

func (cr *CronRunner) List() map[string]JobMeta {
	return cr.jobInfo
}

func (cr *CronRunner) Start() {
	cr.c.Start()
}

func (cr *CronRunner) Stop() {
	cr.c.Stop()
}
