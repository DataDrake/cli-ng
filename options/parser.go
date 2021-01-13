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

// Parse processes arguments and sets flags and subcommand args as needed
func (p *Parser) Parse(rFlags, cFlags, args interface{}) error {
	var err error
	if v := reflect.ValueOf(args); v.IsValid() && !v.IsZero() {
		p.maxArgs = v.Elem().NumField()
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

// ErrMissingFlagName indicates that no flag name was provided
var ErrMissingFlagName = errors.New("missing flag name")

func (p *Parser) parseArg(rFlags, cFlags, args interface{}) error {
	arg := p.raw[0]
	switch {
	case strings.HasPrefix(arg, "--"):
		p.raw = p.raw[1:]
		name := strings.TrimPrefix(arg, "--")
		if len(name) == 0 {
			return ErrMissingFlagName
		}
		return p.setAnyFlag(rFlags, cFlags, name, "long", false)
	case strings.HasPrefix(arg, "-"):
		p.raw = p.raw[1:]
		chars := strings.TrimPrefix(arg, "-")
		if len(chars) == 0 {
			return ErrMissingFlagName
		}
		for i, char := range chars {
			last := i == (len(chars) - 1)
			if err := p.setAnyFlag(rFlags, cFlags, string(char), "short", last); err != nil {
				return err
			}
		}
		return nil
	default:
		return p.setArg(args)
	}
}

// setAnyFlag attempts to set an entry in 'flags', using the unparsed arguments
func (p *Parser) setAnyFlag(rFlags, cFlags interface{}, name, tag string, last bool) error {
	found, err := p.setFlag(rFlags, name, tag, last)
	if found {
		return err
	}
	found, err = p.setFlag(cFlags, name, tag, last)
	if found {
		return err
	}
	return fmt.Errorf("invalid flag '%s'", name)
}

func (p *Parser) setFlag(flags interface{}, name, tag string, last bool) (found bool, err error) {
	// Try setting root flag
	if v := reflect.ValueOf(flags); !v.IsValid() || v.IsZero() {
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
			found = true
			err = p.setField(element, true, last)
			return
		}
	}
	return
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
	if err := p.setField(arg, false, false); err != nil {
		return fmt.Errorf("Failed to parse arg '%s', reason: %s", arg, err)
	}
	return nil
}

// ErrMissingValue indicates that a flag does not have an associated value
var ErrMissingValue = errors.New("missing value for field")

// ErrBadGroup indicated that more than one short flag in a group expects an argument
var ErrBadGroup = errors.New("only the last short flag in a group can have an argument")

// ErrSliceFlag indicates that a flag has been given a slice type
var ErrSliceFlag = errors.New("flags cannot be slices")

// setField set a StructField to a value
func (p *Parser) setField(field reflect.Value, flag, last bool) error {
	switch field.Kind() {
	case reflect.Slice:
		if flag {
			return ErrSliceFlag
		}
		p.appendSlice(field)
	case reflect.Bool:
		field.SetBool(true)
	default:
		if flag && !last {
			return ErrBadGroup
		}
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
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		value, e := strconv.ParseInt(raw, 10, 64)
		if e != nil {
			return fmt.Errorf("'%s' is not a valid int", raw)
		}
		field.SetInt(value)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		value, e := strconv.ParseUint(raw, 10, 64)
		if e != nil {
			return fmt.Errorf("'%s' is not a valid uint", raw)
		}
		field.SetUint(value)
	case reflect.Float32, reflect.Float64:
		value, e := strconv.ParseFloat(raw, 64)
		if e != nil {
			return fmt.Errorf("'%s' is not a valid float", raw)
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
