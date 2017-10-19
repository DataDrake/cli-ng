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

// NewParser does the initial parsing of arguments and returns the resulting Parser
func NewParser(raw []string) (p *Parser, sub string) {
	p = &Parser{make(map[string]Flag), make([]string, 0)}
	flagsDone := false
	i := 0
	for i < len(raw) {
		switch {
		case flagsDone:
			p.args = append(p.args, raw[i])
			i++
		case strings.HasPrefix(raw[i], "--"):
			pieces := strings.Split(raw[i], "=")
			if len(pieces) == 1 {
				p.flags[strings.TrimPrefix(pieces[0], "--")] = Flag{Long, ""}
			} else {
				p.flags[strings.TrimPrefix(pieces[0], "--")] = Flag{Long, pieces[1]}
			}
			i++
		case strings.HasPrefix(raw[i], "-"):
			pieces := strings.Split(raw[i], "=")
			if len(pieces) == 1 {
				p.flags[strings.TrimPrefix(pieces[0], "-")] = Flag{Short, ""}
			} else {
				p.flags[strings.TrimPrefix(pieces[0], "-")] = Flag{Short, pieces[1]}
			}
			i++
		default:
			flagsDone = true
		}
	}
	if len(p.args) == 0 {
		return
	}
	sub = p.args[0]
	p.args = p.args[1:]
	return
}

// SetFlags attempts to set the entries in 'flags', using the previously parsed arguments
func (p *Parser) SetFlags(flags interface{}) {
	flagsElement := reflect.ValueOf(flags).Elem()
	flagsType := flagsElement.Type()
	for i := 0; i < flagsElement.NumField(); i++ {
		typeField := flagsType.Field(i)
		elementField := flagsElement.Field(i)
		var deletion string
		for k, v := range p.flags {
			if k == typeField.Tag.Get(v.kind) && elementField.CanSet() {
				var err error
				switch elementField.Kind() {
				case reflect.Bool:
					elementField.SetBool(true)
				case reflect.String:
					elementField.SetString(v.value)
				case reflect.Int64:
					i, err := strconv.ParseInt(v.value, 10, 64)
					if err != nil {
						elementField.SetInt(i)
					}
				case reflect.Uint64:
					u, err := strconv.ParseUint(v.value, 10, 64)
					if err != nil {
						elementField.SetUint(u)
					}
				case reflect.Float64:
					f, err := strconv.ParseFloat(v.value, 64)
					if err != nil {
						elementField.SetFloat(f)
					}
				default:
					panic("[cli-ng] Unsupported flag type: " + elementField.Kind().String())
				}
				if err != nil {
					panic("Failed to parse flag '" + k + "', reason: " + err.Error())
				}
				deletion = k
				break
			}
		}
		if deletion != "" {
			delete(p.flags, deletion)
		}
	}
}

// SetArgs attempts to set the entries in 'args', using the previously parsed arguments
func (p *Parser) SetArgs(args interface{}) bool {
	argsElement := reflect.ValueOf(args).Elem()
	if len(p.flags) > 0 {
		for name, flag := range p.flags {
			fmt.Fprintf(os.Stderr, "Unrecognized flag '%s' with argument '%s'\n", name, flag.value)
		}
		return false
	}
	if len(p.args) != argsElement.NumField() {
		return false
	}
	for i := 0; i < argsElement.NumField(); i++ {
		elementField := argsElement.Field(i)
		for _, v := range p.args {
			if elementField.CanSet() {
				var err error
				switch elementField.Kind() {
				case reflect.Bool:
					elementField.SetBool(true)
				case reflect.String:
					elementField.SetString(v)
				case reflect.Int64:
					i, err := strconv.ParseInt(v, 10, 64)
					if err != nil {
						elementField.SetInt(i)
					}
				case reflect.Uint64:
					u, err := strconv.ParseUint(v, 10, 64)
					if err != nil {
						elementField.SetUint(u)
					}
				case reflect.Float64:
					f, err := strconv.ParseFloat(v, 64)
					if err != nil {
						elementField.SetFloat(f)
					}
				default:
					panic("[cli-ng] Unsupported arg type: " + elementField.Kind().String())
				}
				if err != nil {
					panic("Failed to parse arg '" + elementField.String() + "', reason: " + err.Error())
				}
				break
			} else {
				panic("[cli-ng] arg '" + elementField.String() + "' must be public")
			}
		}
	}
	return true
}
