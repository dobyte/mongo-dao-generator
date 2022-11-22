package gen

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
	"sort"
	"strings"
)

// 自动填充类型
type autofill int

const (
	objectID autofill = iota + 1 // primitive.NewObjectID()
	dateTime                     // primitive.NewDateTimeFromTime(time.Now())
	autoIncr                     // auto-increment
)

const (
	pkg1 = "time"
	pkg2 = "context"
	pkg3 = "go.mongodb.org/mongo-driver/bson/primitive"
	pkg4 = "go.mongodb.org/mongo-driver/mongo"
	pkg5 = "go.mongodb.org/mongo-driver/mongo/options"
	pkg6 = "errors"
)

type field struct {
	name              string       // 字段名
	column            string       // 列名
	colName           string       // 列名称
	colType           string       // 列类型
	autoFill          autofill     // 列自动填充类型
	autoIncrFieldName string       // 自增字段名
	autoIncrFieldKind reflect.Kind // 自增字段类型
}

type parser struct {
	rv               reflect.Value
	rt               reflect.Type
	rk               reflect.Kind
	fields           []field
	fieldNameMaxLen  int
	imports          map[string]struct{}
	makeAutoIncrCode bool
	opts             *Options
}

func newParser(model interface{}, opts *Options) *parser {
	p := &parser{
		fields:  make([]field, 0),
		imports: make(map[string]struct{}),
		opts:    opts,
	}
	p.parse(model)
	return p
}

func (p *parser) parse(model interface{}) {
	p.rv = reflect.ValueOf(model)
	p.rt = p.rv.Type()
	p.rk = p.rv.Kind()

	for p.rk == reflect.Ptr {
		p.rv = p.rv.Elem()
		p.rt = p.rv.Type()
		p.rk = p.rv.Kind()
	}

	for i := 0; i < p.rv.NumField(); i++ {
		fv := p.rv.Field(i)
		fk := fv.Kind()

		for fk == reflect.Ptr {
			fv = fv.Elem()
			fk = fv.Kind()
		}

		f := field{
			name:    p.rt.Field(i).Name,
			colName: p.rt.Field(i).Tag.Get("bson"),
			colType: "string",
		}

		val, ok := p.rt.Field(i).Tag.Lookup("gen")
		if ok {
			parts := strings.Split(val, ";")
			for _, part := range parts {
				if part == "" {
					continue
				}

				elements := strings.SplitN(part, ":", 2)
				switch elements[0] {
				case "autoFill":
					switch fv.Interface().(type) {
					case primitive.ObjectID:
						f.autoFill = objectID
						p.imports[pkg3] = struct{}{}
					case primitive.DateTime:
						f.autoFill = dateTime
						p.imports[pkg1] = struct{}{}
						p.imports[pkg3] = struct{}{}
					}
				case "autoIncr":
					if len(elements) != 2 || elements[1] == "" {
						continue
					}

					switch fk {
					case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
						f.autoFill = autoIncr
						f.autoIncrFieldName = elements[1]
						f.autoIncrFieldKind = fk
						p.makeAutoIncrCode = true
					}
				}
			}
		}

		if p.makeAutoIncrCode && p.opts.EnableSubPkg {
			pkg := fmt.Sprintf("%s/%s", strings.TrimRight(p.opts.OutputPkg, "/"), toKebabCase(p.opts.CounterName))
			p.imports[pkg] = struct{}{}
		}

		if l := len(f.name); l > p.fieldNameMaxLen {
			p.fieldNameMaxLen = l
		}

		p.fields = append(p.fields, f)
	}
}

// 获取模型类名
func (p *parser) modelClassName() string {
	return toPascalCase(p.rt.Name())
}

// 获取模型变量名
func (p *parser) modelVariableName() string {
	return toCamelCase(p.rt.Name())
}

// 获取模型包名称
func (p *parser) modelPackageName() string {
	parts := strings.Split(p.modelPackagePath(), "/")
	return parts[len(parts)-1]
}

// 获取模型包路径
func (p *parser) modelPackagePath() string {
	return p.rt.PkgPath()
}

// 获取DAO前缀名
func (p *parser) daoPrefixName() string {
	if p.opts.EnableSubPkg {
		return ""
	} else {
		return toPascalCase(p.rt.Name())
	}
}

// 获取DAO类名称
func (p *parser) daoClassName() string {
	return toPascalCase(p.rt.Name())
}

// 获取DAO的变量名
func (p *parser) daoVariableName() string {
	return toCamelCase(p.rt.Name())
}

// 获取集合名称
func (p *parser) collectionName() string {
	return toUnderScoreCase(p.rt.Name())
}

// 获取文件名称
func (p *parser) fileName() string {
	return toUnderScoreCase(p.rt.Name()) + ".go"
}

// 获取要导入的包
func (p *parser) packages() (str string) {
	p.imports[pkg2] = struct{}{}
	p.imports[pkg4] = struct{}{}
	p.imports[pkg5] = struct{}{}
	p.imports[pkg6] = struct{}{}
	p.imports[p.modelPackagePath()] = struct{}{}

	packages := make([]string, 0, len(p.imports))
	for pkg := range p.imports {
		packages = append(packages, pkg)
	}

	sort.Slice(packages, func(i, j int) bool {
		return packages[i] < packages[j]
	})

	for _, pkg := range packages {
		str += fmt.Sprintf("\t\"%s\"\n", pkg)
	}

	str = strings.TrimLeft(str, "\t")
	str = strings.TrimRight(str, "\n")
	return
}

// 获取模型列定义
func (p *parser) modelColumnsDefined() (str string) {
	for i, f := range p.fields {
		str += fmt.Sprintf("\t%s%s%s", f.name, strings.Repeat(" ", p.fieldNameMaxLen-len(f.name)+1), f.colType)
		if i != len(p.fields)-1 {
			str += "\n"
		}
	}

	str = strings.TrimLeft(str, "\t")
	return
}

// 获取模型列实例
func (p *parser) modelColumnsInstance() (str string) {
	for i, f := range p.fields {
		str += fmt.Sprintf("\t%s:%s\"%s\",", f.name, strings.Repeat(" ", p.fieldNameMaxLen-len(f.name)+1), f.colName)
		if i != len(p.fields)-1 {
			str += "\n"
		}
	}

	str = strings.TrimLeft(str, "\t")
	return
}

// 自动填充代码
func (p *parser) autofillCode() (str string) {
	var counterPkgPrefix string
	if p.opts.EnableSubPkg {
		counterPkgPrefix = fmt.Sprintf("%s.", toPackageName(p.opts.CounterName))
	}

	for _, f := range p.fields {
		if f.autoFill == 0 {
			continue
		}

		if str != "" {
			str += "\n\n"
		}

		switch f.autoFill {
		case objectID:
			str += fmt.Sprintf("\tif model.%s.IsZero() {\n", f.name)
			str += fmt.Sprintf("\t\tmodel.%s = primitive.NewObjectID()\n", f.name)
			str += "\t}"
		case dateTime:
			str += fmt.Sprintf("\tif model.%s == 0 {\n", f.name)
			str += fmt.Sprintf("\t\tmodel.%s = primitive.NewDateTimeFromTime(time.Now())\n", f.name)
			str += "\t}"
		case autoIncr:
			str += fmt.Sprintf("\tif model.%s == 0 {\n", f.name)
			str += fmt.Sprintf("\t\tif id, err := %sNew%s(dao.Database).Incr(ctx, \"%s\"); err != nil {\n", counterPkgPrefix, p.opts.CounterName, f.autoIncrFieldName)
			str += "\t\t\treturn err\n"
			str += "\t\t} else {\n"

			switch f.autoIncrFieldKind {
			case reflect.Int:
				str += fmt.Sprintf("\t\t\tmodel.%s = int(id)\n", f.name)
			case reflect.Int8:
				str += fmt.Sprintf("\t\t\tmodel.%s = int8(id)\n", f.name)
			case reflect.Int16:
				str += fmt.Sprintf("\t\t\tmodel.%s = int16(id)\n", f.name)
			case reflect.Int32:
				str += fmt.Sprintf("\t\t\tmodel.%s = int32(id)\n", f.name)
			case reflect.Int64:
				str += fmt.Sprintf("\t\t\tmodel.%s = id\n", f.name)
			case reflect.Uint:
				str += fmt.Sprintf("\t\t\tmodel.%s = uint(id)\n", f.name)
			case reflect.Uint8:
				str += fmt.Sprintf("\t\t\tmodel.%s = uint8(id)\n", f.name)
			case reflect.Uint16:
				str += fmt.Sprintf("\t\t\tmodel.%s = uint16(id)\n", f.name)
			case reflect.Uint32:
				str += fmt.Sprintf("\t\t\tmodel.%s = uint32(id)\n", f.name)
			case reflect.Uint64:
				str += fmt.Sprintf("\t\t\tmodel.%s = uint64(id)\n", f.name)
			}

			str += "\t\t}\n"
			str += "\t}"
		}
	}

	if str != "" {
		str += "\n\n"
	}

	str += "\treturn nil"
	str = strings.TrimLeft(str, "\t")

	return
}

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
