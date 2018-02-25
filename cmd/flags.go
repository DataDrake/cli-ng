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
	"reflect"
	"strconv"
)

// PrintFlags writes out the flags in a struct
func PrintFlags(flags interface{}) {
	t := reflect.TypeOf(flags).Elem()
	if t.NumField() > 0 {
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
		//formatLong := "    -%" + strconv.Itoa(maxShort) + "s, %" + strconv.Itoa(maxLong) + "s%" + strconv.Itoa(maxArg) + "s : %s\n"
		//formatShort := "     %" + strconv.Itoa(maxShort+maxLong+4) + "s %" + strconv.Itoa(maxArg) + "s : %s\n"
		for i := 0; i < t.NumField(); i++ {
			short := t.Field(i).Tag.Get("short")
			long := t.Field(i).Tag.Get("long")
			arg := t.Field(i).Tag.Get("arg")
			desc := t.Field(i).Tag.Get("desc")
			format := " : %s\n"
			if maxArg > 0 {
				if arg == "true" {
					format = " arg" + format
				} else {
					format = "    " + format
				}
			}
			if long != "" {
				format = "    %" + strconv.Itoa(maxShort+1) + "s, %" + strconv.Itoa(maxLong+2) + "s" + format
				fmt.Printf(format, "-"+short, "--"+long, desc)
			} else {
				format = "    %" + strconv.Itoa(maxShort+maxLong+5) + "s" + format
				fmt.Printf(format, "-"+short, desc)
			}
		}
		print("\n\n")
	}
}
