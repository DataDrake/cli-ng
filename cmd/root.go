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
	"github.com/DataDrake/cli-ng/options"
	"github.com/DataDrake/cli-ng/term"
	"os"
	"reflect"
	"sort"
	"strings"
	"text/tabwriter"
)

// Root is the main command that supports multiple Sub commands
type Root struct {
	Name   string
	Short  string
	Flags  interface{}
	Single bool
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

// Usage prints the usage for this program
func (r *Root) Usage() {
	if r.Single {
		fmt.Printf(term.Bold("NAME:")+" %s\n\n", r.Name)
	} else {
		fmt.Printf(term.Bold("USAGE:")+" %s CMD [OPTIONS]\n\n", r.Name)
	}
	if len(r.Short) > 0 {
		fmt.Printf(term.Bold("DESCRIPTION:")+" %s\n\n", r.Short)
	}
	r.printSubcommands()
	if r.Flags != nil {
		fmt.Printf(term.Bold("GLOBAL FLAGS:\n\n"))
		PrintFlags(r.Flags)
	}
	os.Exit(1)
}

func (r *Root) printSubcommands() {
	fmt.Printf(term.Bold("COMMANDS:\n\n"))
	keys := generateKeys()
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	if r.Single {
		fmt.Fprintln(tw, term.Bold("    NAME\tDESCRIPTION"))
		for _, k := range keys {
			fmt.Fprintf(tw, term.Resetln("    %s\t%s"), k, subcommands[k].Short)
		}
	} else {
		fmt.Fprintln(tw, term.Bold("    NAME\tALIAS\tDESCRIPTION"))
		for _, k := range keys {
			fmt.Fprintf(tw, term.Resetln("    %s\t%s\t%s"), k, subcommands[k].Alias, subcommands[k].Short)
		}
	}
	tw.Flush()
	fmt.Println()
}

func generateKeys() (keys []string) {
	for key, cmd := range subcommands {
		if cmd.Hidden {
			continue
		}
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return
}

// SubUsage prints a general usage statement for a subcommand
func (r *Root) SubUsage(c *Sub) {
	// Print the usage line
	if r.Single {
		fmt.Printf(term.Bold("USAGE:")+" %s [OPTIONS]", c.Name)
	} else {
		fmt.Printf(term.Bold("USAGE:")+" %s %s [OPTIONS]", r.Name, c.Name)
	}
	// Print the argument names
	if v := reflect.ValueOf(c.Args); v.IsValid() && !v.IsZero() {
		t := v.Elem().Type()
		for i := 0; i < t.NumField(); i++ {
			name := t.Field(i).Name
			if t.Field(i).Type.Kind() == reflect.Slice {
				fmt.Printf(" [%s1 ... %sN]", name, name)
			} else {
				fmt.Printf(" <%s>", name)
			}
		}
	}
	print("\n\n")
	// Print the description
	fmt.Printf(term.Bold("DESCRIPTION:")+" %s\n\n", c.Short)
	// Print the arguments
	if v := reflect.ValueOf(c.Args); v.IsValid() && !v.IsZero() {
		t := v.Elem().Type()
		if t.NumField() > 0 {
			fmt.Printf(term.Bold("ARGUMENTS:\n\n"))
			tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(tw, term.Bold("    NAME\tDESCRIPTION"))
			for i := 0; i < t.NumField(); i++ {
				fmt.Fprintf(tw, term.Resetln("    %s\t%s"), t.Field(i).Name, t.Field(i).Tag.Get("desc"))
			}
			tw.Flush()
			fmt.Println()
		}
	}
	// Print global flags
	if c.Flags != nil {
		fmt.Printf(term.Bold("%s FLAGS:\n\n"), strings.ToUpper(c.Name))
		PrintFlags(c.Flags)
	}
	// Print global flags
	if r.Flags != nil {
		fmt.Printf(term.Bold("GLOBAL FLAGS:\n\n"))
		PrintFlags(r.Flags)
	}
}
