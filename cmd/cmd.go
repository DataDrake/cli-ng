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

// CMD is a type for all commands
type CMD struct {
	Name   string
	Alias  string
	Short  string
	Hidden bool
	Args   interface{}
	Run    func(r *RootCMD, c *CMD)
}

// Usage prints a general usage statement
func Usage(r *RootCMD, c *CMD) {
	// Print the usage line
	fmt.Printf("USAGE: %s [OPTIONS] %s", r.Name, c.Name)
	// Print the argument names
	t := reflect.TypeOf(c.Args).Elem()
	max := 0
	for i := 0; i < t.NumField(); i++ {
		name := t.Field(i).Name
		fmt.Printf(" <%s>", name)
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
	if r.Flags != nil {
		fmt.Printf("GLOBAL FLAGS:\n\n")
		PrintFlags(r.Flags)
	}
}
