package machine

import "strconv"
import "regexp"

func str2num(str interface{}) int {
	if v, ok := str.(string); ok {
		if num, err := strconv.Atoi(v); err != nil {
			return 0
		} else {
			return num
		}
	}
	if v, ok := str.(int); ok {
		return v
	}
	// error
	return 0
}

func isStringType(text string) bool {
	if m, err := regexp.MatchString("^(\"|').+?\\1$", text); err == nil {
		return m
	}
	return false
}

func cutLabelInCode(code []string) ([][]interface{}, []string) {
	// TODO
	return nil, nil
}
