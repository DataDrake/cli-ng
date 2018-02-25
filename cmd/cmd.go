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
	"github.com/DataDrake/cli-ng/translate"
	"reflect"
	"strconv"
)

// CMD is a type for all commands
type CMD struct {
	Name  string
	Alias string
	Short string
	Args  interface{}
	Run   func(r *RootCMD, c *CMD)
}

// Usage prints a general usage statement
func Usage(r *RootCMD, c *CMD) {
	translate.Printf("SUBCOMMAND USAGE", r.Name, c.Name)
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
	translate.Printf("DESCRIPTION", c.Short)
	print("\n\n")
	if t.NumField() > 0 {
		translate.Printf("ARGUMENTS")
		print("\n\n")
		for i := 0; i < t.NumField(); i++ {
			fmt.Printf("%"+strconv.Itoa(max+4)+"s : %s\n", t.Field(i).Name, t.Field(i).Tag.Get("desc"))
		}
		print("\n")
	}
	if r.Flags != nil {
		translate.Printf("GLOBAL FLAGS")
		print("\n\n")
		PrintFlags(r.Flags)
	}
}
