//
// Copyright 2017 Bryan T. Meyers <bmeyers@datadrake.com>
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

// Help fulfills the "help" subcommand
var Help = CMD{
	Name:  "help",
	Alias: "?",
	Short: "Get help with a specific subcommand",
	Args:  &HelpArgs{},
	Run:   HelpRun,
}

// HelpArgs contains the arguments for the "help" subcommand
type HelpArgs struct {
	Subcommand string `desc:"Command to get help for"`
}

// HelpRun prints the usage for the requested command
func HelpRun(r *RootCMD, c *CMD) {

	args := c.Args.(*HelpArgs)

	sub := r.Subcommands[args.Subcommand]
	if sub == nil {
		alias := r.Aliases[args.Subcommand]
		if alias != "" {
			sub = r.Subcommands[alias]
		}
	}
	if sub == nil {
		fmt.Printf("ERROR: '%s' is not a valid subcommand\n", args.Subcommand)
		Usage(r, c)
		os.Exit(1)
	}
	Usage(r, sub)

}
