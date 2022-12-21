package cmd

import (
	"bufio"
	"encoding/base64"
	"flag"
	"fmt"
	"io"

	"os"
	"strings"

	"github.com/lucasepe/kvs/internal/aes"
	"github.com/lucasepe/kvs/internal/pbdk"
	"github.com/lucasepe/kvs/internal/store"
	"github.com/lucasepe/toolbox/flags/commander"
	"github.com/lucasepe/toolbox/slug"
)

const (
	maxFileSize = 1024 * 1024 // 1 MB
)

func newCmdSet() *cmdSet {
	return &cmdSet{}
}

type cmdSet struct {
	itemKey string
	bucket  string
	store   string
	encrypt bool
}

func (*cmdSet) Name() string { return "set" }
func (*cmdSet) Synopsis() string {
	return "Save a key/value pair to a bucket."
}
func (*cmdSet) Usage() string {
	return strings.ReplaceAll(`{NAME} set [-s store] [-e] -b bucket <key> <value>
  
   Save the value 'my@gmail.com' with the key 'user' into the 'google' bucket:
     {NAME} set -b google user my@gmail.com

   Save the content of the 'doc.txt' file with the key 'doc' using pipes:
     cat doc.txt | {NAME} set -b google doc

   Save a command output using pipes:
     pwgen 14 1 | {NAME} set -b instagram pass`, "{NAME}", appName)
}

func (p *cmdSet) SetFlags(fs *flag.FlagSet) {
	fs.BoolVar(&p.encrypt, "e", false, "encrypt the value")
	fs.StringVar(&p.bucket, "b", "", "bucket name (required)")
	if def, err := defaultStoreFile(); err == nil {
		fs.StringVar(&p.store, "s", def, fmt.Sprintf("storage file (default: %s)", def))
	} else {
		fs.StringVar(&p.store, "s", "", "storage file (required)")
	}
}

func (p *cmdSet) Execute(fs *flag.FlagSet) commander.ExitStatus {
	src, err := p.complete(fs)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return commander.ExitFailure
	}

	if len(src) == 0 {
		return commander.ExitSuccess
	}

	dat, err := p.encryptEventually(src)
	if err != nil {
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

	if err := db.Set(p.itemKey, dat); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return commander.ExitFailure
	}

	fmt.Printf("value with key '%s' successfully saved to '%s'\n", p.itemKey, p.store)

	return commander.ExitSuccess
}

func (p *cmdSet) complete(fs *flag.FlagSet) ([]byte, error) {
	if len(p.bucket) == 0 {
		return nil, fmt.Errorf("bucket name is required")
	}

	if fs.NArg() == 0 {
		return nil, fmt.Errorf("item key is required")
	}

	p.bucket = slug.Slugify(p.bucket)
	p.itemKey = fs.Arg(0)

	var reader io.Reader

	info, err := os.Stdin.Stat()
	if err != nil {
		return nil, err
	}

	if (info.Mode() & os.ModeCharDevice) != os.ModeCharDevice {
		reader = io.LimitReader(bufio.NewReader(os.Stdin), maxFileSize)
	} else if fs.NArg() > 1 {
		reader = strings.NewReader(fs.Arg(1))
	}

	if reader == nil {
		return nil, fmt.Errorf("the value to store has not been specified")
	}

	return io.ReadAll(reader)
}

func (p *cmdSet) encryptEventually(dat []byte) ([]byte, error) {
	if !p.encrypt {
		return dat, nil
	}

	sec, ok := os.LookupEnv(envSecret)
	if !ok || len(sec) == 0 {
		return dat, nil
	}

	key, err := pbdk.DeriveKey([]byte(sec))
	if err != nil {
		return nil, err
	}

	src, err := aes.GcmEncrypt(dat, key)
	if err != nil {
		return nil, err
	}

	enc := base64.RawStdEncoding
	buf := make([]byte, enc.EncodedLen(len(src)))
	enc.Encode(buf, src)

	return buf, nil
}
