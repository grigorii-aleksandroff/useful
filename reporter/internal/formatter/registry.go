package formatter

import (
	"fmt"

	"git.bububla.com/kilogramix/asia/reporter.git/internal/defs"
)

func Register(formatters ...Formatter) (map[defs.FormatterCode]Formatter, map[defs.FormatterCode]bool, error) {
	result := make(map[defs.FormatterCode]Formatter, len(formatters))
	codes := make(map[defs.FormatterCode]bool, len(formatters))
	for _, current := range formatters {
		if _, ok := result[current.Code()]; ok {
			return nil, nil, fmt.Errorf("formatter %q is registered twice", current.Code())
		}
		result[current.Code()] = current
		codes[current.Code()] = true
	}
	return result, codes, nil
}
