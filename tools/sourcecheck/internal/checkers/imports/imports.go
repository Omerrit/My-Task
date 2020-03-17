package imports

import (
	"fmt"
	"gerrit-share.lan/go/tools/sourcecheck/internal/checkers"
	"gerrit-share.lan/go/utils/sets"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"
)

type importChecker struct {
	dirsSettings     map[string]dirSettings     // map[dirName]store
	packagesSettings map[string]packageSettings // map[importPath]setting
}

func NewImportChecker() checkers.Checker {
	return &importChecker{
		dirsSettings:     make(map[string]dirSettings),
		packagesSettings: make(map[string]packageSettings),
	}
}

func (c *importChecker) Setup(dirName string, packageName string) error {
	parentDir := filepath.Dir(dirName)
	parentDirSettings := c.dirsSettings[parentDir]
	packageSettings, dirSettings, err := makeSettings(dirName, packageName, parentDirSettings)
	if err != nil {
		return fmt.Errorf("Error: Import checker: %w", err)
	}
	if dirSettings.RootInternalDir != "" {
		rootPackageName := c.dirsSettings[dirSettings.RootInternalDir].ImportPath
		packageSettings.Update(c.packagesSettings[rootPackageName])
	}
	c.packagesSettings[packageName] = packageSettings
	c.dirsSettings[dirName] = dirSettings
	return nil
}

func makeSettings(dirName string, packageName string, parentDirSettings dirSettings) (packageSettings, dirSettings, error) {
	var (
		packageSettings packageSettings
		dirSettings     dirSettings
		err             error
	)
	config, err := newConfig(dirName, packageName, parentDirSettings)
	if err != nil {
		return packageSettings, dirSettings, err
	}
	packageSettings, err = config.MakePackageSettings()
	if err != nil {
		return packageSettings, dirSettings, err
	}
	dirSettings = config.MakeDirSettings()
	return packageSettings, dirSettings, nil
}

func (c *importChecker) Check(fileName string) checkers.Messages {
	if strings.HasSuffix(fileName, "_test.go") {
		return nil
	}
	dirName := filepath.Dir(fileName)
	dirSettings := c.dirsSettings[dirName]
	packageSettings := c.packagesSettings[dirSettings.ImportPath]
	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, fileName, nil, parser.AllErrors)
	var output checkers.Messages
	if err != nil {
		return output.Append(fileName, fmt.Sprintf("Error: Import checker: %s", err),
			checkers.NoLine, checkers.NoColumn)
	}
	for _, importSpec := range file.Imports {
		position := fileSet.Position(importSpec.Pos())
		importPath := trimQuotes(importSpec.Path.Value)
		if c.isInternalImport(importPath, dirName) {
			continue
		}
		if c.isSameRootInternalDir(importPath, dirSettings.RootInternalDir) {
			continue
		}
		ok, err := packageSettings.Deny.Match(importPath)
		if !ok {
			continue
		}
		if err != nil {
			output.Append(position.Filename, fmt.Sprintf("Error: Import checker: %s", err),
				checkers.NoLine, checkers.NoColumn)
		}
		if c.isAllowed(importPath, dirSettings.ImportPath) {
			continue
		}
		output.Append(position.Filename, fmt.Sprintf("bad import (%s)", importSpec.Path.Value),
			position.Line, position.Column)
	}
	return output
}

func (c *importChecker) getRootInternalDir(importPath string) string {
	return c.dirsSettings[c.packagesSettings[importPath].DirName].RootInternalDir
}

func (c *importChecker) isInternalImport(importPath string, currentDir string) bool {
	rootInternalDir := c.getRootInternalDir(importPath)
	return rootInternalDir != "" && rootInternalDir == currentDir
}

func (c *importChecker) isSameRootInternalDir(importPath string, rootInternalDir string) bool {
	packageRootInternalDir := c.getRootInternalDir(importPath)
	return packageRootInternalDir != "" && packageRootInternalDir == rootInternalDir
}

func trimQuotes(str string) string {
	if len(str) >= 2 {
		if str[0] == '"' && str[len(str)-1] == '"' || str[0] == '`' && str[len(str)-1] == '`' {
			return str[1 : len(str)-1]
		}
	}
	return str
}

func (c *importChecker) isAllowed(importPath string, startPackage string) bool {
	var (
		packagesForVisiting []string
		visitedPackages     sets.String
	)
	packagesForVisiting = append(packagesForVisiting, startPackage)
	return c.deepImportCheck(importPath, packagesForVisiting, visitedPackages)
}

func (c *importChecker) deepImportCheck(importPath string, packagesForVisiting []string, visitedPackages sets.String) bool {
	if len(packagesForVisiting) == 0 {
		return false
	}
	var newPackages []string
	for _, packageName := range packagesForVisiting {
		visitedPackages.Add(packageName)
		settings, exist := c.packagesSettings[packageName]
		if !exist {
			continue
		}
		if settings.AllowedImports.Contains(importPath) {
			return true
		}
		for allowedImport, _ := range settings.AllowedImports {
			if !visitedPackages.Contains(allowedImport) {
				newPackages = append(newPackages, allowedImport)
			}
		}
	}
	return c.deepImportCheck(importPath, newPackages, visitedPackages)
}
