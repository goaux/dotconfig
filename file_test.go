package dotconfig

import (
	"os"
	"testing"
)

func TestFile(t *testing.T) {
	// Save original functions to restore later
	origXdgConfigHome := xdgConfigHome
	origCheckFile := checkFile
	origUserHomeDir := userHomeDir

	// Restore original functions after test
	defer func() {
		xdgConfigHome = origXdgConfigHome
		checkFile = origCheckFile
		userHomeDir = origUserHomeDir
	}()

	t.Run("XDG config file exists", func(t *testing.T) {
		// Mock functions
		xdgConfigHome = func() string { return "/mock/xdg" }
		checkFile = func(path string) fileExists {
			if path == "/mock/xdg/myapp/config.yaml" {
				return FileExists
			}
			return NotExists
		}
		userHomeDir = func() (string, error) {
			return "/mock/home", nil
		}

		path, status := File("myapp", "config.yaml")

		if path != "/mock/xdg/myapp/config.yaml" {
			t.Errorf("Expected path to be '/mock/xdg/myapp/config.yaml', got '%s'", path)
		}
		if status != FileExists {
			t.Errorf("Expected status to be FileExists (%d), got %d", FileExists, status)
		}
	})

	t.Run("XDG config base exists, file not exists", func(t *testing.T) {
		// Mock functions
		xdgConfigHome = func() string { return "/mock/xdg" }
		checkFile = func(path string) fileExists {
			if path == "/mock/xdg/myapp/config.yaml" {
				return BaseExists
			}
			return NotExists
		}
		userHomeDir = func() (string, error) {
			return "/mock/home", nil
		}

		path, status := File("myapp", "config.yaml")

		if path != "/mock/xdg/myapp/config.yaml" {
			t.Errorf("Expected path to be '/mock/xdg/myapp/config.yaml', got '%s'", path)
		}
		if status != BaseExists {
			t.Errorf("Expected status to be BaseExists (%d), got %d", BaseExists, status)
		}
	})

	t.Run("XDG not exist, home .config file exists", func(t *testing.T) {
		// Mock functions
		xdgConfigHome = func() string { return "" }
		checkFile = func(path string) fileExists {
			if path == "/mock/home/.config/myapp/config.yaml" {
				return FileExists
			}
			return NotExists
		}
		userHomeDir = func() (string, error) {
			return "/mock/home", nil
		}

		path, status := File("myapp", "config.yaml")

		if path != "/mock/home/.config/myapp/config.yaml" {
			t.Errorf("Expected path to be '/mock/home/.config/myapp/config.yaml', got '%s'", path)
		}
		if status != FileExists {
			t.Errorf("Expected status to be FileExists (%d), got %d", FileExists, status)
		}
	})

	t.Run("XDG not exist, home lib file exists", func(t *testing.T) {
		// Mock functions
		xdgConfigHome = func() string { return "" }
		checkFile = func(path string) fileExists {
			if path == "/mock/home/lib/myapp/config.yaml" {
				return FileExists
			}
			return NotExists
		}
		userHomeDir = func() (string, error) {
			return "/mock/home", nil
		}

		path, status := File("myapp", "config.yaml")

		if path != "/mock/home/lib/myapp/config.yaml" {
			t.Errorf("Expected path to be '/mock/home/lib/myapp/config.yaml', got '%s'", path)
		}
		if status != FileExists {
			t.Errorf("Expected status to be FileExists (%d), got %d", FileExists, status)
		}
	})

	t.Run("XDG not exist, dot file exists", func(t *testing.T) {
		// Mock functions
		xdgConfigHome = func() string { return "" }
		checkFile = func(path string) fileExists {
			if path == "/mock/home/.myapp.yaml" {
				return FileExists
			}
			return NotExists
		}
		userHomeDir = func() (string, error) {
			return "/mock/home", nil
		}

		path, status := File("myapp", "config.yaml")

		if path != "/mock/home/.myapp.yaml" {
			t.Errorf("Expected path to be '/mock/home/.myapp.yaml', got '%s'", path)
		}
		if status != FileExists {
			t.Errorf("Expected status to be FileExists (%d), got %d", FileExists, status)
		}
	})

	t.Run("No files exist, fallback to first option", func(t *testing.T) {
		// Mock functions
		xdgConfigHome = func() string { return "/mock/xdg" }
		checkFile = func(path string) fileExists { return NotExists }
		userHomeDir = func() (string, error) {
			return "/mock/home", nil
		}

		path, status := File("myapp", "config.yaml")

		if path != "/mock/xdg/myapp/config.yaml" {
			t.Errorf("Expected path to be '/mock/xdg/myapp/config.yaml', got '%s'", path)
		}
		if status != NotExists {
			t.Errorf("Expected status to be NotExists (%d), got %d", NotExists, status)
		}
	})

	t.Run("No locations available, use local dot directory file", func(t *testing.T) {
		// Mock functions
		xdgConfigHome = func() string { return "" }
		checkFile = func(path string) fileExists {
			if path == ".myapp/config.yaml" {
				return FileExists
			}
			return NotExists
		}
		userHomeDir = func() (string, error) {
			return "", os.ErrNotExist
		}

		path, status := File("myapp", "config.yaml")

		if path != ".myapp/config.yaml" {
			t.Errorf("Expected path to be '.myapp/config.yaml', got '%s'", path)
		}
		if status != FileExists {
			t.Errorf("Expected status to be FileExists (%d), got %d", FileExists, status)
		}
	})

	t.Run("No locations available, use local dot file", func(t *testing.T) {
		// Mock functions
		xdgConfigHome = func() string { return "" }
		checkFile = func(path string) fileExists {
			if path == ".myapp.yaml" {
				return FileExists
			}
			return NotExists
		}
		userHomeDir = func() (string, error) {
			return "", os.ErrNotExist
		}

		path, status := File("myapp", "config.yaml")

		if path != ".myapp.yaml" {
			t.Errorf("Expected path to be '.myapp.yaml', got '%s'", path)
		}
		if status != FileExists {
			t.Errorf("Expected status to be FileExists (%d), got %d", FileExists, status)
		}
	})

	t.Run("No locations available, local file not exist", func(t *testing.T) {
		// Mock functions
		xdgConfigHome = func() string { return "" }
		checkFile = func(path string) fileExists { return NotExists }
		userHomeDir = func() (string, error) {
			return "", os.ErrNotExist
		}

		path, status := File("myapp", "config.yaml")

		if path != ".myapp.yaml" {
			t.Errorf("Expected path to be '.myapp.yaml', got '%s'", path)
		}
		if status != NotExists {
			t.Errorf("Expected status to be NotExists (%d), got %d", NotExists, status)
		}
	})

	t.Run("Use app name when file name is '.'", func(t *testing.T) {
		// Mock functions
		xdgConfigHome = func() string { return "/mock/xdg" }
		checkFile = func(path string) fileExists {
			if path == "/mock/xdg/myapp/myapp" {
				return FileExists
			}
			return NotExists
		}
		userHomeDir = func() (string, error) {
			return "/mock/home", nil
		}

		path, status := File("myapp", ".")

		if path != "/mock/xdg/myapp/myapp" {
			t.Errorf("Expected path to be '/mock/xdg/myapp/myapp', got '%s'", path)
		}
		if status != FileExists {
			t.Errorf("Expected status to be FileExists (%d), got %d", FileExists, status)
		}
	})

	t.Run("Use app name when file name is '/'", func(t *testing.T) {
		// Mock functions
		xdgConfigHome = func() string { return "/mock/xdg" }
		checkFile = func(path string) fileExists {
			if path == "/mock/xdg/myapp/myapp" {
				return FileExists
			}
			return NotExists
		}
		userHomeDir = func() (string, error) {
			return "/mock/home", nil
		}

		path, status := File("myapp", "/")

		if path != "/mock/xdg/myapp/myapp" {
			t.Errorf("Expected path to be '/mock/xdg/myapp/myapp', got '%s'", path)
		}
		if status != FileExists {
			t.Errorf("Expected status to be FileExists (%d), got %d", FileExists, status)
		}
	})
}

func TestFileConfig_List(t *testing.T) {
	// Save original functions to restore later
	origXdgConfigHome := xdgConfigHome
	origUserHomeDir := userHomeDir

	// Restore original functions after test
	defer func() {
		xdgConfigHome = origXdgConfigHome
		userHomeDir = origUserHomeDir
	}()

	t.Run("with XDG_CONFIG_HOME", func(t *testing.T) {
		xdgConfigHome = func() string { return "/mock/xdg" }
		userHomeDir = func() (string, error) {
			return "/mock/home", nil
		}

		cfg := newFileConfig("myapp", "config.yaml")
		paths := []string{}
		for path := range cfg.List() {
			paths = append(paths, path)
		}

		expectedPaths := []string{
			"/mock/xdg/myapp/config.yaml",
			"/mock/home/lib/myapp/config.yaml",
			"/mock/home/.myapp/config.yaml",
			"/mock/home/.myapp.yaml",
		}

		if len(paths) != len(expectedPaths) {
			t.Errorf("Expected %d paths, got %d", len(expectedPaths), len(paths))
		}

		for i, expected := range expectedPaths {
			if i >= len(paths) {
				t.Errorf("Missing expected path at index %d: %s", i, expected)
				continue
			}
			if paths[i] != expected {
				t.Errorf("Expected path[%d] to be '%s', got '%s'", i, expected, paths[i])
			}
		}
	})

	t.Run("without XDG_CONFIG_HOME", func(t *testing.T) {
		xdgConfigHome = func() string { return "" }
		userHomeDir = func() (string, error) {
			return "/mock/home", nil
		}

		cfg := newFileConfig("myapp", "config.yaml")
		paths := []string{}
		for path := range cfg.List() {
			paths = append(paths, path)
		}

		expectedPaths := []string{
			"/mock/home/.config/myapp/config.yaml",
			"/mock/home/lib/myapp/config.yaml",
			"/mock/home/.myapp/config.yaml",
			"/mock/home/.myapp.yaml",
		}

		if len(paths) != len(expectedPaths) {
			t.Errorf("Expected %d paths, got %d", len(expectedPaths), len(paths))
		}

		for i, expected := range expectedPaths {
			if i >= len(paths) {
				t.Errorf("Missing expected path at index %d: %s", i, expected)
				continue
			}
			if paths[i] != expected {
				t.Errorf("Expected path[%d] to be '%s', got '%s'", i, expected, paths[i])
			}
		}
	})
}

func TestNewFileConfig(t *testing.T) {
	t.Run("normal file name", func(t *testing.T) {
		cfg := newFileConfig("myapp", "config.yaml")

		if cfg.App != "myapp" {
			t.Errorf("Expected App to be 'myapp', got '%s'", cfg.App)
		}

		if cfg.File != "config.yaml" {
			t.Errorf("Expected File to be 'config.yaml', got '%s'", cfg.File)
		}
	})

	t.Run("file name is dot", func(t *testing.T) {
		cfg := newFileConfig("myapp", ".")

		if cfg.App != "myapp" {
			t.Errorf("Expected App to be 'myapp', got '%s'", cfg.App)
		}

		if cfg.File != "myapp" {
			t.Errorf("Expected File to be 'myapp', got '%s'", cfg.File)
		}
	})

	t.Run("file name is slash", func(t *testing.T) {
		cfg := newFileConfig("myapp", "/")

		if cfg.App != "myapp" {
			t.Errorf("Expected App to be 'myapp', got '%s'", cfg.App)
		}

		if cfg.File != "myapp" {
			t.Errorf("Expected File to be 'myapp', got '%s'", cfg.File)
		}
	})

	t.Run("file name with path", func(t *testing.T) {
		cfg := newFileConfig("myapp", "/path/to/config.yaml")

		if cfg.App != "myapp" {
			t.Errorf("Expected App to be 'myapp', got '%s'", cfg.App)
		}

		if cfg.File != "config.yaml" {
			t.Errorf("Expected File to be 'config.yaml', got '%s'", cfg.File)
		}
	})
}

func TestCheckFile(t *testing.T) {
	testCases := []struct {
		File     string
		Expected fileExists
	}{
		{"file_test.go", FileExists},
		{"not_exists", BaseExists},
		{"not_exists/not_exists", NotExists},
	}

	for _, tc := range testCases {
		if got := checkFile(tc.File); got != tc.Expected {
			t.Errorf("Expected %v, got %v", tc.Expected, got)
		}
	}
}
