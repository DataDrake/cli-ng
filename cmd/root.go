//
// Copyright 2017-2018 Bryan T. Meyers <bmeyers@datadrake.com>
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
	"github.com/DataDrake/cli-ng/options"
	"github.com/DataDrake/cli-ng/translate"
	"os"
	"sort"
	"strconv"
)

// RootCMD is the main command that runs everything
type RootCMD struct {
	Name        string
	Short       string
	Subcommands map[string]*CMD
	Aliases     map[string]string
	Flags       interface{}
}

// RegisterCMD add a sub-command to this program
func (r *RootCMD) RegisterCMD(c *CMD) {
	// Set up the command
	if r.Subcommands == nil {
		r.Subcommands = make(map[string]*CMD)
	}
	r.Subcommands[c.Name] = c

	// Set up the alias
	if r.Aliases == nil {
		r.Aliases = make(map[string]string)
	}
	r.Aliases[c.Alias] = c.Name
}

// Usage prints the usage for this program
func (r *RootCMD) Usage() {
	translate.Printf("ROOT USAGE", r.Name)
	print("\n\n")
	if len(r.Short) > 0 {
		translate.Printf("DESCRIPTION", r.Short)
		print("\n\n")
	}
	translate.Printf("COMMANDS")
	print("\n\n")
	var keys []string
	maxKey := 0
	maxAlias := 0
	for key, cmd := range r.Subcommands {
		keys = append(keys, key)
		if len(key) > maxKey {
			maxKey = len(key)
		}
		if len(cmd.Alias) > maxAlias {
			maxAlias = len(cmd.Alias)
		}
	}
	maxAlias += 2
	sort.Strings(keys)
	format := "    %" + strconv.Itoa(maxKey) + "s %" + strconv.Itoa(maxAlias) + "s : %s\n"
	for _, k := range keys {
		fmt.Printf(format, k, "("+r.Subcommands[k].Alias+")", r.Subcommands[k].Short)
	}
	print("\n")
	if r.Flags != nil {
		translate.Printf("GLOBAL FLAGS")
		print("\n\n")
		PrintFlags(r.Flags)
	}
	os.Exit(1)
}

// Run finds the appropriate CMD and executes it, or prints the global Usage
func (r *RootCMD) Run() {
	if len(os.Args) < 2 {
		r.Usage()
		os.Exit(1)
	}
	p, sub := options.NewParser(os.Args[1:])
	if sub == "" {
		r.Usage()
		os.Exit(1)
	}
	// Get the subcommand if it exists
	c := r.Subcommands[sub]
	if c == nil {
		// Try to get the subcommand by alias if it not found as a normal name
		if alias := r.Aliases[sub]; alias != "" {
			c = r.Subcommands[alias]
		}
		if c == nil {
			r.Usage()
			os.Exit(1)
		}
	}
	// Handle any flags for the RootCMD
	if r.Flags != nil {
		p.SetFlags(r.Flags)
	}
	// Not yet supported
	//p.SubFlags(c)
	// Handle the arguments for the subcommand
	if !p.SetArgs(c.Args) {
		Usage(r, c)
		os.Exit(1)
	}
	c.Run(r, c)
}
