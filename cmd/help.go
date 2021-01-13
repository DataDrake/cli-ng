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

// Help fulfills the "help" subcommand
var Help = Sub{
	Name:    "help",
	Alias:   "?",
	Short:   "Get help with a specific subcommand",
	SkipMan: true,
	Args:    &HelpArgs{},
	Run:     HelpRun,
}

// HelpArgs contains the arguments for the "help" subcommand
type HelpArgs struct {
	Subcommand string `desc:"Command to get help for"`
}

// HelpRun prints the usage for the requested command
func HelpRun(r *Root, c *Sub) {
	// Get the arguments
	args := c.Args.(*HelpArgs)
	// Find the subcommand
	sub := subcommands[args.Subcommand]
	if sub == nil {
		// Find the aliased subcommand
		alias := aliases[args.Subcommand]
		if alias != "" {
			sub = subcommands[alias]
		}
	}
	// Fail if no matches
	if sub == nil {
		fmt.Printf("ERROR: '%s' is not a valid subcommand\n", args.Subcommand)
		r.Usage()
		os.Exit(1)
	}
	// Print usage
	r.SubUsage(sub)
}
