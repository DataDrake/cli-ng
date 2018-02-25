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

package main

import (
	"github.com/DataDrake/cli-ng/cmd"
)

func main() {

	// Global Flags
	flags := struct {
		Debug   bool `short:"d" long:"debug" desc:"Show debugging information"`
		NoColor bool `short:"N" long:"no-color" desc:"Disable coloring of output text"`
		Yes     bool `short:"y" long:"yes-all" desc:"assume yes in all yes/no queries"`
		Verbose bool `short:"v" long:"verbose" desc:"Detailed output"`
	}{}

	// Build Application
	r := &cmd.RootCMD{
		Name:  "cli-ng",
		Short: "An easy to use CLI library for the Go language",
		Flags: &flags,
	}

	// Setup the Sub-Commands
	r.RegisterCMD(&cmd.Help)

	// Run the program
	r.Run()
}
