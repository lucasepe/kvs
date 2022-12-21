package cmd

import (
	"flag"
	"fmt"

	"github.com/lucasepe/toolbox/flags/commander"
)

func newCmdVersion(ver, bld string) *cmdVersion {
	return &cmdVersion{
		version: ver,
		build:   bld,
	}
}

type cmdVersion struct {
	version string
	build   string
}

func (*cmdVersion) Name() string { return "version" }

func (*cmdVersion) Synopsis() string {
	return "Print the current build information."
}

func (*cmdVersion) Usage() string {
	return fmt.Sprintf("%s version", appName)
}

func (p *cmdVersion) SetFlags(fs *flag.FlagSet) {}

func (p *cmdVersion) Execute(fs *flag.FlagSet) commander.ExitStatus {
	fmt.Printf("Key Value Store %s (build: %s)\n", p.version, p.build)
	return commander.ExitSuccess
}
