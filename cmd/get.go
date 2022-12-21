package cmd

import (
	"encoding/base64"
	"encoding/binary"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/lucasepe/kvs/internal/aes"
	"github.com/lucasepe/kvs/internal/pbdk"
	"github.com/lucasepe/kvs/internal/store"
	"github.com/lucasepe/toolbox/flags/commander"
	"github.com/lucasepe/toolbox/slug"
)

func newCmdGet() *cmdGet {
	return &cmdGet{}
}

type cmdGet struct {
	itemKey string
	bucket  string
	store   string
	decrypt bool
}

func (*cmdGet) Name() string { return "get" }
func (*cmdGet) Synopsis() string {
	return "Retrieve a value from a bucket."
}
func (*cmdGet) Usage() string {
	return strings.ReplaceAll(`{NAME} [-s store] [-d] -b bucket <key>

   Get the value of the key 'user' from the 'google' bucket:
     {NAME} get -b google user`, "{NAME}", appName)
}

func (p *cmdGet) SetFlags(fs *flag.FlagSet) {
	fs.BoolVar(&p.decrypt, "d", false, "decrypt the value")
	fs.StringVar(&p.bucket, "b", "", "bucket name (required)")
	if def, err := defaultStoreFile(); err == nil {
		fs.StringVar(&p.store, "s", def, fmt.Sprintf("storage file (default: %s)", def))
	} else {
		fs.StringVar(&p.store, "s", "", "storage file (required)")
	}
}

func (p *cmdGet) Execute(fs *flag.FlagSet) commander.ExitStatus {
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

	data, err := db.Get(p.itemKey)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return commander.ExitFailure
	}

	res, err := p.decryptEventually(data)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return commander.ExitFailure
	}

	binary.Write(os.Stdout, binary.LittleEndian, res)

	return commander.ExitSuccess
}

func (p *cmdGet) complete(fs *flag.FlagSet) error {
	if len(p.bucket) == 0 {
		return fmt.Errorf("bucket name is required")
	}

	if fs.NArg() < 1 {
		return fmt.Errorf("key is required")
	}

	p.bucket = slug.Slugify(p.bucket)
	p.itemKey = fs.Arg(0)

	return nil
}

func (p *cmdGet) decryptEventually(dat []byte) ([]byte, error) {
	if !p.decrypt {
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

	enc := base64.RawStdEncoding
	buf := make([]byte, enc.DecodedLen(len(dat)))
	l, err := enc.Decode(buf, dat)
	if err != nil {
		return nil, err
	}

	return aes.GcmDecrypt(buf[:l], key)
}
