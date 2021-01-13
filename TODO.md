# TODO

## Fixes
 - [x] Switch to `tab/writer` for printing sub-commands, args, and flags
 - [x] Rework flag printing to automatically add types for non-boolean types.
 - [x] Rework manpage generation to be more DRY
 - [x] Allow multiple "short" flags to be specified in a row:
   ```
   -xvf vs -x -v -f
   ```
 - [x] Enforce single-character "short" names
 - [x] Re-add `nil` checks for Flags and Args.
 - [ ] ~Consider allowing flags before the sub-command is specified~
    - This is a bad idea since it would require setting flags twice to parse a sub-command and separate logic for Single binaries
 - [x] Consider allowing an empty slice argument, requires StructTag (yes)
 - [x] Add flag to suppress man page for a sub-command
