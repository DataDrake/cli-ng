package translate

import (
	"github.com/leonelquinteros/gotext"
)

var enUS = `
msgid "ROOT USAGE"
msgstr "USAGE: %s CMD [OPTIONS]"

msgid "SUBCOMMAND USAGE"
msgstr "USAGE: %s %s [OPTIONS]"

msgid "DESCRIPTION"
msgstr "DESCRIPTION: %s"

msgid "COMMANDS"
msgstr "COMMANDS:"

msgid "ARGUMENTS"
msgstr "ARGUMENTS:"

msgid "GLOBAL FLAGS"
msgstr "GLOBAL FLAGS:"
`

func init() {
	po := new(gotext.Po)
	po.Parse([]byte(enUS))
	Internal["en_US"] = po
}
