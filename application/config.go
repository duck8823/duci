package application

import (
	"encoding/json"
	"fmt"
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
	Server *Server `yaml:"server" json:"server"`
}

type Server struct {
	WorkDir    string `yaml:"workdir" json:"workdir"`
	Port       int    `yaml:"port" json:"port"`
	SSHKeyPath string `yaml:"ssh_key_path" json:"sshKeyPath"`
}

func init() {
	Config = &Configuration{
		Server: &Server{
			WorkDir:    path.Join(os.TempDir(), Name),
			Port:       8080,
			SSHKeyPath: path.Join(os.Getenv("HOME"), ".ssh/id_rsa"),
		},
	}
}

func (c *Configuration) String() string {
	bytes, _ := json.Marshal(c)
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
