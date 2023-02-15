package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

type counter struct {
	opts            *options
	modelName       string // 模型名称
	daoClassName    string // DAO类名称
	daoVariableName string // DAO变量名
	daoPkgPath      string // DAO包路径
	daoPkgName      string // DAO包名称
	daoOutputDir    string // DAO文件输出目录
	daoOutputFile   string // DAO文件输出名
	daoPrefixName   string // DAO前缀名
	collectionName  string // 数据库集合名
}

func newCounter(opts *options) *counter {
	c := &counter{}
	c.opts = opts
	c.modelName = toPascalCase(opts.counterName)
	c.daoClassName = toPascalCase(c.modelName)
	c.daoVariableName = toCamelCase(c.modelName)
	c.daoOutputFile = fmt.Sprintf("%s.go", toFileName(c.modelName, c.opts.fileNameStyle))
	c.collectionName = toUnderScoreCase(c.modelName)

	dir := strings.TrimSuffix(opts.daoDir, "/")

	if opts.subpkgEnable {
		c.daoOutputDir = dir + "/" + toPackagePath(c.modelName, c.opts.subpkgStyle)
	} else {
		c.daoOutputDir = dir
		c.daoPrefixName = toPascalCase(c.modelName)
	}

	return c
}

func (c *counter) setDaoPkgPath(path string) {
	if c.opts.subpkgEnable {
		c.daoPkgPath = path + "/" + toPackagePath(c.modelName, c.opts.subpkgStyle)
	} else {
		c.daoPkgPath = path
	}

	c.daoPkgName = toPackageName(filepath.Base(c.daoPkgPath))
}
