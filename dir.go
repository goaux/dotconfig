package dotconfig

import (
	"iter"
	"os"
	"path/filepath"
)

// Dir searches for a configuration directory for the specified application.
// It searches multiple potential locations following standard conventions
// across different operating systems.
//
// The function tries the following locations in order:
//
//  1. $XDG_CONFIG_HOME/<app> (if XDG_CONFIG_HOME is set)
//  2. $HOME/.config/<app> (if XDG_CONFIG_HOME is not set)
//  3. $HOME/lib/<app> (for Plan9 compatibility)
//  4. $HOME/.<app> (if [os.UserHomeDir] returns no error)
//  5. .<app> (in current directory, as last resort)
//
// If an existing directory is found, it returns the directory path and true.
// If no existing directory is found but potential locations were checked,
// it returns the first potential location and false.
// If no locations could be determined, it returns ".<app>" and whether it exists.
//
// Parameters:
//   - app: The application name to search configurations for
//
// Returns:
//   - dir: The configuration directory path
//   - exist: Boolean indicating whether the directory exists on the filesystem
func Dir(app string) (dir string, exist bool) {
	var fallback string
	for dir := range list(app) {
		if dirExists(dir) {
			return dir, true
		}
		if fallback == "" {
			fallback = dir
		}
	}
	if fallback == "" {
		dir = "." + app
		return dir, dirExists(dir)
	}
	return fallback, false
}

func list(app string) iter.Seq[string] {
	return func(yield func(string) bool) {
		if xdg := xdgConfigHome(); xdg != "" {
			listWithXDG(yield, app, xdg)
		} else {
			listWithNoXDG(yield, app)
		}
	}
}

func listWithXDG(yield func(string) bool, app, xdg string) {
	if yield(filepath.Join(xdg, app)) {
		if home, err := userHomeDir(); err == nil { // if NO error
			listHome(yield, home, app)
		}
	}
}

func listWithNoXDG(yield func(string) bool, app string) {
	if home, err := userHomeDir(); err == nil { // if NO error
		if yield(filepath.Join(home, ".config", app)) {
			listHome(yield, home, app)
		}
	}
}

func listHome(yield func(string) bool, home, app string) {
	if yield(filepath.Join(home, "lib", app)) {
		yield(filepath.Join(home, "."+app))
	}
}

var xdgConfigHome = func() string {
	return os.Getenv("XDG_CONFIG_HOME")
}

var dirExists = func(dir string) bool {
	info, err := os.Stat(dir)
	return err == nil && info.IsDir()
}

var userHomeDir = os.UserHomeDir
