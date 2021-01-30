# TODO

 - [ ] PrintFlags doesn't print string from string types

# BACKLOG

 - [ ] Add Version field to `cmd.Root`
 - [ ] Allow slice args to contain things other than strings
 - [ ] Consider allowing slices for flags, delimited by application-specified delimiter
 - [ ] Add copyright notice for man pages
 - [ ] Add version command which also prints copyright notice

# COMPLETED

 - [x] Re-add `nil` checks for flags and args
 - [x] GenManPages panics on nil `interface{}` value
 - [x] PrintFlags should handle flags without 'short' names
 - [x] Missing newline after usage line in GenManPages when a command has no args
 - [x] [Options...] is printed for sub-commands without options in GenManPages
 - [x] Switch to `tab/writer` for printing sub-commands, args, flags
 - [x] Rework flag printing to automatically add types for non-bools
 - [x] Rework manpage generation to be more DRY
 - [x] Allow multiple short flags in a row (e.g. -tvf)
 - [x] Enfore single-character 'short' names
 - [x] add flag to suppress man page for a sub-command
 - [x] Add sub-command for creating cymlinks for Single binaries
 - [x] Allow empty slice arguments (sero struct tag)
 - [x] GenManPages should print type for flags
 - [x] PrintFlags should print type for flags

