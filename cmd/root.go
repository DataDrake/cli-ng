package cmd

import (
	"flag"
	"fmt"
	"github.com/DataDrake/cli-ng/translate"
	"os"
	"reflect"
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
	flag.Parse()
	if flag.NArg() == 0 {
		r.Usage()
		os.Exit(1)
	}
	c := r.Subcommands[os.Args[1]]
	if c == nil {
		if alias := r.Aliases[os.Args[1]]; alias != "" {
			c = r.Subcommands[alias]
		}
	}
	if c == nil {
		r.Usage()
		os.Exit(1)
	}
	if flag.NArg() != (reflect.TypeOf(c.Args).NumField() + 1) {
		Usage(r, c)
		os.Exit(1)
	}
	c.Run(r, c)
}
