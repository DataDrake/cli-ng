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
)

// Hidden fulfills the "hidden" subcommand
var Hidden = CMD{
	Name:   "hidden",
	Alias:  "hi",
	Short:  "Hidden command for testing",
	Hidden: true,
	Args:   &HiddenArgs{},
	Run:    HiddenRun,
}

// HiddenArgs contains the arguments for the "hidden" subcommand
type HiddenArgs struct{}

// HiddenRun prints the usage for the requested command
func HiddenRun(r *RootCMD, c *CMD) {
	// Get the arguments
	// args := c.Args.(*HiddenArgs).Args
	fmt.Println("You didn't see me!!!!!!")
}
