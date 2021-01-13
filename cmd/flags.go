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
	"github.com/DataDrake/cli-ng/term"
	"io"
	"os"
	"reflect"
	"text/tabwriter"
)

// PrintFlags writes out the flags in a struct
func PrintFlags(flags interface{}) {
	// Get all the struct elements
	t := reflect.TypeOf(flags).Elem()
	args := hasArgs(t)
	if t.NumField() > 0 {
		tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		if args {
			fmt.Fprintln(tw, term.Bold("    NAME\tARG\tDESCRIPTION"))
		} else {
			fmt.Fprintln(tw, term.Bold("    NAME\tDESCRIPTION"))
		}
		// Iterate over arguments
		for i := 0; i < t.NumField(); i++ {
			printFlag(tw, t.Field(i), args)
		}
		tw.Flush()
		fmt.Println()
	}
}

func printFlag(tw io.Writer, f reflect.StructField, args bool) {
	short := f.Tag.Get("short")
	desc := f.Tag.Get("desc")
	name := "-" + short
	if long := f.Tag.Get("long"); long != "" {
		name += ", --" + long
	}
	if args {
		fmt.Fprintf(tw, term.Resetln("    %s\t%s\t%s"), name, arg(f), desc)
	} else {
		fmt.Fprintf(tw, term.Resetln("    %s\t%s"), name, desc)
	}
}

func hasArgs(t reflect.Type) bool {
	for i := 0; i < t.NumField(); i++ {
		if arg(t.Field(i)) != "" {
			return true
		}
	}
	return false
}

func arg(f reflect.StructField) string {
	if f.Tag.Get("arg") != "" {
		return f.Type.Kind().String()
	}
	return ""
}
