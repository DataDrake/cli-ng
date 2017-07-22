package translate

import (
	"github.com/leonelquinteros/gotext"
	"os"
	"strings"
)

// Internal holds the internal locales for this library
var Internal map[string]*gotext.Po

var language string

// GetLanguage retrieves the name of the Locale in use
func GetLanguage() (l string) {
	l = language
	if l != "" {
		goto STRIP
	}
	if l = os.Getenv("LANGUAGE"); l != "" {
		goto STRIP
	}
	if l = os.Getenv("LC_ALL"); l != "" {
		goto STRIP
	}
	if l = os.Getenv("LANG"); l != "" {
		goto STRIP
	}
	//fallback to US English
	l = "en_US"
STRIP:
	// Deal with any locale names like en_US.utf8
	if i := strings.IndexRune(l, '.'); i != -1 {
		l = l[0:i]
	}
	return
}

func init() {
	language = GetLanguage()
	Internal = make(map[string]*gotext.Po)
}

// Locale gets a valid internal translation for this library, based on Language
func Locale() *gotext.Po {
	if po := Internal[GetLanguage()]; po != nil {
		return po
	}
	return Internal["en_US"]
}

// Printf writes a formatted string in the right language
func Printf(str string, vars ...interface{}) {
	print(Locale().Get(str, vars...))
}
