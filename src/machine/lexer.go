package machine

import "regexp"

type TokenLexer struct {
	TypeName string
	Regex    *regexp.Regexp
}

var lexerMap = [...]TokenLexer{
	TokenLexer{
		TypeName: "Number",
		Regex:    regexp.MustCompile("^\\d+(\\.\\d)?$"),
	},
	TokenLexer{
		TypeName: "String",
		Regex:    regexp.MustCompile("^(\"').+?\\1$"),
	},
	TokenLexer{
		TypeName: "Operator",
		Regex:    regexp.MustCompile("^[-+*\\/\\^|&]$"),
	},
	TokenLexer{
		TypeName: "Instruction",
		Regex:    regexp.MustCompile("^(int|bool|if|jump|over|print|println|read|return|call)$"),
	},
	TokenLexer{
		TypeName: "Instruction_Args",
		Regex:    regexp.MustCompile("^(dup|swap|exit|load|store)_(\\d+)$"),
	},
}

func GetTokenTypeName(token string) (name string, arg string) {
	for _, tle := range lexerMap {
		if tle.Regex.MatchString(token) {
			match := tle.Regex.FindStringSubmatch(token)
			if len(match) == 3 {
				return tle.TypeName, match[2]
			}
			return tle.TypeName, ""
		}
	}
	return "UNKNOW TOKEN", ""
}
