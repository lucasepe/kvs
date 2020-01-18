/*
Copyright © 2019 Luca Sepe <luca.sepe@gmail.com>

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
	"os"
	"path"
	"strings"

	"github.com/avelino/slugify"
	"github.com/lucasepe/kvs/pkg/cl"

	"github.com/spf13/cobra"
)

const (
	appName    = "kvs"
	appVersion = "1.00"
	appSummary = "Key Val Store"
	banner     = ` _
| | __  __   __    ___ 
| |/ /  \ \ / /  / __|
|   <    \ V /   \__ \
|_|\_\ey  \_/ al |___/ tore`

	optBucket = "bucket"
	optStore  = "store"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	DisableSuggestions:    true,
	DisableFlagsInUseLine: true,
	Version:               appVersion,
	Use:                   fmt.Sprintf("%s <COMMAND>", appName),
	Short:                 appSummary,
	Long:                  banner,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// getStoreFile returns the 'slugified' database path.
func getStoreFile() string {
	name, err := rootCmd.Flags().GetString(optStore)
	cl.HandleErr(err)
	name = strings.TrimSpace(name)
	if len(name) == 0 {
		cl.HandleErr(fmt.Errorf("missing store name")) // TODO better error handling
	}

	if !strings.HasSuffix(name, ".db") {
		name = fmt.Sprintf("%s.db", slugify.Slugify(name))
	}

	dir, err := cl.GetWorkDir(appName)
	cl.HandleErr(err)

	return path.Join(dir, name)
}

func init() {
	rootCmd.PersistentFlags().StringP(optStore, "s", "vault", "store name")

	rootCmd.SetVersionTemplate(`{{with .Name}}{{printf "%s " .}}{{end}}{{printf "%s" .Version}} - Luca Sepe <luca.sepe@gmail.com>
`)
}
