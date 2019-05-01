package machine

import (
	"strconv"
	"strings"
)

type label_body struct {
	Value interface{}
	Idx   int
}

func num2float(num interface{}) float64 {
	if v, ok := num.(float64); ok {
		return v
	}
	if v, ok := num.(int); ok {
		return float64(v)
	}
	if v, ok := num.(string); ok {
		return str2float(v)
	}
	return float64(0)
}

func str2num(str interface{}) interface{} {
	if v, ok := str.(string); ok {
		if !zeroRegex.MatchString(v) {
			if strings.Contains(v, ".") {
				return str2float(v)
			} else {
				return str2int(v)
			}
		}
	}
	return 0
}

func integral(str string) string {
	if m := numRegex.FindStringSubmatch(str); len(m) != 0 {
		return m[1]
	}
	return "0"
}

func str2int(str interface{}) int {
	if v, ok := str.(int); ok {
		return v
	}
	if v, ok := str.(string); ok {
		intv := integral(v)
		if num, err := strconv.Atoi(intv); err == nil {
			return num
		}
	} else if v, ok := str.(int); ok {
		return v
	}
	// error
	return int(0)
}

func str2float(str interface{}) float64 {
	if v, ok := str.(float64); ok {
		return v
	}
	if v, ok := str.(string); ok {
		if num, err := strconv.ParseFloat(v, 64); err == nil {
			return num
		}
	} else if v, ok := str.(float64); ok {
		return v
	}
	// error
	return float64(0)
}

func isStringType(text string) bool {
	m := stringRegex.MatchString(text)
	return m && text[0] == text[len(text)-1]
}

// Remove extra " " or "\n" symbols
// => ["1", " \n "," 2 ","add\n"]
// <= ["1","2","add"]
func washCode(code []string) []string {
	var ret []string
	for _, line := range code {
		trimedLine := strings.Trim(line, " \n")
		if len(line) != 0 {
			ret = append(ret, trimedLine)
		}
	}
	return ret
}

// merge label =>
// => ["label1:","mul"]
// <= ["label1":"mul"]
func labelMerge(code []string) []string {
	var ret []string
	var labelHead = ""
	for _, line := range code {
		if labelHeadRegex.MatchString(line) {
			labelHead = line
		} else {
			ret = append(ret, labelHead+line)
			labelHead = ""
		}
	}
	return ret
}

func cutLabelInCode(code []string) (map[string]*label_body, []string) {
	var clearCode []string
	labels := make(map[string]*label_body)
	// wash
	code = washCode(code)
	code = labelMerge(code)

	for idx, line := range code {
		if m := labelRegex.FindStringSubmatch(line); len(m) != 0 {
			labelKey := strings.ToLower(m[1])
			labels[labelKey] = &label_body{
				Value: line,
				Idx:   idx,
			}
			clearCode = append(clearCode, m[2])
		} else {
			clearCode = append(clearCode, line)
		}
	}
	return labels, clearCode
}

func clearCode(text string) string {
	text = commentRegex.ReplaceAllString(text, "")
	text = mulitEnterRegex.ReplaceAllString(text, "\n")
	return text
}

func codeReader(code string) []string {
	code = clearCode(code)
	return exprRegex.FindAllString(code, -1)
}
