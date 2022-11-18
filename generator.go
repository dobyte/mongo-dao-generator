package gen

import (
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
	varDaoClassNameKey         = "VarDaoClassName"         // dao类名
	varDaoVariableNameKey      = "VarDaoVariableName"      // dao变量名
	varDaoPackageNameKey       = "VarDaoPackageName"       // dao包名
	varDaoPackagePathKey       = "VarDaoPackagePath"       // dao包路径
	varCollectionNameKey       = "VarCollectionName"       // 集合名称
	varAutofillCodeKey         = "VarAutofillCode"         // 自动填充代码
)

type rule struct {
	model interface{} // 模型
	out   string      // 包输出位置
	pkg   string      // 包前缀
}

type Generator struct {
	list   []rule
	prefix string
}

func NewGenerator() *Generator {
	return &Generator{
		list: make([]rule, 0),
	}
}

// AddModel 添加模型
func (g *Generator) AddModel(model interface{}, out string, pkg string) {
	g.list = append(g.list, rule{
		model: model,
		out:   out,
		pkg:   pkg,
	})
}

// MakeDao 批量生成DAO
func (g *Generator) MakeDao() error {
	for _, item := range g.list {
		if err := g.makeDao(item); err != nil {
			return err
		}
	}

	return nil
}

// 生成单个DAO
func (g *Generator) makeDao(rule rule) error {
	if rule.model == nil {
		return nil
	}

	p := newParser(rule.model)

	err := g.makeInternalDao(p, rule.out)
	if err != nil {
		return err
	}

	return g.makeExternalDao(p, rule.out, rule.pkg)
}

// 生成内部DAO
func (g *Generator) makeInternalDao(p *parser, out string) error {
	var (
		dir      = strings.TrimRight(out, "/") + "/internal/"
		replaces = make(map[string]string)
	)

	replaces[varModelClassNameKey] = p.modelClassName()
	replaces[varModelPackageNameKey] = p.modelPackageName()
	replaces[varModelPackagePathKey] = p.modelPackagePath()
	replaces[varModelVariableNameKey] = p.modelVariableName()
	replaces[varDaoClassNameKey] = p.daoClassName()
	replaces[varDaoVariableNameKey] = p.daoVariableName()
	replaces[varCollectionNameKey] = p.collectionName()
	replaces[varModelColumnsDefineKey] = p.modelColumnsDefined()
	replaces[varModelColumnsInstanceKey] = p.modelColumnsInstance()
	replaces[varAutofillCodeKey] = p.autofillCode()
	replaces[varPackagesKey] = p.packages()

	s := os.Expand(InternalTemplate, func(s string) string {
		return replaces[s]
	})

	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(dir+p.fileName(), []byte(strings.TrimLeft(s, "\n")), os.ModePerm)
}

// 生成外部DAO
func (g *Generator) makeExternalDao(p *parser, out string, pkg string) error {
	var (
		dir  = strings.TrimRight(out, "/") + "/"
		file = dir + p.fileName()
	)

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
	replaces[varDaoPackageNameKey] = strings.ReplaceAll(filepath.Base(pkg), "-", "")
	replaces[varDaoPackagePathKey] = pkg

	s := os.Expand(ExternalTemplate, func(s string) string {
		return replaces[s]
	})

	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(file, []byte(strings.TrimLeft(s, "\n")), os.ModePerm)
}
