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
	"io"
	"os"
	"reflect"
	"sort"
	"strings"
)

// GenManPages fulfills the "gen-man-pages" subcommand
var GenManPages = Sub{
	Name:   "gen-man-pages",
	Alias:  "gmp",
	Short:  "Generate man-pages for the root command and each sub-command",
	Run:    GenManPagesRun,
	Hidden: true,
}

// GenManPagesRun prints the usage for the requested command
func GenManPagesRun(r *Root, c *Sub) {
	if err := GenerateRootPage(r); err != nil {
		panic(err)
	}
	if err := GenerateSubPages(r); err != nil {
		panic(err)
	}
}

// GenerateRootPage generates a man-page for the root command
func GenerateRootPage(r *Root) error {
	// Open file
	man, err := os.Create(r.Name + ".1")
	if err != nil {
		return err
	}
	defer man.Close()
	genRootHeader(man, r)
	genRootSynopsis(man, r)
	genRootSubcommands(man, r)
	// Global Flags
	genFlags(man, r.Flags, "GLOBAL FLAGS")
	return nil
}

func genRootHeader(man io.Writer, r *Root) {
	fmt.Fprintf(man, ".TH %s 1\n", r.Name)
	fmt.Fprintln(man, ".SH NAME")
	fmt.Fprintf(man, "%s \\- %s\n", r.Name, r.Short)
}

func genRootSynopsis(man io.Writer, r *Root) {
	fmt.Fprintln(man, ".SH SYNOPSIS")
	if r.Single {
		fmt.Fprintln(man, "\\fICMD\\fR [\\fIOPTIONS...\\fR] [\\fIARGS...\\fR]")
	} else {
		fmt.Fprintf(man, ".B %s \\fICMD\\fR [\\fIOPTIONS...\\fR] [\\fIARGS...\\fR]\n", r.Name)
	}
}

func genRootSubcommands(man io.Writer, r *Root) {
	names := getVisibleSubcommands(r)
	fmt.Fprintln(man, ".SH COMMANDS")
	for _, name := range names {
		sub := subcommands[name]
		genRootSubcommand(man, r, sub, name)
	}
}

func getVisibleSubcommands(r *Root) (names []string) {
	for name, cmd := range subcommands {
		if cmd.Hidden || cmd.SkipMan {
			continue
		}
		names = append(names, name)
	}
	sort.Strings(names)
	return
}

func genRootSubcommand(man io.Writer, r *Root, sub *Sub, name string) {
	fmt.Fprintln(man, ".TP")
	if r.Single {
		fmt.Fprintf(man, ".B %s \n", name)
		fmt.Fprint(man, sub.Short)
		fmt.Fprintf(man, "\n\nSee \\fI%s(1)\\fR for specific usage\n\n", name)
	} else {
		fmt.Fprintf(man, ".B %s (%s) \n", name, sub.Alias)
		fmt.Fprint(man, sub.Short)
		fmt.Fprintf(man, "\n\nSee \\fI%s\\-%s(1)\\fR for specific usage\n\n", r.Name, name)
	}
}

// GenerateSubPages generates a man-page for every subcommand
func GenerateSubPages(r *Root) error {
	for name, cmd := range subcommands {
		if cmd.Hidden || cmd.SkipMan {
			continue
		}
		if err := GenerateSubPage(r, name); err != nil {
			return err
		}
	}
	return nil
}

// GenerateSubPage generates a man-page for a single subcommand
func GenerateSubPage(r *Root, name string) error {
	// Open file
	var filename string
	if r.Single {
		filename = name + ".1"
	} else {
		filename = r.Name + "-" + name + ".1"
	}
	man, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer man.Close()
	sub := subcommands[name]
	genSubHeader(man, r, sub, name)
	genSubSynopsis(man, r, sub, name)
	genSubArgs(man, sub)
	// Sub Flags
	genFlags(man, sub.Flags, strings.ToUpper(name)+" FLAGS")
	// Global Flags
	genFlags(man, r.Flags, "GLOBAL FLAGS")
	return nil
}

func genSubHeader(man io.Writer, r *Root, sub *Sub, name string) {
	if r.Single {
		fmt.Fprintf(man, ".TH %s 1\n", name)
	} else {
		fmt.Fprintf(man, ".TH %s\\-%s 1\n", r.Name, name)
	}
	// Name
	fmt.Fprintln(man, ".SH NAME")
	fmt.Fprintf(man, "%s \\- %s\n", name, sub.Short)
}

func genSubSynopsis(man io.Writer, r *Root, sub *Sub, name string) {
	fmt.Fprintln(man, ".SH SYNOPSIS")
	if r.Single {
		fmt.Fprintf(man, ".B %s\n", name)
	} else {
		fmt.Fprintf(man, ".B %s \\fI%s\\fR\n", r.Name, name)
	}
	hasFlags := false
	if v := reflect.ValueOf(r.Flags); v.IsValid() && !v.IsZero() {
		hasFlags = true
	}
	if v := reflect.ValueOf(sub.Flags); v.IsValid() && !v.IsZero() {
		hasFlags = true
	}
	if hasFlags {
		fmt.Fprint(man, "[\\fIOPTIONS...\\fR]")
	}
}

func genSubArgs(man io.Writer, sub *Sub) {
	if v := reflect.ValueOf(sub.Args); !v.IsValid() || v.IsZero() {
		fmt.Fprintf(man, "\n\n")
		return
	}
	args := reflect.TypeOf(sub.Args).Elem()
	for i := 0; i < args.NumField(); i++ {
		field := args.Field(i)
		if k := field.Type.Kind(); k == reflect.Slice {
			if field.Tag.Get("zero") != "" {
				fmt.Fprintf(man, " [\\fI%s...\\fR] ", strings.ToUpper(field.Name))
			} else {
				fmt.Fprintf(man, " \\fI%s...\\fR ", strings.ToUpper(field.Name))
			}
		} else {
			fmt.Fprintf(man, " \\fI%s\\fR ", strings.ToUpper(field.Name))
		}
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
}

// genFlags prints out Flag structs in man-page format
func genFlags(man io.Writer, flags interface{}, name string) {
	if v := reflect.ValueOf(flags); !v.IsValid() || v.IsZero() {
		return
	}
	fmt.Fprintf(man, ".SH %s\n", name)
	v := reflect.ValueOf(flags).Elem()
	flagsType := v.Type()
	for i := 0; i < flagsType.NumField(); i++ {
		genFlag(man, flagsType.Field(i))
	}
}

func genFlag(man io.Writer, field reflect.StructField) {
	fmt.Fprintln(man, ".TP")
	fmt.Fprint(man, ".BR ")
	tag := field.Tag
	var short string
	if short = tag.Get("short"); len(short) > 0 {
		fmt.Fprintf(man, "\\-%s", short)
	}
	if long := tag.Get("long"); len(long) > 0 {
		if len(short) > 0 {
			fmt.Fprintf(man, " \", \" \\-\\-%s", long)
		} else {
			fmt.Fprintf(man, "\\-\\-%s", long)
		}
	}
	if k := field.Type.Kind(); k != reflect.Bool {
		fmt.Fprintf(man, " \" \\fI%s\\fR\n", strings.ToUpper(k.String()))
	} else {
		fmt.Fprintln(man, "\\fR")
	}
	fmt.Fprintf(man, "%s\n\n", tag.Get("desc"))
}
