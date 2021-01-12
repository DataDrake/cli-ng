# TODO

## Fixes
 - [ ] Right align subcommands in cmd.Root.Usage()
 - [ ] Rework manpage generation to be more DRY, maybe use `text/template`
 - [ ] Rework flag printing to automatically add "arg" for non-boolean types.
 - [ ] Allow multiple "short" flags to be specified in a row:
   
   ```
   -xvf vs -x -v -f
   ```
   - This will require enforcing single-character "short" names
 - [x] Re-add `nil` checks for Flags and Args.
 - [ ] Consider allowing an empty slice argument
 - [ ] Add flag to suppress man page for a sub-command
