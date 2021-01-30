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

package main

import (
	"fmt"
	"github.com/DataDrake/cli-ng/cmd"
)

type level uint8

const license = `Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.`

func main() {

	// Global Flags
	flags := struct {
		Debug   bool  `short:"d" long:"debug" desc:"Show debugging information"`
		NoColor bool  `short:"N" long:"no-color" desc:"Disable coloring of output text"`
		Yes     bool  `short:"y" desc:"Assume yes in all yes/no queries"`
		Verbose bool  `short:"v" long:"verbose" desc:"Detailed output"`
		Level   level `short:"l" arg:"true" long:"level" desc:"Level of something"`
	}{}

	// Build Application
	r := &cmd.Root{
		Name:  "cli-ng",
		Short: "An easy to use CLI library for the Go language",
		Flags: &flags,
        Version: "2.0.0",
        Copyright: "© 2017-2021 Bryan T. Meyers <root@datadrake.com>",
        License: license,
	}

	// Setup the Sub-Commands
	cmd.Register(&cmd.Help)
	cmd.Register(&cmd.Example)
	cmd.Register(&cmd.Hidden)
	cmd.Register(&cmd.GenManPages)
	cmd.Register(&cmd.Version)

	// Run the program
	r.Run()
	if flags.Debug {
		fmt.Println("Debug is on!")
	}
	fmt.Printf("Level is %d\n", flags.Level)
}
