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
		max := 0
		for i := 0; i < t.NumField(); i++ {
			if long := t.Field(i).Tag.Get("long"); len(long) > max {
				max = len(long)
			}
		}
		for i := 0; i < t.NumField(); i++ {
			short := t.Field(i).Tag.Get("short")
			long := t.Field(i).Tag.Get("long")
			arg := t.Field(i).Tag.Get("arg")
			desc := t.Field(i).Tag.Get("desc")
			if arg == "true" {
				arg = "arg"
			}
			if long != "" {
				fmt.Printf("    -%s,%"+strconv.Itoa(max+3)+"s %3s : %s\n", short, "--"+long, arg, desc)
			} else {
				fmt.Printf("    -%"+strconv.Itoa(max+8)+"s %3s : %s\n", short, arg, desc)
			}
		}
		print("\n\n")
	}
}
