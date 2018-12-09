package application

import (
	"context"
	"github.com/duck8823/duci/application/service/job"
	"github.com/duck8823/duci/domain/model/job"
	"github.com/duck8823/duci/domain/model/job/target/git"
	"github.com/duck8823/duci/domain/model/job/target/github"
	"github.com/pkg/errors"
)

// Initialize singleton instances that are needed by application
func Initialize() error {
	if err := git.InitializeWithHTTP(func(ctx context.Context, log job.Log) {
		for line, err := log.ReadLine(); err == nil; line, err = log.ReadLine() {
			println(line.Message)
		}
	}); err != nil {
		return errors.WithStack(err)
	}

	if err := github.Initialize(Config.GitHub.APIToken.String()); err != nil {
		return errors.WithStack(err)
	}

	if err := job_service.Initialize(Config.Server.DatabasePath); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
