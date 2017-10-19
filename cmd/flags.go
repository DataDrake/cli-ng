package cmd

import (
	"fmt"
	"reflect"
	"strconv"
)

// PrintFlags writes out the flags in a struct
func PrintFlags(flags interface{}) {
	t := reflect.TypeOf(flags).Elem()
	if t.NumField() > 0 {
		max := 0
		for i := 0; i < t.NumField(); i++ {
			if long := t.Field(i).Tag.Get("long"); len(long) > max {
				max = len(long)
			}
		}
		for i := 0; i < t.NumField(); i++ {
			short := t.Field(i).Tag.Get("short")
			long := t.Field(i).Tag.Get("long")
			arg := t.Field(i).Tag.Get("arg")
			desc := t.Field(i).Tag.Get("desc")
			if arg == "true" {
				arg = "arg"
			}
			if long != "" {
				fmt.Printf("    -%s,%"+strconv.Itoa(max+3)+"s %3s : %s\n", short, "--"+long, arg, desc)
			} else {
				fmt.Printf("    -%"+strconv.Itoa(max+8)+"s %3s : %s\n", short, arg, desc)
			}
		}
		print("\n\n")
	}
}
