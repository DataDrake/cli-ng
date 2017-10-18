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
	i := -1
	for k := range r.Subcommands {
		keys = append(keys, k)
		if len(k) > i {
			i = len(k)
		}
	}
	sort.Strings(keys)
	i += 4
	for _, k := range keys {
		fmt.Printf("%"+strconv.Itoa(i)+"s (%s) : %s\n", k, r.Subcommands[k].Alias, r.Subcommands[k].Short)
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
	// Get the subcommand if it exists
	c := r.Subcommands[os.Args[1]]
	if c == nil {
		// Try to get the subcommand by alias if it not found as a normal name
		if alias := r.Aliases[os.Args[1]]; alias != "" {
			c = r.Subcommands[alias]
		}
		if c == nil {
			r.Usage()
			os.Exit(1)
		}
	}
	if len(os.Args) > 2 {
		raw := make([]string, 0)
		copy(raw, os.Args[2:])
		p := options.NewParser(raw)
		// Handle any flags for the RootCMD
		p.SetFlags(&r.Flags)
		// Not yet supported
		//p.SubFlags(c)
		// Handle the arguments for the subcommand
		if !p.SetArgs(c) {
            Usage(r,c)
            os.Exit(1)
        }
	}
	c.Run(r, c)
}
