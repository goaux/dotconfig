// Package dotconfig provides utilities for finding application configuration directories
// and files following standard conventions across different operating systems.
//
// While Go's standard library offers [os.UserConfigDir] to find the user's configuration
// directory, this package extends that functionality by searching multiple potential
// locations and considering traditional dot-directory patterns. This approach provides
// more flexibility and better compatibility with existing applications and different
// environment configurations, including Plan9.
//
// It implements the XDG Base Directory Specification (when XDG_CONFIG_HOME is set)
// and falls back to traditional home directory locations. The package's primary functions
// are [Dir], which locates the appropriate configuration directory for an application,
// and [File], which locates specific configuration files.
//
// The search order for configuration directories is:
//
//  1. $XDG_CONFIG_HOME/<app> (if XDG_CONFIG_HOME is set)
//  2. $HOME/.config/<app> (if XDG_CONFIG_HOME is not set)
//  3. $HOME/lib/<app> (for Plan9 compatibility)
//  4. $HOME/.<app> (if [os.UserHomeDir] returns no error)
//  5. .<app> (in current directory, as last resort)
//
// For files, similar locations are searched but with the file name appended to directories
// or with the file extension appended to dot-prefixed application names.
//
// Unlike [os.UserConfigDir] which only returns a single directory recommendation,
// this package actively searches for existing configuration directories and files, providing
// a recommended path even when no directory or file exists yet, making it easier to handle
// both reading from existing configurations and creating new configuration files.
package dotconfig
