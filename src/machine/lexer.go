package machine

import "regexp"

type TokenLexer struct {
	TypeName string
	Regex    *regexp.Regexp
}

var labelRegex = regexp.MustCompile("^(.+?):(.+?)$")
var labelHeadRegex = regexp.MustCompile("^(.+?):$")
var stringRegex = regexp.MustCompile("^(\"|').+?(\"|')$")
var zeroRegex = regexp.MustCompile("^0(.0+)?|.0+$")
var numRegex = regexp.MustCompile("^(-?\\d+)(\\.\\d)?$")

var lexerMap = [...]TokenLexer{
	TokenLexer{
		TypeName: "Number",
		Regex:    regexp.MustCompile("^-?\\d+(\\.\\d)?$"),
	},
	TokenLexer{
		TypeName: "String",
		Regex:    stringRegex,
	},
	TokenLexer{
		TypeName: "Label_Pointer",
		Regex:    regexp.MustCompile("^&.+$"),
	},
	TokenLexer{
		TypeName: "Operator",
		Regex:    regexp.MustCompile("^(mul|plus|minus|div|eql|equal|greater|less|or|and|xor|not)$"),
	},
	TokenLexer{
		TypeName: "Instruction",
		Regex:    regexp.MustCompile("^(null|int|float|num|bool|if|jump|over|print|println|read|return|call|dup|swap|exit|load|store)$"),
	},
	TokenLexer{
		TypeName: "Instruction_Args",
		Regex:    regexp.MustCompile("^(dup|swap|exit)_(.+)$"),
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
