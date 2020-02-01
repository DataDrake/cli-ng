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

package options

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// Parser can be used to read and convert the raw program arguments
type Parser struct {
	flags map[string]Flag
	args  []string
}

// parseOption adds an option to our internal table
func (p *Parser) parseOption(option, prefix, kind string) {
	pieces := strings.Split(option, "=")
	if len(pieces) == 1 {
		p.flags[strings.TrimPrefix(pieces[0], prefix)] = Flag{kind, ""}
	} else {
		p.flags[strings.TrimPrefix(pieces[0], prefix)] = Flag{kind, pieces[1]}
	}
}

// NewParser does the initial parsing of arguments and returns the resulting Parser
func NewParser(raw []string) (p *Parser, sub string) {
	// Check for subcommand
	if len(raw) < 1 {
		panic("Must use a subcommand")
	}
	// Init parser
	p = &Parser{make(map[string]Flag), make([]string, 0)}
	// Parse options
	for _, curr := range raw {
		switch {
		case strings.HasPrefix(curr, "--"):
			// Parse long option
			p.parseOption(curr, "--", Long)
		case strings.HasPrefix(curr, "-"):
			// Parse short option
			p.parseOption(curr, "-", Short)
		default:
			// Get subcommand
			if sub == "" {
				sub = curr
			} else {
				// get arguments
				p.args = append(p.args, curr)
			}
		}
	}
	return
}

// setField set a StructField to a value
func setField(field reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.Bool:
		// Bools
		field.SetBool(true)
	case reflect.String:
		//String
		field.SetString(value)
	case reflect.Int64:
		// Int64
		i, e := strconv.ParseInt(value, 10, 64)
		if e != nil {
			return e
		}
		field.SetInt(i)
	case reflect.Uint64:
		// Uint64
		u, e := strconv.ParseUint(value, 10, 64)
		if e != nil {
			return e
		}
		field.SetUint(u)
	case reflect.Float64:
		// Float64
		f, e := strconv.ParseFloat(value, 64)
		if e != nil {
			return e
		}
		field.SetFloat(f)
	default:
		return fmt.Errorf("[cli-ng] Unsupported field type: %s", field.Kind().String())
	}
	return nil
}

func setSlice(field reflect.Value, value []string) error {
	kind := field.Type().Elem().Kind()
	switch kind {
	case reflect.String:
		field.Set(reflect.ValueOf(value))
	default:
		return fmt.Errorf("[cli-ng] Unsupported arg slice type '%s'", kind.String())
	}
	return nil
}

// SetFlags attempts to set the entries in 'flags', using the previously parsed arguments
func (p *Parser) SetFlags(flags interface{}) {
	// Get the struct element values
	flagsElement := reflect.ValueOf(flags).Elem()
	// Get the struct element types
	flagsType := flagsElement.Type()
	// Iterate over struct fields
	for i := 0; i < flagsElement.NumField(); i++ {
		tags := flagsType.Field(i).Tag
		element := flagsElement.Field(i)
		if !element.CanSet() {
			continue
		}
		var deletion string
		for k, v := range p.flags {
			if k == tags.Get(v.kind) {
				if err := setField(element, v.value); err != nil {
					panic("Failed to parse flag '" + k + "', reason: " + err.Error())
				}
				deletion = k
				break
			}
		}
		// Remove if a match is found (speed, duplication)
		if deletion != "" {
			delete(p.flags, deletion)
		}
	}
}

// UnknownFlags checks for unregistered flags that are set
func (p *Parser) UnknownFlags() {
	if len(p.flags) > 0 {
		for name, flag := range p.flags {
			fmt.Fprintf(os.Stderr, "Unrecognized flag '%s' with argument '%s'\n", name, flag.value)
		}
	}
}

// SetArgs attempts to set the entries in 'args', using the previously parsed arguments
func (p *Parser) SetArgs(args interface{}) bool {
	argsElement := reflect.ValueOf(args).Elem()
	num := argsElement.NumField()
	if num > 0 {
		if arg := argsElement.Field(num - 1); arg.Kind() == reflect.Slice {
			num--
		}
	}
	if len(p.args) < num {
		return false
	}
	for i := 0; i < argsElement.NumField(); i++ {
		arg := argsElement.Field(i)
		if !arg.CanSet() {
			continue
		}
		if arg.Kind() == reflect.Slice {
			if i != (argsElement.NumField() - 1) {
				panic("[cli-ng] arg slice must be the last argument")
			}
			if err := setSlice(arg, p.args[i:]); err != nil {
				panic("Failed to parse arg '" + arg.String() + "', reason: " + err.Error())
			}
		} else {
			if err := setField(arg, p.args[i]); err != nil {
				panic("Failed to parse arg '" + arg.String() + "', reason: " + err.Error())
			}
		}
	}
	return true
}
