package gen

import (
	"fmt"
	"github.com/dobyte/gen-mongo-dao/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	varPackagesKey             = "VarPackages"             // 导入的包
	varModelClassNameKey       = "VarModelClassName"       // 模型类名
	varModelPackageNameKey     = "VarModelPackageName"     // 模型包名
	varModelPackagePathKey     = "VarModelPackagePath"     // 模型包路径
	varModelVariableNameKey    = "VarModelVariableName"    // 模型变量名
	varModelColumnsDefineKey   = "VarModelColumnsDefine"   // 模型列定义
	varModelColumnsInstanceKey = "VarModelColumnsInstance" // 模型列实例
	varModelClassDefineKey     = "VarModelClassDefine"     // 模型类定义
	varDaoClassNameKey         = "VarDaoClassName"         // dao类名
	varDaoVariableNameKey      = "VarDaoVariableName"      // dao变量名
	varDaoPackageNameKey       = "VarDaoPackageName"       // dao包名
	varDaoPackagePathKey       = "VarDaoPackagePath"       // dao包路径
	varDaoPrefixNameKey        = "VarDaoPrefixName"        // dao包前缀
	varCollectionNameKey       = "VarCollectionName"       // 集合名称
	varAutofillCodeKey         = "VarAutofillCode"         // 自动填充代码
)

type Generator struct {
	opts   *Options
	models []interface{}
}

type Options struct {
	OutputDir    string `json:"output_dir"`     // 输出目录
	OutputPkg    string `json:"output_pkg"`     // 输出包路径
	EnableSubPkg bool   `json:"enable_sub_pkg"` // 是否启用子包
	CounterName  string `json:"counter_name"`   // 计数器名称，默认为"Counter"
}

func NewGenerator(opts *Options) *Generator {
	if opts.CounterName == "" {
		opts.CounterName = "Counter"
	} else {
		opts.CounterName = toPascalCase(opts.CounterName)
	}

	return &Generator{
		opts:   opts,
		models: make([]interface{}, 0),
	}
}

// AddModels 添加模型
func (g *Generator) AddModels(models ...interface{}) {
	g.models = append(g.models, models...)
}

// MakeDao 批量生成DAO
func (g *Generator) MakeDao() error {
	var isMakeCounter bool

	for _, model := range g.models {
		if model == nil {
			continue
		}

		p := newParser(model, g.opts)

		if err := g.makeInternalDao(p); err != nil {
			return err
		}

		if err := g.makeExternalDao(p); err != nil {
			return err
		}

		if p.makeAutoIncrCode {
			isMakeCounter = true
		}
	}

	if isMakeCounter {
		if err := g.makeInternalCounterDao(); err != nil {
			return err
		}

		if err := g.makeExternalCounterDao(); err != nil {
			return err
		}
	}

	return nil
}

// 生成内部DAO
func (g *Generator) makeInternalDao(p *parser) error {
	dir := strings.TrimRight(g.opts.OutputDir, "/")
	if g.opts.EnableSubPkg {
		dir += fmt.Sprintf("/%s/internal/", p.modelPackageName())
	} else {
		dir += "/internal/"
	}

	replaces := make(map[string]string)
	replaces[varModelClassNameKey] = p.modelClassName()
	replaces[varModelPackageNameKey] = p.modelPackageName()
	replaces[varModelPackagePathKey] = p.modelPackagePath()
	replaces[varModelVariableNameKey] = p.modelVariableName()
	replaces[varDaoPrefixNameKey] = p.daoPrefixName()
	replaces[varDaoClassNameKey] = p.daoClassName()
	replaces[varDaoVariableNameKey] = p.daoVariableName()
	replaces[varCollectionNameKey] = p.collectionName()
	replaces[varModelColumnsDefineKey] = p.modelColumnsDefined()
	replaces[varModelColumnsInstanceKey] = p.modelColumnsInstance()
	replaces[varAutofillCodeKey] = p.autofillCode()
	replaces[varPackagesKey] = p.packages()

	return doWrite(dir+p.fileName(), template.InternalTemplate, replaces)
}

// 生成外部DAO
func (g *Generator) makeExternalDao(p *parser) error {
	var (
		dir = strings.TrimRight(g.opts.OutputDir, "/")
		pkg = strings.TrimRight(g.opts.OutputPkg, "/")
	)

	if g.opts.EnableSubPkg {
		dir += "/" + p.modelPackageName()
		pkg += "/" + p.modelPackageName()
	}

	file := dir + "/" + p.fileName()

	_, err := os.Stat(file)
	if err != nil {
		switch {
		case os.IsNotExist(err):
		// ignore
		case os.IsExist(err):
			return nil
		default:
			return err
		}
	} else {
		return nil
	}

	replaces := make(map[string]string)
	replaces[varDaoClassNameKey] = p.daoClassName()
	replaces[varDaoPackageNameKey] = toPackageName(filepath.Base(pkg))
	replaces[varDaoPackagePathKey] = pkg

	return doWrite(file, template.ExternalTemplate, replaces)
}

// 生成计数器内部DAO
func (g *Generator) makeInternalCounterDao() error {
	dir := strings.TrimRight(g.opts.OutputDir, "/")
	if g.opts.EnableSubPkg {
		dir += fmt.Sprintf("/%s/internal/", toKebabCase(g.opts.CounterName))
	} else {
		dir += "/internal/"
	}

	file := dir + "/" + toUnderScoreCase(g.opts.CounterName) + ".go"

	var modelClassDefine string
	if g.opts.EnableSubPkg {
		modelClassDefine += "type Model struct {\n"
	} else {
		modelClassDefine += fmt.Sprintf("type %sModel struct {\n", toPascalCase(g.opts.CounterName))
	}
	modelClassDefine += "\tID    string `bson:\"_id\"`\n"
	modelClassDefine += "\tValue int64  `bson:\"value\"`\n"
	modelClassDefine += "}"

	replaces := make(map[string]string)
	replaces[varDaoClassNameKey] = toPascalCase(g.opts.CounterName)
	replaces[varDaoVariableNameKey] = toCamelCase(g.opts.CounterName)
	replaces[varCollectionNameKey] = toUnderScoreCase(g.opts.CounterName)
	replaces[varModelClassDefineKey] = modelClassDefine

	if !g.opts.EnableSubPkg {
		replaces[varDaoPrefixNameKey] = toPascalCase(g.opts.CounterName)
	}

	return doWrite(file, template.CounterInternalTemplate, replaces)
}

// 生成计数器外部DAO
func (g *Generator) makeExternalCounterDao() error {
	var (
		dir = strings.TrimRight(g.opts.OutputDir, "/")
		pkg = strings.TrimRight(g.opts.OutputPkg, "/")
	)

	if g.opts.EnableSubPkg {
		packageName := toKebabCase(g.opts.CounterName)
		dir += "/" + packageName
		pkg += "/" + packageName
	}

	file := dir + "/" + toUnderScoreCase(g.opts.CounterName) + ".go"

	_, err := os.Stat(file)
	if err != nil {
		switch {
		case os.IsNotExist(err):
		// ignore
		case os.IsExist(err):
			return nil
		default:
			return err
		}
	} else {
		return nil
	}

	replaces := make(map[string]string)
	replaces[varDaoClassNameKey] = toPascalCase(g.opts.CounterName)
	replaces[varDaoPackageNameKey] = toPackageName(filepath.Base(pkg))
	replaces[varDaoPackagePathKey] = pkg

	return doWrite(file, template.CounterExternalTemplate, replaces)
}

// 写文件
func doWrite(file string, tpl string, replaces map[string]string) error {
	s := os.Expand(tpl, func(s string) string {
		return replaces[s]
	})

	if err := os.MkdirAll(filepath.Dir(file), os.ModePerm); err != nil {
		return err
	}

	return ioutil.WriteFile(file, []byte(strings.TrimLeft(s, "\n")), os.ModePerm)
}
