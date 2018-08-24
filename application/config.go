package application

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"os"
	"path"
	"runtime"
	"time"
)

const (
	Name = "duci"
)

var (
	Config *Configuration
)

type maskString string

func (s maskString) MarshalJSON() ([]byte, error) {
	return []byte(`"***"`), nil
}

type Configuration struct {
	Server *Server `yaml:"server" json:"server"`
	GitHub *GitHub `yaml:"github" json:"github"`
	Job    *Job    `yaml:"job" json:"job"`
}

type Server struct {
	WorkDir      string `yaml:"workdir" json:"workdir"`
	Port         int    `yaml:"port" json:"port"`
	DatabasePath string `yaml:"database_path" json:"databasePath"`
}

type GitHub struct {
	SSHKeyPath string     `yaml:"ssh_key_path" json:"sshKeyPath"`
	APIToken   maskString `yaml:"api_token" json:"apiToken"`
}

type Job struct {
	Timeout     int64 `yaml:"timeout" json:"timeout"`
	Concurrency int   `yaml:"concurrency" json:"concurrency"`
}

func init() {
	Config = &Configuration{
		Server: &Server{
			WorkDir:      path.Join(os.TempDir(), Name),
			Port:         8080,
			DatabasePath: path.Join(os.Getenv("HOME"), ".duci/db"),
		},
		GitHub: &GitHub{
			SSHKeyPath: path.Join(os.Getenv("HOME"), ".ssh/id_rsa"),
			APIToken:   maskString(os.Getenv("GITHUB_API_TOKEN")),
		},
		Job: &Job{
			Timeout:     600,
			Concurrency: runtime.NumCPU(),
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

func (c *Configuration) Timeout() time.Duration {
	return time.Duration(c.Job.Timeout) * time.Second
}
