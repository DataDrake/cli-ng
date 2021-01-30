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
	"os"
)

// Example fulfills the "example" subcommand
var Example = Sub{
	Name:  "example",
	Alias: "ex",
	Short: "Example command for testing",
	Flags: &ExampleFlags{},
	Args:  &ExampleArgs{},
	Run:   ExampleRun,
}

// ExampleFlags contains the additional flags for the "example" subcommand
type ExampleFlags struct {
	Boop bool `short:"b" long:"boop" desc:"You saw nothing"`
	Bop  bool `long:"bop" desc:"Ouch!"`
}

// ExampleArgs contains the arguments for the "example" subcommand
type ExampleArgs struct {
	Args []uint8 `zero:"yes" desc:"Slice o' Args"`
}

// ExampleRun prints the usage for the requested command
func ExampleRun(r *Root, c *Sub) {
	// Get the arguments
	args := c.Args.(*ExampleArgs).Args
	flags := c.Flags.(*ExampleFlags)
	if flags.Boop {
		fmt.Println("You got booped!!!")
	}
	if flags.Bop {
		fmt.Println("Stop hitting yourself!!!")
	}
	if len(args) == 0 {
		fmt.Println("You get nothing!!!")
		os.Exit(1)
	}
	for _, arg := range args {
		fmt.Printf("You get a '%d'!!!!!!\n", arg)
	}
}
