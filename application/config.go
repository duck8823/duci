package application

import (
	"encoding/json"
	"fmt"
	"github.com/duck8823/duci/infrastructure/logger"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"os"
	"path"
)

const Name = "duci"

var (
	Config *Configuration
)

type Configuration struct {
	Server *Server `yaml:"server"`
}

type Server struct {
	WorkDir string `yaml:"workdir"`
	Port    int    `yaml:"port"`
}

func init() {
	Config = &Configuration{
		Server: &Server{
			WorkDir: path.Join(os.TempDir(), Name),
			Port:    8080,
		},
	}
}

func (c *Configuration) String() string {
	bytes, err := json.Marshal(c)
	if err != nil {
		logger.Errorf(uuid.UUID{}, "%+v", err)
		return ""
	}
	return string(bytes)
}

func (c *Configuration) Set(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return errors.WithStack(err)
	}
	return yaml.NewDecoder(file).Decode(c)
}

func (c *Configuration) Addr() string {
	return fmt.Sprintf(":%d", c.Server.Port)
}
