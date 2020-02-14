/*
Copyright © 2020 Luca Sepe <luca.sepe@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/avelino/slugify"
	"github.com/lucasepe/kvs"
	"github.com/lucasepe/kvs/pkg/aes"
	"github.com/lucasepe/kvs/pkg/cl"
	"github.com/lucasepe/kvs/pkg/pbdk"
	"github.com/spf13/cobra"
)

const (
	maxFileSize = 1024 * 1024 // 1 MB
	optEncrypt  = "encrypt"
)

// pushCmd represents the set command
var pushCmd = &cobra.Command{
	Use:                   "push <key> <value>",
	DisableSuggestions:    true,
	DisableFlagsInUseLine: false,
	Args:                  cobra.MinimumNArgs(1),
	Short:                 "Update a store adding an item with the specified key in a bucket",
	Example: `  Push a file content (using pipes):
    cat avatar.jpg | kvs push avatar --bucket google --store accounts

  Push a command output (using pipes):
    pwgen 14 1 | kvs push -b instagram -s 'my secrets'

  Store a 'google' account 'user' field:
    kvs push --bucket google --store accounts user luca.sepe@gmail.com`,
	Run: func(cmd *cobra.Command, args []string) {
		bucket, err := cmd.Flags().GetString(optBucket)
		cl.HandleErr(err)
		bucket = slugify.Slugify(bucket)

		key := slugify.Slugify(args[0])

		var reader io.Reader

		info, err := os.Stdin.Stat()
		cl.HandleErr(err)

		if (info.Mode() & os.ModeCharDevice) != os.ModeCharDevice {
			reader = io.LimitReader(bufio.NewReader(os.Stdin), maxFileSize)
		} else if len(args) > 1 {
			reader = strings.NewReader(args[1])
		}

		if reader == nil {
			fmt.Fprintln(os.Stderr, "undefined item to store")
			os.Exit(1)
		}

		data, err := ioutil.ReadAll(reader)
		cl.HandleErr(err)
		if len(data) == 0 {
			return
		}

		if ok, _ := cmd.Flags().GetBool(optEncrypt); ok {
			secret, err := cl.GetSecret("Secret phrase: ")
			cl.HandleErr(err)
			secret2, err := cl.GetSecret("Secret phrase again: ")
			cl.HandleErr(err)
			if ok := bytes.Equal(secret, secret2); !ok {
				fmt.Fprintf(os.Stderr, "secret phrases do not match\n")
				os.Exit(1)
			}

			res, err := encrypt(data, secret)
			cl.HandleErr(err)

			data = res
		}

		store, err := kvs.Open(getStoreFile())
		cl.HandleErr(err)
		defer store.Close()

		if err = store.Put(bucket, key, data); err != nil {
			fmt.Fprintf(os.Stderr, "failed to store item with key '%s' (bucket: %s)\n", key, bucket)
			os.Exit(1)
		}

		fmt.Printf("item successfully stored in bucket '%s' with key '%s'\n", bucket, key)
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)

	pushCmd.Flags().BoolP(optEncrypt, "e", false, "encrypt data using AES algorithm in CBC mode with PKCS7 padding")

	pushCmd.Flags().StringP(optBucket, "b", "", "bucket name")
	pushCmd.MarkFlagRequired(optBucket)
}

func encrypt(data, secret []byte) ([]byte, error) {
	key, err := pbdk.DeriveKey(secret)
	if err != nil {
		return nil, err
	}

	enc, err := aes.GcmEncrypt(data, key)
	if err != nil {
		return nil, err
	}

	res := make([]byte, base64.RawStdEncoding.EncodedLen(len(enc)))
	base64.RawStdEncoding.Encode(res, enc)

	return res, nil
}
