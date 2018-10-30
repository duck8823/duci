package application

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"time"
)

const (
	// Name is a application name.
	Name = "duci"
	// DefaultConfigPath is a path to configuration file
	DefaultConfigurationPath = "./config.yml"
)

var (
	// Config is a application configuration.
	Config *Configuration
)

type maskString string

// MarshalJSON returns masked string.
func (s maskString) MarshalJSON() ([]byte, error) {
	return []byte(`"***"`), nil
}

// Configuration of application.
type Configuration struct {
	Server *Server `yaml:"server" json:"server"`
	GitHub *GitHub `yaml:"github" json:"github"`
	Job    *Job    `yaml:"job" json:"job"`
}

// Server describes a configuration of server.
type Server struct {
	WorkDir      string `yaml:"workdir" json:"workdir"`
	Port         int    `yaml:"port" json:"port"`
	DatabasePath string `yaml:"database_path" json:"databasePath"`
}

// GitHub describes a configuration of github.
type GitHub struct {
	SSHKeyPath string     `yaml:"ssh_key_path" json:"sshKeyPath"`
	APIToken   maskString `yaml:"api_token" json:"apiToken"`
}

// Job describes a configuration of each jobs.
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
			SSHKeyPath: "",
			APIToken:   maskString(os.Getenv("GITHUB_API_TOKEN")),
		},
		Job: &Job{
			Timeout:     600,
			Concurrency: runtime.NumCPU(),
		},
	}
}

// String returns default config path
func (c *Configuration) String() string {
	return ""
}

// Set configuration with file path
func (c *Configuration) Set(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return errors.WithStack(err)
	}
	data, _ := ioutil.ReadAll(file)
	data = []byte(os.ExpandEnv(string(data)))
	return yaml.NewDecoder(bytes.NewReader(data)).Decode(c)
}

// Type returns value type of itself
func (c *Configuration) Type() string {
	return "string"
}

// Addr returns a string of server port
func (c *Configuration) Addr() string {
	return fmt.Sprintf(":%d", c.Server.Port)
}

// Timeout returns timeout duration.
func (c *Configuration) Timeout() time.Duration {
	return time.Duration(c.Job.Timeout) * time.Second
}
