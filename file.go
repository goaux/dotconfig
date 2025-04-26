package dotconfig

import (
	"iter"
	"os"
	"path/filepath"
)

type fileExists int

const (
	// NotExists indicates that neither the file nor its base directory exists
	NotExists fileExists = iota

	// BaseExists indicates that the file's base directory exists but the file itself does not
	BaseExists

	// FileExists indicates that the file exists
	FileExists
)

// File searches for a configuration file for the specified application.
// It follows similar conventions to [Dir] but locates specific files rather than directories.
//
// The function tries the following locations in order:
//
//  1. $XDG_CONFIG_HOME/<app>/<name> (if XDG_CONFIG_HOME is set)
//  2. $HOME/.config/<app>/<name> (if XDG_CONFIG_HOME is not set)
//  3. $HOME/lib/<app>/<name> (for Plan9 compatibility)
//  4. $HOME/.<app>/<name> (if [os.UserHomeDir] returns no error)
//  5. $HOME/.<app><ext> (where <ext> is the file extension of <name>, if [os.UserHomeDir] returns no error)
//  6. .<app>/<name> (in current directory)
//  7. .<app><ext> (in current directory, as last resort)
//
// If the file name parameter is "." or "/", the application name is used as the file name.
//
// Parameters:
//   - app: The application name to search configurations for
//   - name: The name of the configuration file to find
//
// Returns:
//   - path: The configuration file path
//   - status: A fileExists constant indicating whether the file exists, only its base directory exists, or neither exists
func File(app, name string) (path string, status fileExists) {
	cfg := newFileConfig(app, name)
	var fallback string
	for file := range cfg.List() {
		if check := checkFile(file); check == FileExists {
			return file, check
		}
		if fallback == "" {
			fallback = file
		}
	}
	if fallback == "" {
		file := filepath.Join("."+app, cfg.File)
		if check := checkFile(file); check == FileExists {
			return file, check
		}
		fallback = "." + app + filepath.Ext(cfg.File)
	}
	return fallback, checkFile(fallback)
}

var checkFile = func(name string) fileExists {
	if _, err := os.Stat(name); err == nil { // if NO error
		return FileExists
	}
	if _, err := os.Stat(filepath.Dir(name)); err == nil { // if NO error
		return BaseExists
	}
	return NotExists
}

type fileConfig struct {
	App  string
	File string
}

func newFileConfig(app, file string) *fileConfig {
	file = filepath.Base(file)
	if file == "." || file == "/" {
		file = app
	}
	return &fileConfig{App: app, File: file}
}

func (cfg *fileConfig) List() iter.Seq[string] {
	return func(yield func(string) bool) {
		if xdg := xdgConfigHome(); xdg != "" {
			cfg.ListWithXDG(yield, xdg)
		} else {
			cfg.ListWithNoXDG(yield)
		}
	}
}

func (cfg *fileConfig) ListWithXDG(yield func(string) bool, xdg string) {
	if yield(filepath.Join(xdg, cfg.App, cfg.File)) {
		if home, err := userHomeDir(); err == nil { // if NO error
			cfg.ListHome(yield, home)
		}
	}
}

func (cfg *fileConfig) ListWithNoXDG(yield func(string) bool) {
	if home, err := userHomeDir(); err == nil { // if NO error
		if yield(filepath.Join(home, ".config", cfg.App, cfg.File)) {
			cfg.ListHome(yield, home)
		}
	}
}

func (cfg *fileConfig) ListHome(yield func(string) bool, home string) {
	if yield(filepath.Join(home, "lib", cfg.App, cfg.File)) {
		if yield(filepath.Join(home, "."+cfg.App, cfg.File)) {
			yield(filepath.Join(home, "."+cfg.App+filepath.Ext(cfg.File)))
		}
	}
}
