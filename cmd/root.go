//
// Copyright 2017-2020 Bryan T. Meyers <root@datadrake.com>
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
	"os"
	"sort"
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
	// Print Usage
	fmt.Printf("USAGE: %s [OPTIONS] CMD\n\n", r.Name)
	// Print Description
	if len(r.Short) > 0 {
		fmt.Printf("DESCRIPTION: %s\n\n", r.Short)
	}
	// Print sub-commands
	fmt.Printf("COMMANDS:\n\n")
	// Key the names of the sub-commands and find the longest command and alias
	var keys []string
	maxKey := 0
	maxAlias := 0
	for key, cmd := range r.Subcommands {
		if cmd.Hidden {
			continue
		}
		keys = append(keys, key)
		if len(key) > maxKey {
			maxKey = len(key)
		}
		if len(cmd.Alias) > maxAlias {
			maxAlias = len(cmd.Alias)
		}
	}
	sort.Strings(keys)
	// Add spacing for ()
	format := fmt.Sprintf("    %%%ds (%%%ds) : %%s\n", maxKey, maxAlias)
	for _, k := range keys {
		fmt.Printf(format, k, r.Subcommands[k].Alias, r.Subcommands[k].Short)
	}
	print("\n")
	// Print the global flags
	if r.Flags != nil {
		fmt.Printf("GLOBAL FLAGS:\n\n")
		PrintFlags(r.Flags)
	}
	os.Exit(1)
}

// Run finds the appropriate CMD and executes it, or prints the global Usage
func (r *RootCMD) Run() {
	if len(os.Args) < 2 {
		r.Usage()
	}
	p, sub := options.NewParser(os.Args[1:])
	if sub == "" {
		r.Usage()
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
		}
	}
	// Handle any flags for the RootCMD
	if r.Flags != nil {
		p.SetFlags(r.Flags)
	}
	// Not yet supported
	if c.Flags != nil {
		p.SetFlags(c.Flags)
	}
	// Check for unknown flags
	p.UnknownFlags()
	// Handle the arguments for the subcommand
	if !p.SetArgs(c.Args) {
		Usage(r, c)
		os.Exit(1)
	}
	c.Run(r, c)
}
