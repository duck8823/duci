package application

import (
	"context"
	jobService "github.com/duck8823/duci/application/service/job"
	"github.com/duck8823/duci/domain/model/job"
	"github.com/duck8823/duci/domain/model/job/target/git"
	"github.com/duck8823/duci/domain/model/job/target/github"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Initialize singleton instances that are needed by application
func Initialize() error {
	switch {
	case len(Config.GitHub.SSHKeyPath) == 0:
		if err := git.InitializeWithHTTP(printLog); err != nil {
			return errors.WithStack(err)
		}
	default:
		if err := git.InitializeWithSSH(Config.GitHub.SSHKeyPath, printLog); err != nil {
			return errors.WithStack(err)
		}
	}

	if err := github.Initialize(Config.GitHub.APIToken.String()); err != nil {
		return errors.WithStack(err)
	}

	if err := jobService.Initialize(Config.Server.DatabasePath); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func printLog(_ context.Context, log job.Log) {
	for line, err := log.ReadLine(); err == nil; line, err = log.ReadLine() {
		logrus.Info(line.Message)
	}
}
