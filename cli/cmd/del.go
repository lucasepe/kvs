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
	"fmt"

	"github.com/avelino/slugify"
	"github.com/lucasepe/kvs"
	"github.com/lucasepe/kvs/pkg/cl"
	"github.com/spf13/cobra"
)

// delCmd represents the delete command
var delCmd = &cobra.Command{
	Use:                   "del <key>",
	DisableSuggestions:    true,
	DisableFlagsInUseLine: false,
	Args:                  cobra.MaximumNArgs(1),
	Short:                 "Removes from a store the item with the specified key from a bucket",
	Run: func(cmd *cobra.Command, args []string) {
		bucket, err := cmd.Flags().GetString(optBucket)
		cl.HandleErr(err)
		bucket = slugify.Slugify(bucket)

		store, err := kvs.Open(getStoreFile())
		cl.HandleErr(err)
		defer store.Close()

		if len(args) > 0 {
			key := slugify.Slugify(args[0])
			if err := store.Delete(bucket, key); err == nil {
				fmt.Printf("item with key '%s' successfully removed from bucket '%s'\n", key, bucket)
			}
		} else {
			if err := store.DeleteBucket(bucket); err == nil {
				fmt.Printf("bucket '%s' successfully removed\n", bucket)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(delCmd)

	delCmd.Flags().StringP(optBucket, "b", "", "bucket name")
	delCmd.MarkFlagRequired(optBucket)
}
