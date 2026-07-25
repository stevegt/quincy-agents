// Intent: Parse AGENTS.toml configuration with categories, tags, and module sources.
// Source: DI-jusuk

package toml

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Config   ConfigSection       `toml:"config"`
	Order    OrderSection        `toml:"order"`
	Category map[string]Category `toml:"category"`
	Activate ActivateSection     `toml:"activate"`
	Output   OutputSection       `toml:"output"`
}

type ConfigSection struct {
	ModuleDir string `toml:"module_dir"`
}

type OrderSection struct {
	Categories []string `toml:"categories"`
}

type Category struct {
	Tags    []string `toml:"tags"`
	Modules []Module `toml:"module"`
}

type Module struct {
	Name   string   `toml:"name"`
	Source string   `toml:"source"`
	Tags   []string `toml:"tags"`
}

type ActivateSection struct {
	Scopes []string `toml:"scopes"`
}

type OutputSection struct {
	Path string `toml:"path"`
}

func Parse(path string) (*Config, error) {
	var cfg Config
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return nil, err
	}

	if cfg.Config.ModuleDir == "" {
		cfg.Config.ModuleDir = ".mogent/modules"
	}

	if cfg.Output.Path == "" {
		cfg.Output.Path = "AGENTS.md"
	}

	return &cfg, nil
}

func (c *Config) ResolveModulePath(source string) string {
	if filepath.IsAbs(source) {
		return source
	}

	if len(source) > 0 && source[0] == '~' {
		home, err := os.UserHomeDir()
		if err != nil {
			return source
		}
		return filepath.Join(home, source[1:])
	}

	return filepath.Join(c.Config.ModuleDir, source)
}
