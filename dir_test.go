package dotconfig

import (
	"os"
	"testing"
)

func TestDir(t *testing.T) {
	// Save original functions to restore later
	origXdgConfigHome := xdgConfigHome
	origDirExists := dirExists
	origUserHomeDir := userHomeDir

	// Restore original functions after test
	defer func() {
		xdgConfigHome = origXdgConfigHome
		dirExists = origDirExists
		userHomeDir = origUserHomeDir
	}()

	t.Run("XDG config exists", func(t *testing.T) {
		// Mock functions
		xdgConfigHome = func() string { return "/mock/xdg" }
		dirExists = func(dir string) bool {
			return dir == "/mock/xdg/myapp"
		}
		userHomeDir = func() (string, error) {
			return "/mock/home", nil
		}

		dir, exist := Dir("myapp")

		if dir != "/mock/xdg/myapp" {
			t.Errorf("Expected dir to be '/mock/xdg/myapp', got '%s'", dir)
		}
		if !exist {
			t.Error("Expected exist to be true")
		}
	})

	t.Run("XDG config not exist, fallback to home dot dir", func(t *testing.T) {
		// Mock functions
		xdgConfigHome = func() string { return "/mock/xdg" }
		dirExists = func(dir string) bool {
			return dir == "/mock/home/.myapp"
		}
		userHomeDir = func() (string, error) {
			return "/mock/home", nil
		}

		dir, exist := Dir("myapp")

		if dir != "/mock/home/.myapp" {
			t.Errorf("Expected dir to be '/mock/home/.myapp', got '%s'", dir)
		}
		if !exist {
			t.Error("Expected exist to be true")
		}
	})

	t.Run("No XDG, home .config exists", func(t *testing.T) {
		// Mock functions
		xdgConfigHome = func() string { return "" }
		dirExists = func(dir string) bool {
			return dir == "/mock/home/.config/myapp"
		}
		userHomeDir = func() (string, error) {
			return "/mock/home", nil
		}

		dir, exist := Dir("myapp")

		if dir != "/mock/home/.config/myapp" {
			t.Errorf("Expected dir to be '/mock/home/.config/myapp', got '%s'", dir)
		}
		if !exist {
			t.Error("Expected exist to be true")
		}
	})

	t.Run("No XDG, no .config, fallback to home dot dir", func(t *testing.T) {
		// Mock functions
		xdgConfigHome = func() string { return "" }
		dirExists = func(dir string) bool {
			return dir == "/mock/home/.myapp"
		}
		userHomeDir = func() (string, error) {
			return "/mock/home", nil
		}

		dir, exist := Dir("myapp")

		if dir != "/mock/home/.myapp" {
			t.Errorf("Expected dir to be '/mock/home/.myapp', got '%s'", dir)
		}
		if !exist {
			t.Error("Expected exist to be true")
		}
	})

	t.Run("No directories exist, fallback to first option", func(t *testing.T) {
		// Mock functions
		xdgConfigHome = func() string { return "/mock/xdg" }
		dirExists = func(dir string) bool { return false }
		userHomeDir = func() (string, error) {
			return "/mock/home", nil
		}

		dir, exist := Dir("myapp")

		if dir != "/mock/xdg/myapp" {
			t.Errorf("Expected dir to be '/mock/xdg/myapp', got '%s'", dir)
		}
		if exist {
			t.Error("Expected exist to be false")
		}
	})

	t.Run("No locations available, use local dir", func(t *testing.T) {
		// Mock functions
		xdgConfigHome = func() string { return "" }
		dirExists = func(dir string) bool {
			return dir == ".myapp"
		}
		userHomeDir = func() (string, error) {
			return "", os.ErrNotExist
		}

		dir, exist := Dir("myapp")

		if dir != ".myapp" {
			t.Errorf("Expected dir to be '.myapp', got '%s'", dir)
		}
		if !exist {
			t.Error("Expected exist to be true")
		}
	})

	t.Run("No locations available, local dir not exist", func(t *testing.T) {
		// Mock functions
		xdgConfigHome = func() string { return "" }
		dirExists = func(dir string) bool { return false }
		userHomeDir = func() (string, error) {
			return "", os.ErrNotExist
		}

		dir, exist := Dir("myapp")

		if dir != ".myapp" {
			t.Errorf("Expected dir to be '.myapp', got '%s'", dir)
		}
		if exist {
			t.Error("Expected exist to be false")
		}
	})
}

func TestListXDG(t *testing.T) {
	origUserHomeDir := userHomeDir
	defer func() {
		userHomeDir = origUserHomeDir
	}()

	called := 0
	paths := []string{}

	yield := func(path string) bool {
		called++
		paths = append(paths, path)
		return true // Continue iteration
	}

	userHomeDir = func() (string, error) {
		return "/mock/home", nil
	}

	listWithXDG(yield, "myapp", "/mock/xdg")

	if called != 3 {
		t.Errorf("Expected yield to be called 3 times, got %d", called)
	}

	expectedPaths := []string{
		"/mock/xdg/myapp",
		"/mock/home/lib/myapp",
		"/mock/home/.myapp",
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
}

func TestListHome(t *testing.T) {
	origUserHomeDir := userHomeDir
	defer func() {
		userHomeDir = origUserHomeDir
	}()

	called := 0
	paths := []string{}

	yield := func(path string) bool {
		called++
		paths = append(paths, path)
		return true // Continue iteration
	}

	userHomeDir = func() (string, error) {
		return "/mock/home", nil
	}

	listWithNoXDG(yield, "myapp")

	if called != 3 {
		t.Errorf("Expected yield to be called 3 times, got %d", called)
	}

	expectedPaths := []string{
		"/mock/home/.config/myapp",
		"/mock/home/lib/myapp",
		"/mock/home/.myapp",
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
}

func TestListHomeEarlyTermination(t *testing.T) {
	origUserHomeDir := userHomeDir
	defer func() {
		userHomeDir = origUserHomeDir
	}()

	called := 0

	yield := func(path string) bool {
		called++
		return false // Stop iteration after first call
	}

	userHomeDir = func() (string, error) {
		return "/mock/home", nil
	}

	listWithNoXDG(yield, "myapp")

	if called != 1 {
		t.Errorf("Expected yield to be called 1 time, got %d", called)
	}
}

func TestListXDGEarlyTermination(t *testing.T) {
	origUserHomeDir := userHomeDir
	defer func() {
		userHomeDir = origUserHomeDir
	}()

	called := 0

	yield := func(path string) bool {
		called++
		return false // Stop iteration after first call
	}

	userHomeDir = func() (string, error) {
		return "/mock/home", nil
	}

	listWithXDG(yield, "myapp", "/mock/xdg")

	if called != 1 {
		t.Errorf("Expected yield to be called 1 time, got %d", called)
	}
}

func TestList(t *testing.T) {
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

		paths := []string{}
		for path := range list("myapp") {
			paths = append(paths, path)
		}

		expectedPaths := []string{
			"/mock/xdg/myapp",
			"/mock/home/lib/myapp",
			"/mock/home/.myapp",
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

		paths := []string{}
		for path := range list("myapp") {
			paths = append(paths, path)
		}

		expectedPaths := []string{
			"/mock/home/.config/myapp",
			"/mock/home/lib/myapp",
			"/mock/home/.myapp",
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

func TestXdgConfigHome(t *testing.T) {
	expected := "/mock/xdg"
	t.Setenv("XDG_CONFIG_HOME", expected)
	got := xdgConfigHome()
	if got != expected {
		t.Errorf("Expected %v, got %v", expected, got)
	}
}

func TestDirExists(t *testing.T) {
	if !dirExists(".") {
		t.Errorf("Expected true, got false")
	}
	if dirExists("not_exists") {
		t.Errorf("Expected false, got true")
	}
}
