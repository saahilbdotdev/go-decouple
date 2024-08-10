package utils

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

var TRUE_VALUES = []string{"y", "yes", "t", "true", "on", "1"}
var FALSE_VALUES = []string{"n", "no", "f", "false", "off", "0"}

func StringToBool(value any) (bool, error) {
	switch value.(type) {
	case bool:
		return value.(bool), nil
	case string:
		valueStr := strings.ToLower(value.(string))

		if slices.Contains(TRUE_VALUES, valueStr) {
			return true, nil
		}

		if slices.Contains(FALSE_VALUES, valueStr) {
			return false, nil
		}
	}

	return false, errors.New(fmt.Sprintf("Invalid truth value: %v", value))
}

func EnvToMap(envVariables []string) map[string]string {
	envMap := make(map[string]string)

	for _, env := range envVariables {
		key, value, found := strings.Cut(env, "=")

		if found {
			envMap[key] = value
		}
	}

	return envMap
}
