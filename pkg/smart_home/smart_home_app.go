package smart_home

import (
	"ameresii_smart_home/internal/err_stack"
	"ameresii_smart_home/pkg/config_manager"
	"ameresii_smart_home/pkg/cronrunner"
	"ameresii_smart_home/pkg/jobs_manager"
	"context"
	"fmt"
	"github.com/labstack/gommon/log"
)

type App struct {
	cronrunner *cronrunner.CronRunner
}

func NewSmartHomeApp(ctx context.Context) App {
	smartHome := App{
		cronrunner: cronrunner.New(ctx),
	}
	go smartHome.cronrunner.Start()
	return smartHome
}

func (a *App) Run(ctx context.Context) error {
	for _, jobConfig := range config_manager.SmartHomeConfig.Config().Jobs {
		job, err := a.produceJob(ctx, jobConfig.Method)
		if err != nil {
			return err_stack.WithStack(err)
		}

		if err = a.cronrunner.AddJob(
			jobConfig.Schedule,
			cronrunner.JobMeta{
				Name:        jobConfig.Name,
				Schedule:    jobConfig.Schedule,
				Description: jobConfig.Description,
			},
			job,
		); err != nil {
			return err_stack.WithStack(err)
		}
	}

	return nil
}

func (a *App) produceJob(_ context.Context, method string) (func(ctx context.Context), error) {
	funcDef := jobs_manager.FindFuncDefinition(method)
	if funcDef == nil {
		return nil, err_stack.WithStack(fmt.Errorf("\"method %s not found", method))
	}

	jobFunc := func(ctx context.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Info(ctx, err)
			}
		}()

		funcDef.HandlerValue.Call(nil)
	}

	return jobFunc, nil
}
