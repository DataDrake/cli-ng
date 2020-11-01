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
	"reflect"
	"sort"
	"strings"
)

// Root is the main command that supports multiple Sub commands
type Root struct {
	Name   string
	Short  string
	Flags  interface{}
	Single bool
}

func generateKeys() (keys []string, maxKey, maxAlias int) {
	for key, cmd := range subcommands {
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
	return
}

func (r *Root) printSubcommands() {
	fmt.Printf("COMMANDS:\n\n")
	keys, maxKey, maxAlias := generateKeys()
	// Add spacing for ()
	if r.Single {
		format := fmt.Sprintf("    %%%ds : %%s\n", maxAlias)
		for _, k := range keys {
			fmt.Printf(format, k, subcommands[k].Short)
		}
	} else {
		format := fmt.Sprintf("    %%%ds (%%%ds) : %%s\n", maxKey, maxAlias)
		for _, k := range keys {
			fmt.Printf(format, k, subcommands[k].Alias, subcommands[k].Short)
		}
	}
	print("\n")
}

// Usage prints the usage for this program
func (r *Root) Usage() {
	if r.Single {
		fmt.Printf("NAME: %s\n\n", r.Name)
	} else {
		fmt.Printf("USAGE: %s CMD [OPTIONS]\n\n", r.Name)
	}
	if len(r.Short) > 0 {
		fmt.Printf("DESCRIPTION: %s\n\n", r.Short)
	}
	r.printSubcommands()
	if r.Flags != nil {
		fmt.Printf("GLOBAL FLAGS:\n\n")
		PrintFlags(r.Flags)
	}
	os.Exit(1)
}

// SubUsage prints a general usage statement for a subcommand
func (r *Root) SubUsage(c *Sub) {
	// Print the usage line
	if r.Single {
		fmt.Printf("USAGE: %s [OPTIONS]", c.Name)
	} else {
		fmt.Printf("USAGE: %s %s [OPTIONS]", r.Name, c.Name)
	}
	// Print the argument names
	t := reflect.TypeOf(c.Args).Elem()
	max := 0
	for i := 0; i < t.NumField(); i++ {
		name := t.Field(i).Name
		if t.Field(i).Type.Kind() == reflect.Slice {
			fmt.Printf(" [%s1 ... %sN]", name, name)
		} else {
			fmt.Printf(" <%s>", name)
		}
		if len(name) > max {
			max = len(name)
		}
	}
	print("\n\n")
	// Print the description
	fmt.Printf("DESCRIPTION: %s\n\n", c.Short)
	// Print the arguments
	format := fmt.Sprintf("%%%ds : %%s\n", max+4)
	if t.NumField() > 0 {
		fmt.Printf("ARGUMENTS:\n\n")
		for i := 0; i < t.NumField(); i++ {
			fmt.Printf(format, t.Field(i).Name, t.Field(i).Tag.Get("desc"))
		}
		print("\n")
	}
	// Print global flags
	if c.Flags != nil {
		fmt.Printf("%s FLAGS:\n\n", strings.ToUpper(c.Name))
		PrintFlags(c.Flags)
	}
	// Print global flags
	if r.Flags != nil {
		fmt.Printf("GLOBAL FLAGS:\n\n")
		PrintFlags(r.Flags)
	}
}

// Run finds the appropriate CMD and executes it, or prints the global Usage
func (r *Root) Run() {
	if !r.Single && len(os.Args) < 2 {
		r.Usage()
	}
	args := os.Args
	if !r.Single {
		args = args[1:]
	}
	p, sub := options.NewParser(args, r.Single)
	if sub == "" {
		r.Usage()
	}
	// Get the subcommand if it exists
	c := subcommands[sub]
	if c == nil {
		// Try to get the subcommand by alias if it not found as a normal name
		if alias := aliases[sub]; alias != "" {
			c = subcommands[alias]
		}
		if c == nil {
			r.Usage()
		}
	}
	// Parser flags
	if err := p.Parse(r.Flags, c.Flags, c.Args); err != nil {
		fmt.Printf("Error: %s\n\n", err)
		r.SubUsage(c)
		os.Exit(1)
	}
	c.Run(r, c)
}
