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
	"errors"
	"fmt"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
)

// Parser can be used to read and convert the raw program arguments
type Parser struct {
	raw     []string
	numArgs int
	maxArgs int
}

// NewParser does the initial parsing of arguments and returns the resulting Parser
func NewParser(raw []string, single bool) (p *Parser, sub string) {
	if len(raw) < 1 {
		panic("Must use a subcommand")
	}
	// pop subcommand off the front
	sub, raw = filepath.Base(raw[0]), raw[1:]
	p = &Parser{
		raw: raw,
	}
	return
}

// ErrInsufficientArgs indicates that not enough arguments were provided for this subcommand
var ErrInsufficientArgs = errors.New("Missing argument(s)")

func (p *Parser) parseArg(rFlags, cFlags, args interface{}) error {
	arg := p.raw[0]
	switch {
	case strings.HasPrefix(arg, "--"):
		return p.setAnyFlag(rFlags, cFlags, "long")
	case strings.HasPrefix(arg, "-"):
		return p.setAnyFlag(rFlags, cFlags, "short")
	default:
		return p.setArg(args)
	}
}

// Parse processes arguments and sets flags and subcommand args as needed
func (p *Parser) Parse(rFlags, cFlags, args interface{}) error {
	var err error
	if args != nil && !reflect.ValueOf(args).IsNil() {
		p.maxArgs = reflect.ValueOf(args).Elem().NumField()
	}
	for len(p.raw) > 0 {
		if err = p.parseArg(rFlags, cFlags, args); err != nil {
			return err
		}
	}
	if p.numArgs < p.maxArgs {
		return ErrInsufficientArgs
	}
	if len(p.raw) > 0 {
		return ErrTooManyArgs
	}
	return err
}

func (p *Parser) setFlag(flags interface{}, tag, name string) (ok bool, err error) {
	// Try setting root flag
	if flags == nil || reflect.ValueOf(flags).IsNil() {
		return
	}
	flagsElement := reflect.ValueOf(flags).Elem()
	flagsType := flagsElement.Type()
	// Iterate over struct fields
	for i := 0; i < flagsElement.NumField(); i++ {
		tags := flagsType.Field(i).Tag
		element := flagsElement.Field(i)
		if !element.CanSet() {
			continue
		}
		if name == tags.Get(tag) {
			p.raw = p.raw[1:]
			if err = p.setField(element); err == nil {
				ok = true
			}
			return
		}
	}
	return
}

// setAnyFlag attempts to set an entry in 'flags', using the unparsed arguments
func (p *Parser) setAnyFlag(rFlags, cFlags interface{}, tag string) error {
	name := strings.TrimLeft(p.raw[0], "-")
	ok, err := p.setFlag(rFlags, tag, name)
	if ok || err != nil {
		return err
	}
	ok, err = p.setFlag(cFlags, tag, name)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("invalid flag '%s'", name)
	}
	return nil
}

// ErrTooManyArgs indicates that too many arguments were provided for this subcommand
var ErrTooManyArgs = errors.New("too many arguments")

// setArg attempts to set an entry in 'args', using an unparsed argument
func (p *Parser) setArg(args interface{}) error {
	argsElement := reflect.ValueOf(args).Elem()
	num := argsElement.NumField()
	if num == 0 {
		return ErrTooManyArgs
	}
	if p.numArgs >= num {
		arg := argsElement.Field(num - 1)
		if arg.Kind() != reflect.Slice {
			return ErrTooManyArgs
		}
		p.numArgs++
		p.appendSlice(arg)
		return nil
	}
	arg := argsElement.Field(p.numArgs)
	if !arg.CanSet() {
		return fmt.Errorf("Failed to set arg '%s', unsettable", arg)
	}
	p.numArgs++
	if err := p.setField(arg); err != nil {
		return fmt.Errorf("Failed to parse arg '%s', reason: %s", arg, err)
	}
	return nil
}

// ErrMissingValue indicates that a flag does not have an associated value
var ErrMissingValue = errors.New("missing value for field")

// setField set a StructField to a value
func (p *Parser) setField(field reflect.Value) error {
	switch field.Kind() {
	case reflect.Slice:
		p.appendSlice(field)
	case reflect.Bool:
		field.SetBool(true)
	default:
		return p.setFieldValue(field)
	}
	return nil
}

func (p *Parser) setFieldValue(field reflect.Value) error {
	if len(p.raw) < 1 {
		return ErrMissingValue
	}
	var raw string
	raw, p.raw = p.raw[0], p.raw[1:]
	if strings.HasPrefix("-", raw) {
		return ErrMissingValue
	}
	switch field.Kind() {
	case reflect.String:
		field.SetString(raw)
	case reflect.Int64:
		value, e := strconv.ParseInt(raw, 10, 64)
		if e != nil {
			return fmt.Errorf("'%s' is not a valid int64", raw)
		}
		field.SetInt(value)
	case reflect.Uint64:
		value, e := strconv.ParseUint(raw, 10, 64)
		if e != nil {
			return fmt.Errorf("'%s' is not a valid uint64", raw)
		}
		field.SetUint(value)
	case reflect.Float64:
		value, e := strconv.ParseFloat(raw, 64)
		if e != nil {
			return fmt.Errorf("'%s' is not a valid float64", raw)
		}
		field.SetFloat(value)
	default:
		return fmt.Errorf("[cli-ng] Unsupported field type: %s", field.Kind().String())
	}
	return nil
}

func (p *Parser) appendSlice(field reflect.Value) {
	var value string
	value, p.raw = p.raw[0], p.raw[1:]
	field.Set(reflect.Append(field, reflect.ValueOf(value)))
}
