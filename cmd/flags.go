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
	"reflect"
)

// PrintFlags writes out the flags in a struct
func PrintFlags(flags interface{}) {
	// Get all the struct elements
	t := reflect.TypeOf(flags).Elem()
	if t.NumField() > 0 {
		// Find all the string lengths
		maxArg := 0
		maxLong := 0
		maxShort := 0
		for i := 0; i < t.NumField(); i++ {
			if t.Field(i).Tag.Get("arg") != "" {
				maxArg = 4
			}
			if long := len(t.Field(i).Tag.Get("long")); long > maxLong {
				maxLong = long
			}
			if short := len(t.Field(i).Tag.Get("short")); short > maxShort {
				maxShort = short
			}
		}
		// Generate format strings
		formatLong := fmt.Sprintf("    %%-%ds, %%-%ds %%3s : %%s\n", maxShort+1, maxLong+2)
		formatShort := fmt.Sprintf("    %%-%ds %%3s : %%s\n", maxShort+maxLong+5)
		// Iterate over arguments
		for i := 0; i < t.NumField(); i++ {
			short := t.Field(i).Tag.Get("short")
			long := t.Field(i).Tag.Get("long")
			arg := t.Field(i).Tag.Get("arg")
			desc := t.Field(i).Tag.Get("desc")
			format := ""
			if maxArg > 0 && arg == "true" {
				format = "arg"
			}
			if long != "" {
				fmt.Printf(formatLong, "-"+short, "--"+long, format, desc)
			} else {
				fmt.Printf(formatShort, "-"+short, format, desc)
			}
		}
		print("\n\n")
	}
}
