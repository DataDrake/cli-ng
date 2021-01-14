//
// Copyright 2017-2021 Bryan T. Meyers <root@datadrake.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// GenSingleLinks fulfills the "gen-single-links" subcommand
var GenSingleLinks = Sub{
	Name:    "gen-single-links",
	Short:   "GenSingleLinks creates symlinks for a Single binary",
	Hidden:  true,
	SkipMan: true,
	Args:    &GenSingleLinksArgs{},
	Run:     GenSingleLinksRun,
}

// GenSingleLinksArgs specifies the output path for the links
type GenSingleLinksArgs struct {
	Path string `desc:"output dir for symlinks"`
}

// GenSingleLinksRun prints the usage for the requested command
func GenSingleLinksRun(r *Root, c *Sub) {
	args := c.Args.(*GenSingleLinksArgs)
	var keys []string
	for k := range subcommands {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {
		path := filepath.Join(args.Path, key)
		if err := os.Symlink(r.Name, path); err != nil {
			if !os.IsExist(err) {
				fmt.Printf("Failed to make symlink, reason: %s\n", err)
				os.Exit(1)
			}
			fmt.Printf("Warning: '%s' already exists.\n", path)
		}
	}
}
