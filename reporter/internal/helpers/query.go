package helpers

import (
	"fmt"
	"strings"
)

type RunnerFilter struct {
	Type string
	ID   string
}

func RunnerCondition(typeColumn, idColumn string, runners []RunnerFilter) (string, []interface{}) {
	if len(runners) == 0 {
		return "1 = 0", nil
	}

	parts := make([]string, 0, len(runners))
	args := make([]interface{}, 0, len(runners)*2)
	for _, runner := range runners {
		parts = append(parts, fmt.Sprintf("(%s = ? AND %s = ?)", typeColumn, idColumn))
		args = append(args, runner.Type, runner.ID)
	}

	return strings.Join(parts, " OR "), args
}
