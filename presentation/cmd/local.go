package cmd

import (
	"context"
	"github.com/duck8823/duci/domain/model/docker"
	"github.com/duck8823/duci/domain/model/job"
	"github.com/duck8823/duci/domain/model/runner"
	"github.com/labstack/gommon/random"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var localCmd = createCmd("local", "Run locally", runLocal)

func runLocal(cmd *cobra.Command, args []string) {
	readConfiguration(cmd)

	runner := runner.DefaultDockerRunnerBuilder().
		LogFunc(func(_ context.Context, log job.Log) {
			for line, err := log.ReadLine(); err == nil; line, err = log.ReadLine() {
				logrus.Info(line.Message)
			}
		}).
		Build()

	if err := runner.Run(context.Background(), ".", docker.Tag(random.String(12, random.Hex)), args); err != nil {
		logrus.Fatalf("an error occurred: %+v", err)
	}
}
