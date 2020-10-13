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
	raw []string
}

// NewParser does the initial parsing of arguments and returns the resulting Parser
func NewParser(raw []string, single bool) (p *Parser, sub string) {
	// Check for subcommand
	if len(raw) < 1 {
		panic("Must use a subcommand")
	}
	// get subcommand
	sub, raw = filepath.Base(raw[0]), raw[1:]
	// Init parser
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

	numArgs := 0
	i := 0

	for i < len(p.raw) {
		arg := p.raw[i]
		switch {
		case strings.HasPrefix(arg, "--"):
			name := strings.TrimPrefix(arg, "--")
			i, err = p.setFlag(rFlags, cFlags, i, name, "long")
		case strings.HasPrefix(arg, "-"):
			name := strings.TrimPrefix(arg, "-")
			i, err = p.setFlag(rFlags, cFlags, i, name, "short")
		default:
			err = p.setArg(args, i, numArgs)
			numArgs++
			i++
		}
		if err != nil {
			return err
		}
	}
	if numArgs < reflect.ValueOf(args).Elem().NumField() {
		return ErrInsufficientArgs
	}
	return err
}

// setFlag attempts to set an entry in 'flags', using the unparsed arguments
func (p *Parser) setFlag(rFlags, cFlags interface{}, i int, name, tag string) (idx int, err error) {
	idx = i
	// Try setting root flag
	flagsElement := reflect.ValueOf(rFlags).Elem()
	flagsType := flagsElement.Type()
	// Iterate over struct fields
	for i := 0; i < flagsElement.NumField(); i++ {
		tags := flagsType.Field(i).Tag
		element := flagsElement.Field(i)
		if !element.CanSet() {
			continue
		}
		if name == tags.Get(tag) {
			idx, err = p.setField(element, idx)
			return
		}
	}

	// Try setting subcommand flag
	flagsElement = reflect.ValueOf(cFlags).Elem()
	flagsType = flagsElement.Type()
	// Iterate over struct fields
	for i := 0; i < flagsElement.NumField(); i++ {
		tags := flagsType.Field(i).Tag
		element := flagsElement.Field(i)
		if !element.CanSet() {
			continue
		}
		if name == tags.Get(tag) {
			idx, err = p.setField(element, idx)
			return
		}
	}
	err = fmt.Errorf("flag '%s' not found", name)
	return
}

// ErrTooManyArgs indicates that too many arguments were provided for this subcommand
var ErrTooManyArgs = errors.New("too many arguments")

// setArg attempts to set an entry in 'args', using an unparsed argument
func (p *Parser) setArg(args interface{}, i, numArg int) error {
	argsElement := reflect.ValueOf(args).Elem()
	num := argsElement.NumField()
	if num == 0 {
		return ErrTooManyArgs
	}
	if numArg >= num {
		arg := argsElement.Field(num - 1)
		if arg.Kind() != reflect.Slice {
			return ErrTooManyArgs
		}
		p.appendSlice(arg, p.raw[i])
		return nil
	}
	arg := argsElement.Field(numArg)
	if !arg.CanSet() {
		return fmt.Errorf("Failed to set arg '%s', unsettable", arg)
	}
	var err error
	if arg.Kind() == reflect.Slice {
		p.appendSlice(arg, p.raw[i])
	} else {
		_, err = p.setField(arg, i)
	}
	if err != nil {
		err = fmt.Errorf("Failed to parse arg '%s', reason: %s", arg, err)
	}
	return err
}

// ErrMissingValue indicates that a flag does not have an associated value
var ErrMissingValue = errors.New("missing value for field")

// setField set a StructField to a value
func (p *Parser) setField(field reflect.Value, i int) (idx int, err error) {
	idx = i + 1
	if field.Kind() == reflect.Bool {
		field.SetBool(true)
		return
	}
	if idx >= len(p.raw) {
		err = ErrMissingValue
		return
	}
	next := p.raw[idx]
	if strings.HasPrefix("-", next) {
		err = ErrMissingValue
		return
	}
	idx++
	switch field.Kind() {
	case reflect.String:
		//String
		field.SetString(next)
	case reflect.Int64:
		// Int64
		value, e := strconv.ParseInt(next, 10, 64)
		if e != nil {
			err = fmt.Errorf("'%s' is not a valid int64", next)
			return
		}
		field.SetInt(value)
	case reflect.Uint64:
		// Uint64
		value, e := strconv.ParseUint(next, 10, 64)
		if e != nil {
			err = fmt.Errorf("'%s' is not a valid uint64", next)
			return
		}
		field.SetUint(value)
	case reflect.Float64:
		// Float64
		value, e := strconv.ParseFloat(next, 64)
		if e != nil {
			err = fmt.Errorf("'%s' is not a valid float64", next)
			return
		}
		field.SetFloat(value)
	default:
		err = fmt.Errorf("[cli-ng] Unsupported field type: %s", field.Kind().String())
		return
	}
	return
}

func (p *Parser) appendSlice(field reflect.Value, value string) error {
	kind := field.Type().Elem().Kind()
	switch kind {
	case reflect.String:
		field.Set(reflect.Append(field, reflect.ValueOf(value)))
	default:
		return fmt.Errorf("[cli-ng] Unsupported arg slice type '%s'", kind.String())
	}
	return nil
}
