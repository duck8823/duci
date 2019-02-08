package runner

import (
	"bytes"
	"github.com/duck8823/duci/domain/model/docker"
	"github.com/duck8823/duci/domain/model/job"
	"github.com/duck8823/duci/infrastructure/archive/tar"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

// createTarball creates a tar archive
func createTarball(workDir job.WorkDir) (*os.File, error) {
	tarFilePath := filepath.Join(workDir.String(), "duci.tar")
	writeFile, err := os.OpenFile(tarFilePath, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer writeFile.Close()

	if err := tar.Create(workDir.String(), writeFile); err != nil {
		return nil, errors.WithStack(err)
	}

	readFile, _ := os.Open(tarFilePath)
	return readFile, nil
}

// dockerfilePath returns a path to dockerfile for duci using
func dockerfilePath(workDir job.WorkDir) docker.Dockerfile {
	dockerfile := "./Dockerfile"
	if exists(filepath.Join(workDir.String(), ".duci/Dockerfile")) {
		dockerfile = ".duci/Dockerfile"
	}
	return docker.Dockerfile(filepath.Join(workDir.String(), dockerfile))
}

// exists indicates whether the file exists
func exists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

// runtimeOptions parses a config.yml and returns a docker runtime options
func runtimeOptions(workDir job.WorkDir) (docker.RuntimeOptions, error) {
	var opts docker.RuntimeOptions

	if !exists(filepath.Join(workDir.String(), ".duci/config.yml")) {
		return opts, nil
	}
	content, err := ioutil.ReadFile(filepath.Join(workDir.String(), ".duci/config.yml"))
	if err != nil {
		return opts, errors.WithStack(err)
	}
	content = []byte(os.ExpandEnv(string(content)))
	if err := yaml.NewDecoder(bytes.NewReader(content)).Decode(&opts); err != nil {
		return opts, errors.WithStack(err)
	}
	return opts, nil
}
