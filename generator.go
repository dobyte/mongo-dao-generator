package main

import (
	"fmt"
    "github.com/dobyte/mongo-dao-generator/template"
	"go/ast"
	"go/token"
	"golang.org/x/tools/go/packages"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

const (
	symbolBacktick = "`" // 反引号
)

const (
	symbolBacktickKey = "SymbolBacktick" // 反引号
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
	varDaoPrefixNameKey        = "VarDaoPrefixName"        // dao包前缀
	varCollectionNameKey       = "VarCollectionName"       // 集合名称
	varAutofillCodeKey         = "VarAutofillCode"         // 自动填充代码
)

const defaultCounterName = "Counter"

type options struct {
	modelDir      string   // 模型所在目录。必填。
	modelPkgPath  string   // 模型的包路径。选填，默认为空；不填写则默认根据模块和模型所在目录动态计算得到。
	modelPkgAlias string   // 模型的包别名。选填，默认为空；不填写则默认使用模型的包名。
	modelNames    []string // 模型类型名称。
	daoDir        string   // DAO文件输出目录
	daoPkgPath    string   // DAO文件输出目录对应的包路径
	subpkgEnable  bool     // 是否启用子包。
	subpkgStyle   style    // 子包风格。
	counterName   string   // 计数器类型名。选填，默认为counter
	fileNameStyle style    // 文件名风格
}

type generator struct {
	opts       *options
	counter    *counter
	modelNames map[string]struct{}
}

func newGenerator(opts *options) *generator {
	modelNames := make(map[string]struct{}, len(opts.modelNames))
	for _, modelName := range opts.modelNames {
		if isExportable(modelName) {
			modelNames[modelName] = struct{}{}
		}
	}

	if len(modelNames) == 0 {
		log.Fatalf("error: %d model type names found", len(modelNames))
	}

	if opts.counterName == "" {
		opts.counterName = defaultCounterName
	}

	return &generator{
		opts:       opts,
		counter:    newCounter(opts),
		modelNames: modelNames,
	}
}

func (g *generator) makeDao() {
	models := g.parseModels()

	for _, m := range models {
		g.makeModelInternalDao(m)

		g.makeModelExternalDao(m)

		if !m.isDependCounter {
			continue
		}

		g.makeCounterInternalDao()

		g.makeCounterExternalDao()
	}
}

// 生成内部DAO
func (g *generator) makeModelInternalDao(m *model) {
	replaces := make(map[string]string)
	replaces[varModelClassNameKey] = m.modelClassName
	replaces[varModelPackageNameKey] = m.modelPkgName
	replaces[varModelPackagePathKey] = m.modelPkgPath
	replaces[varModelVariableNameKey] = m.modelVariableName
	replaces[varDaoPrefixNameKey] = m.daoPrefixName
	replaces[varDaoClassNameKey] = m.daoClassName
	replaces[varDaoVariableNameKey] = m.daoVariableName
	replaces[varCollectionNameKey] = m.collectionName
	replaces[varModelColumnsDefineKey] = m.modelColumnsDefined()
	replaces[varModelColumnsInstanceKey] = m.modelColumnsInstance()
	replaces[varAutofillCodeKey] = m.autoFillCode()
	replaces[varPackagesKey] = m.packages()

	file := m.daoOutputDir + "/internal/" + m.daoOutputFile

	err := doWrite(file, template.InternalTemplate, replaces)
	if err != nil {
		log.Fatal(err)
	}
}

// 生成外部DAO
func (g *generator) makeModelExternalDao(m *model) {
	file := m.daoOutputDir + "/" + m.daoOutputFile

	_, err := os.Stat(file)
	if err != nil {
		switch {
		case os.IsNotExist(err):
		// ignore
		case os.IsExist(err):
			return
		default:
			log.Fatal(err)
		}
	} else {
		return
	}

	replaces := make(map[string]string)
	replaces[varDaoClassNameKey] = m.daoClassName
	replaces[varDaoPrefixNameKey] = m.daoPrefixName
	replaces[varDaoPackageNameKey] = m.daoPkgName
	replaces[varDaoPackagePathKey] = m.daoPkgPath

	err = doWrite(file, template.ExternalTemplate, replaces)
	if err != nil {
		log.Fatal(err)
	}
}

// 生成计数器内部DAO
func (g *generator) makeCounterInternalDao() {
	replaces := make(map[string]string)
	replaces[varDaoClassNameKey] = g.counter.daoClassName
	replaces[varDaoPrefixNameKey] = g.counter.daoPrefixName
	replaces[varDaoVariableNameKey] = g.counter.daoVariableName
	replaces[varCollectionNameKey] = g.counter.collectionName
	replaces[symbolBacktickKey] = symbolBacktick

	file := g.counter.daoOutputDir + "/internal/" + g.counter.daoOutputFile

	err := doWrite(file, template.CounterInternalTemplate, replaces)
	if err != nil {
		log.Fatal(err)
	}
}

// 生成计数器外部DAO
func (g *generator) makeCounterExternalDao() {
	file := g.counter.daoOutputDir + "/" + g.counter.daoOutputFile

	_, err := os.Stat(file)
	if err != nil {
		switch {
		case os.IsNotExist(err):
		// ignore
		case os.IsExist(err):
			return
		default:
			log.Fatal(err)
		}
	} else {
		return
	}

	replaces := make(map[string]string)
	replaces[varDaoClassNameKey] = g.counter.daoClassName
	replaces[varDaoPackageNameKey] = g.counter.daoPkgName
	replaces[varDaoPackagePathKey] = g.counter.daoPkgPath

	err = doWrite(file, template.CounterExternalTemplate, replaces)
	if err != nil {
		log.Fatal(err)
	}
}

// 解析模型
func (g *generator) parseModels() []*model {
	var (
		pkg          = g.loadPackage()
		models       = make([]*model, 0, len(pkg.Syntax))
		daoPkgPath   = g.opts.daoPkgPath
		modelPkgPath = g.opts.modelPkgPath
		modelPkgName = g.opts.modelPkgAlias
	)

	if g.opts.daoPkgPath == "" && pkg.Module != nil {
		outPath, err := filepath.Abs(g.opts.daoDir)
		if err != nil {
			log.Fatal(err)
		}
		daoPkgPath = pkg.Module.Path + outPath[len(pkg.Module.Dir):]
	}

	g.counter.setDaoPkgPath(daoPkgPath)

	for _, file := range pkg.Syntax {
		if g.opts.modelPkgPath == "" && pkg.Module != nil && pkg.Fset != nil {
			filePath := filepath.Dir(pkg.Fset.Position(file.Package).Filename)
			modelPkgPath = pkg.Module.Path + filePath[len(pkg.Module.Dir):]
		}

		modelPkgName = file.Name.Name

		ast.Inspect(file, func(node ast.Node) bool {
			decl, ok := node.(*ast.GenDecl)
			if !ok || decl.Tok != token.TYPE {
				return true
			}

			for _, s := range decl.Specs {
				spec, ok := s.(*ast.TypeSpec)
				if !ok {
					continue
				}

				_, ok = g.modelNames[spec.Name.Name]
				if !ok {
					continue
				}

				st, ok := spec.Type.(*ast.StructType)
				if !ok {
					continue
				}

				model := newModel(g.opts)
				model.setModelName(spec.Name.Name)
				model.setModelPkg(modelPkgName, modelPkgPath)
				model.setDaoPkgPath(daoPkgPath)

				for _, item := range st.Fields.List {
					name := item.Names[0].Name

					if !isExportable(name) {
						continue
					}

					field := &field{name: name, column: name}

					if item.Tag != nil && len(item.Tag.Value) > 2 {
						runes := []rune(item.Tag.Value)
						if runes[0] != '`' || runes[len(runes)-1] != '`' {
							continue
						}

						tag := reflect.StructTag(runes[1 : len(runes)-1])

						if column := tag.Get("bson"); column != "" {
							field.column = column
						}

						val, ok := tag.Lookup("gen")
						if ok {
							parts := strings.Split(val, ";")
							for _, part := range parts {
								if part == "" {
									continue
								}

								switch eles := strings.SplitN(part, ":", 2); eles[0] {
								case "autoFill":
									expr, ok := item.Type.(*ast.SelectorExpr)
									if !ok {
										continue
									}

									switch fmt.Sprintf("%s.%s", expr.X.(*ast.Ident).Name, expr.Sel.Name) {
									case "primitive.ObjectID":
										field.autoFill = objectID
										model.addImport(pkg3)
									case "primitive.DateTime":
										field.autoFill = dateTime
										model.addImport(pkg1)
										model.addImport(pkg3)
									}
								case "autoIncr":
									if len(eles) != 2 || eles[1] == "" {
										continue
									}

									expr, ok := item.Type.(*ast.Ident)
									if !ok {
										continue
									}

									switch expr.Name {
									case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64":
										field.autoFill = autoIncr
										field.autoIncrFieldName = eles[1]

										if g.opts.subpkgEnable {
											model.addImport(g.counter.daoPkgPath)
										}

										switch expr.Name {
										case "int":
											field.autoIncrFieldKind = reflect.Int
										case "int8":
											field.autoIncrFieldKind = reflect.Int8
										case "int16":
											field.autoIncrFieldKind = reflect.Int16
										case "int32":
											field.autoIncrFieldKind = reflect.Int32
										case "int64":
											field.autoIncrFieldKind = reflect.Int64
										case "uint":
											field.autoIncrFieldKind = reflect.Uint
										case "uint8":
											field.autoIncrFieldKind = reflect.Uint8
										case "uint16":
											field.autoIncrFieldKind = reflect.Uint16
										case "uint32":
											field.autoIncrFieldKind = reflect.Uint32
										case "uint64":
											field.autoIncrFieldKind = reflect.Uint64
										}
									}
								}
							}
						}
					}

					if item.Doc != nil {
						field.documents = make([]string, 0, len(item.Doc.List))
						for _, doc := range item.Doc.List {
							field.documents = append(field.documents, doc.Text)
						}
					}

					if item.Comment != nil {
						field.comment = item.Comment.List[0].Text
					}

					model.addFields(field)
				}

				models = append(models, model)
			}

			return true
		})
	}

	return models
}

// 加载包
func (g *generator) loadPackage() *packages.Package {
	cfg := &packages.Config{
		Mode:  packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles | packages.NeedImports | packages.NeedTypes | packages.NeedTypesSizes | packages.NeedSyntax | packages.NeedTypesInfo | packages.NeedModule,
		Tests: false,
	}
	pkgs, err := packages.Load(cfg, g.opts.modelDir)
	if err != nil {
		log.Fatal(err)
	}
	if len(pkgs) != 1 {
		log.Fatalf("error: %d packages found", len(pkgs))
	}

	return pkgs[0]
}
