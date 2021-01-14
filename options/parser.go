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
	raw     List
	numArgs int
	maxArgs int
	minArgs int
}

// NewParser does the initial parsing of arguments and returns the resulting Parser
func NewParser(args []string, single bool) (p *Parser, sub string) {
	if len(args) < 1 {
		panic("Must use a subcommand")
	}
	// pop subcommand off the front
	sub, args = filepath.Base(args[0]), args[1:]
	p = &Parser{
		raw: NewList(args),
	}
	return
}

// ErrInsufficientArgs indicates that not enough arguments were provided for this subcommand
var ErrInsufficientArgs = errors.New("Missing argument(s)")

// Parse processes arguments and sets flags and subcommand args as needed
func (p *Parser) Parse(rFlags, cFlags, args interface{}) (err error) {
	if rFlags, err = p.verifyFlags(rFlags); err != nil {
		return
	}
	if cFlags, err = p.verifyFlags(cFlags); err != nil {
		return
	}
	if args, err = p.verifyArgs(args); err != nil {
		return
	}
	for !p.raw.IsEmpty() {
		if err = p.parseArg(rFlags, cFlags, args); err != nil {
			return
		}
	}
	if p.numArgs < p.minArgs {
		err = ErrInsufficientArgs
		return
	}
	if !p.raw.IsEmpty() {
		err = ErrTooManyArgs
	}
	return
}

func (p *Parser) verifyFlags(flags interface{}) (out interface{}, err error) {
	if v := reflect.ValueOf(flags); v.IsValid() && !v.IsZero() {
		t := v.Elem().Type()
		for i := 0; i < t.NumField(); i++ {
			if name := t.Field(i).Tag.Get("short"); len(name) > 1 {
				err = fmt.Errorf("short flags must have one character names, found: %s", name)
				return
			}
		}
		out = flags
	}
	return
}

func (p *Parser) verifyArgs(args interface{}) (out interface{}, err error) {
	if v := reflect.ValueOf(args); v.IsValid() && !v.IsZero() {
		p.maxArgs = v.Elem().NumField()
		p.minArgs = p.maxArgs
		if p.maxArgs == 0 {
			return
		}
		if last := v.Elem().Type().Field(p.maxArgs - 1); last.Type.Kind() == reflect.Slice {
			if last.Tag.Get("zero") != "" {
				p.minArgs = p.maxArgs - 1
			}
		}
		out = args
	}
	return
}

// ErrMissingFlagName indicates that no flag name was provided
var ErrMissingFlagName = errors.New("missing flag name")

func (p *Parser) parseArg(rFlags, cFlags, args interface{}) error {
	arg := p.raw.Peek()
	switch {
	case strings.HasPrefix(arg, "--"):
		return p.parseLongFlag(rFlags, cFlags)
	case strings.HasPrefix(arg, "-"):
		return p.parseShortFlags(rFlags, cFlags)
	default:
		return p.setArg(args)
	}
}

func (p *Parser) parseLongFlag(rFlags, cFlags interface{}) error {
	name := strings.TrimPrefix(p.raw.Next(), "--")
	if len(name) == 0 {
		return ErrMissingFlagName
	}
	return p.setAnyFlag(rFlags, cFlags, name, "long", false)
}

func (p *Parser) parseShortFlags(rFlags, cFlags interface{}) error {
	chars := strings.TrimPrefix(p.raw.Next(), "-")
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
}

// setAnyFlag attempts to set an entry in 'flags', using the unparsed arguments
func (p *Parser) setAnyFlag(rFlags, cFlags interface{}, name, tag string, last bool) error {
	if found, err := p.setFlag(rFlags, name, tag, last); found {
		return err
	}
	if found, err := p.setFlag(cFlags, name, tag, last); found {
		return err
	}
	return fmt.Errorf("invalid flag '%s'", name)
}

func (p *Parser) setFlag(flags interface{}, name, tag string, last bool) (found bool, err error) {
	// Try setting root flag
	if flags == nil {
		return
	}
	flagsElement := reflect.ValueOf(flags).Elem()
	flagsType := flagsElement.Type()
	// Iterate over struct fields
	for i := 0; i < flagsType.NumField(); i++ {
		element := flagsElement.Field(i)
		if !element.CanSet() {
			continue
		}
		if tags := flagsType.Field(i).Tag; name == tags.Get(tag) {
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
	if args == nil {
		return ErrTooManyArgs
	}
	argsElement := reflect.ValueOf(args).Elem()
	if p.numArgs >= p.maxArgs {
		arg := argsElement.Field(p.maxArgs - 1)
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
	if err := p.setField(arg, false, false); err != nil {
		return fmt.Errorf("Failed to parse arg '%s', reason: %s", arg, err)
	}
	p.numArgs++
	return nil
}

// ErrMissingValue indicates that a flag does not have an associated value
var ErrMissingValue = errors.New("missing value for field")

// ErrBadGroup indicated that more than one short flag in a group expects an argument
var ErrBadGroup = errors.New("only the last short flag in a group can have an argument")

// ErrSliceFlag indicates that a flag has been given a slice type
var ErrSliceFlag = errors.New("flags cannot be slices")

// ErrBoolArg indicates that an argument is a bool (unsupported)
var ErrBoolArg = errors.New("args cannot be bool values")

// setField set a StructField to a value
func (p *Parser) setField(field reflect.Value, flag, last bool) error {
	switch field.Kind() {
	case reflect.Slice:
		if flag {
			return ErrSliceFlag
		}
		p.appendSlice(field)
	case reflect.Bool:
		if !flag {
			return ErrBoolArg
		}
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
	if p.raw.IsEmpty() {
		return ErrMissingValue
	}
	raw := p.raw.Next()
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
		return fmt.Errorf("unsupported field type: %s", field.Kind().String())
	}
	return nil
}

func (p *Parser) appendSlice(field reflect.Value) {
	value := p.raw.Next()
	field.Set(reflect.Append(field, reflect.ValueOf(value)))
}
