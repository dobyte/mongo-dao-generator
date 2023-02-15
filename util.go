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
	kebabCase      style = "kebab"       // 短横线命名
	underScoreCase style = "under-score" // 下划线命名
	camelCase      style = "camel"       // 小驼峰命名
	pascalCase     style = "pascal"      // 大驼峰命名
	lowerCase      style = "lower"       // 小写命名
)

// 下划线命名
func toUnderScoreCase(s string) string {
	return toLowerCase(s, 95)
}

// 短横线命名
func toKebabCase(s string) string {
	return toLowerCase(s, 45)
}

// 小驼峰命名
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

// 大驼峰命名
func toPascalCase(s string) string {
	s = toCamelCase(s)
	return strings.ToUpper(string(s[0])) + s[1:]
}

// 转小写
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

// 转包名
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

// 转包路径
func toPackagePath(s string, style style) string {
	switch style {
	case kebabCase:
		return toKebabCase(s)
	case underScoreCase:
		return toUnderScoreCase(s)
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

// 转文件名
func toFileName(s string, style style) string {
	switch style {
	case kebabCase:
		return toKebabCase(s)
	case underScoreCase:
		return toUnderScoreCase(s)
	case camelCase:
		return toCamelCase(s)
	case pascalCase:
		return toPascalCase(s)
	case lowerCase:
		return toPackageName(s)
	default:
        return toUnderScoreCase(s)
	}
}

// 写文件
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
