package cmd

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/lucasepe/toolbox/flags/commander"
	"github.com/lucasepe/toolbox/xdg"
)

const (
	appName   = "kvs"
	envSecret = "KVS_SECRET"
	banner    = `┬┌─  ┬  ┬    ┌─┐
├┴┐  └┐┌┘    └─┐
┴ ┴ey └┘alue └─┘tore
`
)

func Run(ver, bld string) {
	app := commander.New(flag.CommandLine, appName)
	app.Banner = banner
	app.Register(app.HelpCommand(), "")
	app.Register(newCmdVersion(ver, bld), "")
	app.Register(newCmdSet(), "")
	app.Register(newCmdList(), "")
	app.Register(newCmdGet(), "")
	app.Register(newCmdDelete(), "")

	flag.Parse()

	os.Exit(int(app.Execute()))
}

func storeDir() (string, error) {
	dir := filepath.Join(xdg.ConfigDir(), appName)
	err := os.MkdirAll(dir, os.ModePerm)
	return dir, err
}

func storeFile(name string) (string, error) {
	dir, err := storeDir()
	if err != nil {
		return "", err
	}

	fn := name[:len(name)-len(filepath.Ext(name))]

	return filepath.Join(dir, fmt.Sprintf("%s.kvs", fn)), nil
}

func defaultStoreFile() (string, error) {
	return storeFile("secrets")
}
