package cmd

import (
	"context"
	"fmt"
	"github.com/duck8823/duci/domain/model/docker"
	"github.com/duck8823/duci/domain/model/job"
	"github.com/duck8823/duci/domain/model/runner"
	"github.com/labstack/gommon/random"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
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

	workspace := filepath.Join(os.TempDir(), random.String(16))
	if err := copyDir(".", workspace); err != nil {
		logrus.Fatalf("an error occurred: %+v", err)
	}

	if err := runner.Run(context.Background(), job.WorkDir(workspace), docker.Tag(random.String(12, random.Hex)), args); err != nil {
		logrus.Fatalf("an error occurred: %+v", err)
	}
	defer os.RemoveAll(workspace)
}

func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return errors.WithStack(err)
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return errors.WithStack(err)
	}

	srcInfo, err := os.Stat(src)
	if err != nil {
		return errors.WithStack(err)
	}
	if err := os.Chmod(dst, srcInfo.Mode()); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func copyDir(src, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return errors.WithStack(err)
	}

	if err = os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return errors.WithStack(err)
	}

	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return errors.WithStack(err)
	}
	for _, entry := range entries {
		srcFilePath := filepath.Join(src, entry.Name())
		dstFilePath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err = copyDir(srcFilePath, dstFilePath); err != nil {
				fmt.Println(err)
			}
		} else {
			if err = copyFile(srcFilePath, dstFilePath); err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
}
