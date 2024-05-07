package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

type counter struct {
	opts            *options
	modelName       string
	daoClassName    string
	daoVariableName string
	daoPkgPath      string
	daoPkgName      string
	daoOutputDir    string
	daoOutputFile   string
	daoPrefixName   string
	collectionName  string
}

func newCounter(opts *options) *counter {
	c := &counter{}
	c.opts = opts
	c.modelName = toPascalCase(opts.counterName)
	c.daoClassName = toPascalCase(c.modelName)
	c.daoVariableName = toCamelCase(c.modelName)
	c.daoOutputFile = fmt.Sprintf("%s.go", toFileName(c.modelName, c.opts.fileNameStyle))
	c.collectionName = toUnderscoreCase(c.modelName)

	dir := strings.TrimSuffix(opts.daoDir, "/")

	if opts.subPkgEnable {
		c.daoOutputDir = dir + "/" + toPackagePath(c.modelName, c.opts.subPkgStyle)
	} else {
		c.daoOutputDir = dir
		c.daoPrefixName = toPascalCase(c.modelName)
	}

	return c
}

func (c *counter) setDaoPkgPath(path string) {
	if c.opts.subPkgEnable {
		c.daoPkgPath = path + "/" + toPackagePath(c.modelName, c.opts.subPkgStyle)
	} else {
		c.daoPkgPath = path
	}

	c.daoPkgName = toPackageName(filepath.Base(c.daoPkgPath))
}
