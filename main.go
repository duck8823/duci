package main

import (
	"context"
	"fmt"
	"github.com/duck8823/duci/application/service/executor"
	"github.com/duck8823/duci/domain/model/job"
	"github.com/duck8823/duci/domain/model/job/target"
	"github.com/duck8823/duci/domain/model/runner"
)

func main() {
	e := &executor.JobExecutor{
		DockerRunner: runner.
			DefaultDockerRunnerBuilder().
			LogFunc(func(_ context.Context, log job.Log) {
				for line, err := log.ReadLine(); err == nil; line, err = log.ReadLine() {
					println(line.Message)
				}
			}).
			Build(),
		StartFunc: func(_ context.Context) {
			println("Job Started")
		},
		EndFunc: func(_ context.Context, err error) {
			if err != nil {
				println(fmt.Sprintf("%+v", err))
			}
			println("Job End")
		},
	}

	if err := e.Execute(context.Background(), &target.Local{Path: "."}); err != nil {
		println(fmt.Sprintf("%+v", err))
	}
}
