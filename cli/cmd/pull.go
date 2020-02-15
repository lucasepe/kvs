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
	"encoding/base64"
	"encoding/binary"
	"os"

	"github.com/avelino/slugify"
	"github.com/lucasepe/kvs"
	"github.com/lucasepe/kvs/pkg/aes"
	"github.com/lucasepe/kvs/pkg/cl"
	"github.com/lucasepe/kvs/pkg/pbdk"
	"github.com/spf13/cobra"
)

const (
	optDecrypt = "decrypt"
)

// pullCmd represents the get command
var pullCmd = &cobra.Command{
	Use:                   "pull <key>",
	DisableSuggestions:    true,
	DisableFlagsInUseLine: false,
	Args:                  cobra.MinimumNArgs(1),
	Short:                 "Fetch from a store the item with the specified key in a bucket",
	Run: func(cmd *cobra.Command, args []string) {
		bucket, err := cmd.Flags().GetString(optBucket)
		cl.HandleErr(err)
		bucket = slugify.Slugify(bucket)

		key := slugify.Slugify(args[0])

		store, err := kvs.Open(getStoreFile())
		cl.HandleErr(err)
		defer store.Close()

		data, err := store.Get(bucket, key)
		cl.HandleErr(err)

		if ok, _ := cmd.Flags().GetBool(optDecrypt); ok {
			secret, err := cl.GetSecret("Secret phrase: ", envSecretKey)
			cl.HandleErr(err)

			res, err := decrypt(data, secret)
			cl.HandleErr(err)

			data = res
		}

		binary.Write(os.Stdout, binary.LittleEndian, data)
	},
}

func init() {
	rootCmd.AddCommand(pullCmd)

	pullCmd.Flags().StringP(optBucket, "b", "", "bucket name")
	pullCmd.MarkFlagRequired(optBucket)

	pullCmd.Flags().BoolP(optDecrypt, "d", false, "decrypt the value")
}

func decrypt(text, secret []byte) ([]byte, error) {
	key, err := pbdk.DeriveKey(secret)
	if err != nil {
		return nil, err
	}

	var l int
	data := make([]byte, base64.RawStdEncoding.DecodedLen(len(text)))
	if l, err = base64.RawStdEncoding.Decode(data, text); err != nil {
		return nil, err
	}

	return aes.GcmDecrypt(data[:l], key)
}
