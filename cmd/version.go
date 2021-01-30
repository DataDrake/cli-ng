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
	"github.com/DataDrake/cli-ng/term"
	"strings"
)

// Version fulfills the "hidden" subcommand
var Version = Sub{
	Name:    "version",
	Short:   "Print the version of this command",
	SkipMan: true,
	Run:     VersionRun,
}

// VersionRun prints the version and optional copyright for a command
func VersionRun(r *Root, c *Sub) {
	fmt.Printf("%s %s\n", r.Name, r.Version)
	if len(r.Copyright) > 0 {
		if strings.Contains(r.Copyright, "\n") {
			fmt.Printf(term.Bold("\nCOPYRIGHT:")+"\n\n%s\n", r.Copyright)
		} else {
			fmt.Printf(term.Bold("\nCOPYRIGHT:")+"\n    %s\n", r.Copyright)
		}
	}
	if len(r.License) > 0 {
		if strings.Contains(r.License, "\n") {
			fmt.Printf(term.Bold("\nLICENSE:")+"\n\n%s\n", r.License)
		} else {
			fmt.Printf(term.Bold("\nLICENSE:")+"\n    %s\n", r.License)
		}
	}
}
