package main

import (
	"os"
	"path/filepath"
	"strings"
	"unicode"
	"unicode/utf8"
)

type style string

const (
	kebabCase      style = "kebab"
	underscoreCase style = "underscore"
	camelCase      style = "camel"
	pascalCase     style = "pascal"
	lowerCase      style = "lower"
)

// convert to underscore style, example: UserProfile > user_profile
func toUnderscoreCase(s string) string {
	return toLowerCase(s, 95)
}

// convert to kebab style, example: UserProfile > user-profile
func toKebabCase(s string) string {
	return toLowerCase(s, 45)
}

// convert to camel style, example: user-profile > userProfile
func toCamelCase(s string) string {
	chars := make([]rune, 0, len(s))
	upper := false
	first := true

	for i := 0; i < len(s); i++ {
		switch {
		case s[i] >= 65 && s[i] <= 90:
			if first {
				chars = append(chars, rune(s[i]+32))
			} else {
				chars = append(chars, rune(s[i]))
			}
			first = false
			upper = false
		case s[i] >= 97 && s[i] <= 122:
			if upper && !first {
				chars = append(chars, rune(s[i]-32))
			} else {
				chars = append(chars, rune(s[i]))
			}
			first = false
			upper = false
		case s[i] == 45:
			upper = true
		case s[i] == 95:
			upper = true
		}
	}

	return string(chars)
}

// convert to pascal style, example: user-profile > UserProfile
func toPascalCase(s string) string {
	s = toCamelCase(s)
	return strings.ToUpper(string(s[0])) + s[1:]
}

func toLowerCase(s string, c rune) string {
	chars := make([]rune, 0)

	for i := 0; i < len(s); i++ {
		if s[i] >= 65 && s[i] <= 90 {
			if i == 0 {
				chars = append(chars, rune(s[i]+32))
			} else {
				chars = append(chars, c, rune(s[i]+32))
			}
		} else {
			chars = append(chars, rune(s[i]))
		}
	}

	return string(chars)
}

func toPackageName(s string) string {
	chars := make([]rune, 0, len(s))
	for i := 0; i < len(s); i++ {
		switch {
		case s[i] >= 65 && s[i] <= 90:
			chars = append(chars, rune(s[i]+32))
		case s[i] >= 97 && s[i] <= 122:
			chars = append(chars, rune(s[i]))
		}
	}

	return string(chars)
}

func toPackagePath(s string, style style) string {
	switch style {
	case kebabCase:
		return toKebabCase(s)
	case underscoreCase:
		return toUnderscoreCase(s)
	case camelCase:
		return toCamelCase(s)
	case pascalCase:
		return toPascalCase(s)
	case lowerCase:
		return toPackageName(s)
	default:
		return toKebabCase(s)
	}
}

func toFileName(s string, style style) string {
	switch style {
	case kebabCase:
		return toKebabCase(s)
	case underscoreCase:
		return toUnderscoreCase(s)
	case camelCase:
		return toCamelCase(s)
	case pascalCase:
		return toPascalCase(s)
	case lowerCase:
		return toPackageName(s)
	default:
        return toUnderscoreCase(s)
	}
}

func doWrite(file string, tpl string, replaces map[string]string) error {
	s := os.Expand(tpl, func(s string) string {
		switch {
		case len(s) >= 3 && s[:3] == "Var":
			return replaces[s]
		case len(s) >= 6 && s[:6] == "Symbol":
			return replaces[s]
		default:
			return "$" + s
		}
	})

	if err := os.MkdirAll(filepath.Dir(file), os.ModePerm); err != nil {
		return err
	}

	return os.WriteFile(file, []byte(strings.TrimPrefix(s, "\n")), os.ModePerm)
}

func isExportable(s string) bool {
	r, _ := utf8.DecodeRuneInString(s)
	return unicode.IsUpper(r)
}
