package style

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"example.com/m/config"
)

const (
	Reset     = "\033[0m"
	Black     = "\033[30m"
	Red       = "\033[31m"
	Green     = "\033[32m"
	Yellow    = "\033[33m"
	Blue      = "\033[34m"
	Magenta   = "\033[35m"
	Cyan      = "\033[36m"
	Gray      = "\033[37m"
	White     = "\033[97m"
	Bold      = "\033[1m"
	Italic    = "\033[3m"
	Underline = "\033[4m"
	Invert    = "\033[7m"
)

func CheckColourSupport() {
	_, err := fmt.Fprint(os.Stdout, "\x1b[31m")
	if err != nil {
		config.ProvideColours = false
	}

	config.ProvideColours = !strings.Contains(fmt.Sprintf("%v", os.Stdout), "stderr")
	fmt.Fprint(os.Stdout, Reset)
}

func Colour(input interface{}, colour ...string) string {
	var s string
	c := ""
	for i := range colour {
		c = c + colour[i]
	}
	switch v := input.(type) {
	case int:
		s = c + strconv.Itoa(v) + Reset
	case bool:
		s = c + strconv.FormatBool(v) + Reset
	case []string:
		s = c + strings.Join(v, ", ") + Reset
	case string:
		if config.ProvideColours {
			s = c + v + Reset
		} else {
			s = v
		}
	default:
		fmt.Printf("Unsupported type provided to Colour func - %T\n", v)
	}
	return s
}
