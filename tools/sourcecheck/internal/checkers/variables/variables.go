package variables

import (
	"fmt"
	"gerrit-share.lan/go/tools/sourcecheck/internal/checkers"
	"go/ast"
	"go/parser"
	"go/token"
)

type exportedVarsChecker struct{}

func NewExportVarsChecker() checkers.Checker {
	return exportedVarsChecker{}
}

func (c exportedVarsChecker) Setup(string, string) error { return nil }

func (c exportedVarsChecker) Check(fileName string) checkers.Messages {
	var output checkers.Messages
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, fileName, nil, parser.AllErrors)
	if err != nil {
		return output.Append(fileName, fmt.Sprintf("Error: Export checker: %s", err), checkers.NoLine, checkers.NoColumn)
	}
	localPackageName := file.Name.String()
	if localPackageName == "main" {
		return nil
	}
	output = c.findExportedVars(fset, file)
	return output
}

func (c exportedVarsChecker) findExportedVars(fileSet *token.FileSet, file *ast.File) checkers.Messages {
	var output checkers.Messages
	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.VAR {
			continue
		}
		for _, spec := range genDecl.Specs {
			valueSpec, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}
			for _, name := range valueSpec.Names {
				if !name.IsExported() {
					continue
				}
				position := fileSet.Position(name.Pos())
				output.Append(position.Filename, fmt.Sprintf("non const export: %s", name),
					position.Line, position.Column)
			}
		}
	}
	return output
}
