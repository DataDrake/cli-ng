# cli-ng
An easy to use CLI library for the Go language

[![Go Report Card](https://goreportcard.com/badge/github.com/DataDrake/cli-ng)](https://goreportcard.com/report/github.com/DataDrake/cli-ng)
[![license](https://img.shields.io/github/license/DataDrake/cli-ng.svg)]()
[![Documentation](https://godoc.org/github.com/DataDrake/cli-ng?status.svg)](http://godoc.org/github.com/DataDrake/cli-ng)

## Motivation

`cli-ng` or "CLI Next Generation" is a library intended to allow developers to quickly and concisely build powerful CLI programs. It implements a fully-custom options parser to allow it to easily handle many sub-commands. It may also be used in Single Binary mode to allow a single executable to provide several programs.

## Usage

### Root Command

At the highest level, a `cmd.Root` specifies an `cli-ng` executable. You can assign it all the fields shown here and more. Running the executable is as simple as calling `Root.Run()`.

``` Go
var Root = &cmd.Root {
    Name:      "example",
    Short:     "An example CLI to show off cli-ng",
    Version:   "1.0.0",
    Copyright: "© 2017-2021 Bryan T. Meyers <root@datadrake.com>",
    License:   "Licensed under the Apache License, Version 2.0",
}

func main() {
    Root.Run()
}
```

### Sub Command

A `cli-ng` executable is composed of one or more sub-commands, each with their own purpose and Run function. Sub-commands also support aliases for less typing. Adding a sub-command is a simple as registering it with `cmd.Register()` during initialization. Your `cmd.Root` will automatically know all about it. If `cmd.Root` is run without specifying a sub-command, a Usage message is printed with a listing of all sub-commands and any global flags.

``` Go

func init() {
    cmd.Register(&Sub1)
}

var Sub1 = cmd.Sub {
    Name:  "sub1",
    Alias: "s1",
    Short: "An example sub-command for showing off cli-ng",
    Run: Sub1Run,
}

func Sub1Run(r *cmd.Root, s *cmd.Sub) {
    // do stuff
}

```

### Flags

Flags can be specified both for `cmd.Root` and `cmd.Sub` using nothing but a struct and some tags. The `short` tag specifies a single-character switch for the flag (e.g. -v). The `long` tag specifies a multi-character name for a flag (e.g. --verbose). Flags must specify at least one of the `short` or `long` tags, but both are not required. The `desc` tag provides a short description for the flag. Boolean flags do not accept an argument. If they are specified, the flag is set to `true`. Other types of flags that supported include:

- string
- float32, float64
- int, int8, int16, int32, int64
- uint, uint8, uint16, uint32, uint64

These types of flags can be set by specifying an additional argument, without an `=` (e.g. -v 8).

Structs **MUST** be assigned by pointer.

``` Go
type GlobalFlags struct {
    Danger  bool  `long:"--yes-i-am-really-sure-this-is-what-i-want" desc:"All safeties are off"`
    Debug   bool  `short:"D" desc:"Enable debug logging (e.g. -v 8)"`
    Verbose uint8 `short:"v" long:"verbose" desc:"Enable verbose logging"`
}

var Root = &cmd.Root {
    Name:      "example",
    Short:     "An example CLI to show off cli-ng",
    Version:   "1.0.0",
    Copyright: "© 2017-2021 Bryan T. Meyers <root@datadrake.com>",
    License:   "Licensed under the Apache License, Version 2.0",
    Flags: &GlobalFlags{},
}
```

### Arguments

Arguments can be specified for each `cmd.Sub` by using nothing but a struct and some tags. The `desc` tag provides a short description of the argument. The following types of arguments are supported:

- string
- float32, float64
- int, int8, int16, int32, int64
- uint, uint8, uint16, uint32, uint64

The last value in the struct may also be a Slice of any of these types. By default, this slice must contain at least one element. Setting the `zero` struct tag allows the slice to be empty.

Structs **MUST** be assigned by pointer.

``` Go
type Sub1Args struct {
    Arg1 string  `desc:"First argument"`
    Args []uint8 `zero:"true" desc:"Zero or more byte-sized integers"`
}

var Sub1 = cmd.Sub {
    Name:  "sub1",
    Alias: "s1",
    Short: "An example sub-command for showing off cli-ng",
    Args: &Sub1Args{},
    Run: Sub1Run,
}

```

## Single Binary Mode

Often, Go binaries can be quite large compared to some languages. There are many good reasons for this that are the subject of discussion for some other time. `cli-ng`, however, is able to help with this problem by supporting something called Single Binary mode. In Single Binary mode, `cli-ng` can act as one or more different executables, each with their own flags and arguments. This approach is similar to the one used by `busybox`.

When `cmd.Root.Single` is set to `true`, sub-commands are now treated as executable names. In other words, the name provided by `os.Args[0]` will be expected to match a sub-command. This is easily achieved by symlinking to the executable with links of those names. `cmd.GenSingleLinks` provides this functionality as a sub-command for ease of use. In this mode, `cmd.GenManPages` will also behave accordingly, generating docs for the sub-commands without the name of `cmd.Root` as a prefix.

## Built-in Sub-Commands

The following sub-commands are provided as a part of `cli-ng`'s `cmd` package. The may be optionally registered before running the `cmd.Root`.

### cmd.Help

This is by far the most useful sub-command. Adding it to your program allows users to type `help` or `?` and then the name or alias of a `cmd.Sub` for a Usage message.

### cmd.GenManPages

The `gen-man-pages` sub-command can be used to generate man pages for each of the sub-commands. It is `Hidden` (`cmd.Sub.Hidden` is set to `true`) and will not show up in any Usage messages unless called by `help`. Like `help`, some sub-commands may wish to omit their man page. This can be easily achieved by setting `cmd.Sub.SkipMan` to `true`.

### cmd.GenSingleLinks

The `gen-single-links` sub-command can be used to generate symlinks for each of the sub-commands when running in `Single` mode. It is `Hidden` (`cmd.Sub.Hidden` is set to `true`) and will not show up in any Usage messages unless called by `help`. It also will not show up as a man page. `gen-single-links` accepts a single argument for the directory to install the links to. It expects that the single-binary is installed there as well.

### cmd.Version

THe `version` sub-command prints out the name and version of the program, followed by optional copyright and license notices.

## Projects Using cli-ng

- github.com/alecbcs/lookout/cli
- github.com/arkenproject/ait
- github.com/DataDrake/cuppa
- github.com/DataDrake/go-base (Single Binary)
- github.com/DataDrake/proc-maps
- github.com/DataDrake/static-cling
- github.com/DataDrake/todo
- github.com/EbonJaeger/mcsmanager
- github.com/getsolus/ferryd
- github.com/getsolus/usysconf

## License
Copyright 2017-2021 Bryan T. Meyers <root@datadrake.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
