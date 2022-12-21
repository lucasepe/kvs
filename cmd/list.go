package cmd

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/lucasepe/kvs/internal/store"
	"github.com/lucasepe/toolbox/flags/commander"
	"github.com/lucasepe/toolbox/textcol"
)

func newCmdList() *cmdList {
	return &cmdList{}
}

type cmdList struct {
	bucket string
	store  string
}

func (*cmdList) Name() string { return "list" }
func (*cmdList) Synopsis() string {
	return "List all buckets or all keys in a bucket."
}
func (*cmdList) Usage() string {
	return strings.ReplaceAll(`{NAME} list [-b bucket]

   List all keys from the 'google' bucket:
     {NAME} list -b google

   List all buckets:
     {NAME} list`, "{NAME}", appName)
}

func (p *cmdList) SetFlags(fs *flag.FlagSet) {
	fs.StringVar(&p.bucket, "b", "", "bucket name")
	if def, err := defaultStoreFile(); err == nil {
		fs.StringVar(&p.store, "s", def, fmt.Sprintf("storage file (default: %s)", def))
	} else {
		fs.StringVar(&p.store, "s", "", "storage file (required)")
	}
}

func (p *cmdList) Execute(fs *flag.FlagSet) commander.ExitStatus {
	db, err := store.New(store.Options{
		BucketName: p.bucket,
		Path:       p.store,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return commander.ExitFailure
	}
	defer db.Close()

	var names []string
	if len(p.bucket) == 0 {
		names = db.Buckets()
	} else {
		names = db.Keys()
	}

	textcol.PrintColumns(os.Stdout, &names, 3)

	return commander.ExitSuccess
}
