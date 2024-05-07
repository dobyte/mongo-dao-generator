package main

import (
	"fmt"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
)

const (
	defaultModelPkgAlias     = "modelpkg"
	defaultModelVariableName = "model"
)

type autoFill int

const (
	objectID autoFill = iota + 1 // primitive.NewObjectID()
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
	pkg7 = "go.mongodb.org/mongo-driver/bson"
)

type field struct {
	name              string
	column            string
	comment           string
	documents         []string
	autoFill          autoFill
	autoIncrFieldName string
	autoIncrFieldKind reflect.Kind
}

type model struct {
	opts               *options
	fields             []*field
	imports            map[string]string
	modelName          string
	modelClassName     string
	modelVariableName  string
	modelPkgPath       string
	modelPkgName       string
	daoClassName       string
	daoVariableName    string
	daoPkgPath         string
	daoPkgName         string
	daoOutputDir       string
	daoOutputFile      string
	daoPrefixName      string
	collectionName     string
	fieldNameMaxLen    int
	fieldComplexMaxLen int
	isDependCounter    bool
}

func newModel(opts *options) *model {
	m := &model{
		opts:    opts,
		fields:  make([]*field, 0),
		imports: make(map[string]string, 8),
	}

	m.addImport(pkg2)
	m.addImport(pkg4)
	m.addImport(pkg5)
	m.addImport(pkg6)
	m.addImport(pkg7)

	return m
}

func (m *model) setModelName(name string) {
	m.modelName = name
	m.modelClassName = toPascalCase(m.modelName)
	m.modelVariableName = toCamelCase(m.modelName)
	m.daoClassName = toPascalCase(m.modelName)
	m.daoVariableName = toCamelCase(m.modelName)
	m.daoOutputFile = fmt.Sprintf("%s.go", toFileName(m.modelName, m.opts.fileNameStyle))
	m.collectionName = toUnderscoreCase(m.modelName)

	dir := strings.TrimSuffix(m.opts.daoDir, "/")

	if m.opts.subPkgEnable {
		m.daoOutputDir = dir + "/" + toPackagePath(m.modelName, m.opts.subPkgStyle)
	} else {
		m.daoOutputDir = dir
		m.daoPrefixName = toPascalCase(m.modelName)
	}
}

func (m *model) setModelPkg(name, path string) {
	m.modelPkgPath = path

	if m.opts.modelPkgAlias != "" {
		m.modelPkgName = m.opts.modelPkgAlias
		m.addImport(m.modelPkgPath, m.modelPkgName)
	} else {
		m.modelPkgName = name
		m.addImport(m.modelPkgPath)
	}

	if m.modelPkgName == defaultModelVariableName {
		m.modelPkgName = defaultModelPkgAlias
		m.addImport(m.modelPkgPath, m.modelPkgName)
	}
}

func (m *model) setDaoPkgPath(path string) {
	if m.opts.subPkgEnable {
		m.daoPkgPath = path + "/" + toPackagePath(m.modelName, m.opts.subPkgStyle)
	} else {
		m.daoPkgPath = path
	}

	m.daoPkgName = toPackageName(filepath.Base(m.daoPkgPath))
}

func (m *model) addImport(pkg string, alias ...string) {
	if len(alias) > 0 {
		m.imports[pkg] = alias[0]
	} else {
		m.imports[pkg] = ""
	}
}

func (m *model) addFields(fields ...*field) {
	for _, f := range fields {
		if l := len(f.name); l > m.fieldNameMaxLen {
			m.fieldNameMaxLen = l
		}

		if l := len(f.name) + len(f.column) + 5; l > m.fieldComplexMaxLen {
			m.fieldComplexMaxLen = l
		}

		if f.autoFill == autoIncr {
			m.isDependCounter = true
		}
	}

	m.fields = append(m.fields, fields...)
}

func (m *model) modelColumnsDefined() (str string) {
	for i, f := range m.fields {
		str += fmt.Sprintf("\t%s%s%s %s", f.name, strings.Repeat(" ", m.fieldNameMaxLen-len(f.name)+1), "string", f.comment)
		if i != len(m.fields)-1 {
			str += "\n"
		}
	}

	str = strings.TrimPrefix(str, "\t")
	return
}

func (m *model) modelColumnsInstance() (str string) {
	for i, f := range m.fields {
		s := fmt.Sprintf("%s:%s\"%s\",", f.name, strings.Repeat(" ", m.fieldNameMaxLen-len(f.name)+1), f.column)
		s += strings.Repeat(" ", m.fieldComplexMaxLen-len(s)+1) + f.comment
		str += "\t" + s
		if i != len(m.fields)-1 {
			str += "\n"
		}
	}

	str = strings.TrimLeft(str, "\t")
	return
}

func (m *model) packages() (str string) {
	packages := make([]string, 0, len(m.imports))
	for pkg := range m.imports {
		packages = append(packages, pkg)
	}

	sort.Slice(packages, func(i, j int) bool {
		return packages[i] < packages[j]
	})

	for _, pkg := range packages {
		if alias := m.imports[pkg]; alias != "" {
			str += fmt.Sprintf("\t%s \"%s\"\n", alias, pkg)
		} else {
			str += fmt.Sprintf("\t\"%s\"\n", pkg)
		}
	}

	str = strings.TrimPrefix(str, "\t")
	str = strings.TrimSuffix(str, "\n")
	return
}

func (m *model) autoFillCode() (str string) {
	var (
		counterName      = toPascalCase(m.opts.counterName)
		counterPkgPrefix string
	)

	if m.opts.subPkgEnable {
		counterPkgPrefix = fmt.Sprintf("%s.", toPackageName(counterName))
	}

	for _, f := range m.fields {
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
			str += fmt.Sprintf("\t\tif id, err := %sNew%s(dao.Database).Incr(ctx, \"%s\"); err != nil {\n", counterPkgPrefix, counterName, f.autoIncrFieldName)
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
	str = strings.TrimPrefix(str, "\t")

	return
}
