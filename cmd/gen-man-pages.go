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
	"io"
	"os"
	"reflect"
	"sort"
	"strings"
)

// GenManPages fulfills the "gen-man-pages" subcommand
var GenManPages = CMD{
	Name:   "gen-man-pages",
	Alias:  "gmp",
	Short:  "Generate man-pages for the root command and each sub-command",
	Args:   &GenManPagesArgs{},
	Run:    GenManPagesRun,
	Hidden: true,
}

// GenFlags prints out Flag structs in man-page format
func GenFlags(man io.Writer, v interface{}) error {
	flags := reflect.ValueOf(v).Elem()
	flagsType := flags.Type()
	for i := 0; i < flags.NumField(); i++ {
		tags := flagsType.Field(i).Tag
		fmt.Fprintln(man, ".TP")
		fmt.Fprint(man, ".BR ")
		var short string
		if short = tags.Get("short"); len(short) > 0 {
			fmt.Fprintf(man, "\\-%s", short)
		}
		if long := tags.Get("long"); len(long) > 0 {
			if len(short) > 0 {
				fmt.Fprintf(man, " \", \" \\-\\-%s", long)
			} else {
				fmt.Fprintf(man, "\\-\\-%s", long)
			}
		}
		if tags.Get("arg") != "" {
			fmt.Fprintln(man, " =\\fIARG\\fR")
		} else {
			fmt.Fprintln(man)
		}
		fmt.Fprintf(man, "%s\n\n", tags.Get("desc"))
	}
	return nil
}

// GenerateSubPage generates a man-page for a single subcommand
func GenerateSubPage(r *RootCMD, name string) error {
	// Open file
	man, err := os.Create(r.Name + "-" + name + ".1")
	if err != nil {
		return err
	}
	defer man.Close()
	sub := r.Subcommands[name]
	// Header
	fmt.Fprintf(man, ".TH %s\\-%s 1\n", r.Name, name)
	// Name
	fmt.Fprintln(man, ".SH NAME")
	fmt.Fprintf(man, "%s \\- %s\n", name, sub.Short)
	// Synopsis
	fmt.Fprintln(man, ".SH SYNOPSIS")
	fmt.Fprintf(man, ".B %s\n", r.Name)
	fmt.Fprint(man, "[\\fIOPTIONS...\\fR]")
	fmt.Fprintf(man, " \\fI%s\\fR", name)
	args := reflect.TypeOf(sub.Args).Elem()
	for i := 0; i < args.NumField(); i++ {
		field := args.Field(i)
		fmt.Fprintf(man, " \\fI%s\\fR ", strings.ToUpper(field.Name))
	}
	fmt.Fprintln(man)
	// Arguments
	if args.NumField() > 0 {
		fmt.Fprintln(man, ".SH ARGUMENTS")
		for i := 0; i < args.NumField(); i++ {
			field := args.Field(i)
			tags := args.Field(i).Tag
			fmt.Fprintln(man, ".TP")
			fmt.Fprintf(man, ".B %s\n", strings.ToUpper(field.Name))
			fmt.Fprintf(man, "%s\n\n", tags.Get("desc"))
		}
	}
	// Sub Flags
	if sub.Flags != nil {
		fmt.Fprintf(man, ".SH %s FLAGS\n", strings.ToUpper(name))
		if err := GenFlags(man, sub.Flags); err != nil {
			return err
		}
	}
	// Global Flags
	fmt.Fprintln(man, ".SH GLOBAL FLAGS")
	return GenFlags(man, r.Flags)
}

// GenerateSubPages generates a man-page for every subcommand
func GenerateSubPages(r *RootCMD) error {
	for name, cmd := range r.Subcommands {
		if cmd.Hidden {
			continue
		}
		GenerateSubPage(r, name)
	}
	return nil
}

// GenerateRootPage generates a man-page for the root command
func GenerateRootPage(r *RootCMD) error {
	// Open file
	man, err := os.Create(r.Name + ".1")
	if err != nil {
		return err
	}
	defer man.Close()
	// Header
	fmt.Fprintf(man, ".TH %s 1\n", r.Name)
	// Name
	fmt.Fprintln(man, ".SH NAME")
	fmt.Fprintf(man, "%s \\- %s\n", r.Name, r.Short)
	// Synopsis
	fmt.Fprintln(man, ".SH SYNOPSIS")
	fmt.Fprintf(man, ".B %s\n", r.Name)
	fmt.Fprintln(man, "[\\fIOPTIONS...\\fR] \\fICMD\\fR [\\fIARGS...\\fR]")
	// Subcommands
	names := make([]string, 0)
	for name, cmd := range r.Subcommands {
		if cmd.Hidden {
			continue
		}
		names = append(names, name)
	}
	sort.Strings(names)
	fmt.Fprintln(man, ".SH COMMANDS")
	for _, name := range names {
		sub := r.Subcommands[name]
		fmt.Fprintln(man, ".TP")
		fmt.Fprintf(man, ".B %s (%s) \n", name, sub.Alias)
		fmt.Fprint(man, sub.Short)
		fmt.Fprintf(man, "\n\nSee \\fI%s\\-%s(1)\\fR for specific usage\n\n", r.Name, name)
	}
	// Global Flags
	fmt.Fprintln(man, ".SH GLOBAL FLAGS")
	return GenFlags(man, r.Flags)
}

// GenManPagesArgs contains the arguments for the "gen-man-pages" subcommand
type GenManPagesArgs struct{}

// GenManPagesRun prints the usage for the requested command
func GenManPagesRun(r *RootCMD, c *CMD) {
	// Get the arguments
	// args := c.Args.(*GenManPagesArgs)
	if err := GenerateRootPage(r); err != nil {
		panic(err.Error())
	}
	if err := GenerateSubPages(r); err != nil {
		panic(err.Error())
	}
}
