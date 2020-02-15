package cl

import (
	"fmt"
	"os"
	"path"
	"syscall"

	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh/terminal"
)

// GetSecret read a password from the terminal
func GetSecret(prompt, envKey string) ([]byte, error) {
	if val, ok := os.LookupEnv(envKey); ok && val != "" {
		return []byte(val), nil
	}

	fmt.Fprint(os.Stderr, prompt)
	var fd int
	if terminal.IsTerminal(int(syscall.Stdin)) {
		fd = int(syscall.Stdin)
	} else {
		tty, err := os.Open("/dev/tty")
		if err != nil {
			return nil, errors.Wrap(err, "error allocating terminal")
		}
		defer tty.Close()
		fd = int(tty.Fd())
	}

	pass, err := terminal.ReadPassword(fd)
	fmt.Fprintln(os.Stderr)
	return pass, err
}

// HandleErr check for an error and eventually exit the app
func HandleErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// GetWorkDir returns the application working directory
// creating it if does not exists.
func GetWorkDir(appName string) (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}

	workdir := path.Join(home, fmt.Sprintf(".%s", appName))
	if _, err := os.Stat(workdir); os.IsNotExist(err) {
		err = os.MkdirAll(workdir, os.ModePerm)
		if err != nil {
			return "", err
		}
	}

	return workdir, nil
}

// FileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
