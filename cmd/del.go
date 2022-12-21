package cmd

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/lucasepe/kvs/internal/store"
	"github.com/lucasepe/toolbox/flags/commander"
	"github.com/lucasepe/toolbox/slug"
)

func newCmdDelete() *cmdDelete {
	return &cmdDelete{}
}

type cmdDelete struct {
	bucket  string
	itemKey string
	store   string
}

func (*cmdDelete) Name() string { return "del" }
func (*cmdDelete) Synopsis() string {
	return "Delete a bucket or a key from a bucket."
}
func (*cmdDelete) Usage() string {
	return strings.ReplaceAll(`{NAME} del [-s store] -b google [key]

   Delete the 'google' bucket:
     {NAME} del -b google

   Delete the key 'user' from the 'google' bucket:
     {NAME} del -b google user`, "{NAME}", appName)
}

func (p *cmdDelete) SetFlags(fs *flag.FlagSet) {
	fs.StringVar(&p.bucket, "b", "", "bucket name")
	fs.StringVar(&p.itemKey, "k", "", "key (required)")
	if def, err := defaultStoreFile(); err == nil {
		fs.StringVar(&p.store, "s", def, fmt.Sprintf("storage file (default: %s)", def))
	} else {
		fs.StringVar(&p.store, "s", "", "storage file (required)")
	}
}

func (p *cmdDelete) Execute(fs *flag.FlagSet) commander.ExitStatus {
	if err := p.complete(fs); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return commander.ExitFailure
	}

	db, err := store.New(store.Options{
		BucketName: p.bucket,
		Path:       p.store,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return commander.ExitFailure
	}
	defer db.Close()

	if len(p.itemKey) == 0 {
		err := db.DeleteBucket(p.bucket)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return commander.ExitFailure
		}

		return commander.ExitSuccess
	}

	err = db.Delete(p.itemKey)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return commander.ExitFailure
	}

	return commander.ExitSuccess
}

func (p *cmdDelete) complete(fs *flag.FlagSet) error {
	if len(p.bucket) == 0 {
		return fmt.Errorf("bucket name is required")
	}
	p.bucket = slug.Slugify(p.bucket)

	if fs.NArg() > 0 {
		p.itemKey = fs.Arg(0)
	}

	return nil
}
