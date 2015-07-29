package main

import (
	"go/ast"
	"go/build"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path"
	"strings"
)

func generate(dir, pkg string) {
	p, err := build.Default.Import(templateDir, dir, build.ImportMode(0))
	if err != nil {
		log.Fatalf("Import %s failed: %s", templateDir, err)
	}

	if len(p.GoFiles) != 1 {
		log.Fatalf("Expecting only 1 go file in dir %s", templateDir)
	}

	templateFilePath := path.Join(p.Dir, p.GoFiles[0])
	fset, f := parseFile(templateFilePath)

	// Change the package to the local package name
	f.Name.Name = pkg

	newDecls := []ast.Decl{}
	for _, Decl := range f.Decls {
		keep := true
		switch d := Decl.(type) {
		case *ast.GenDecl:
			// A general definition
			switch d.Tok {
			case token.TYPE:
				for _, spec := range d.Specs {
					typeSpec := spec.(*ast.TypeSpec)
					keep = !strings.HasPrefix(typeSpec.Name.Name, "Replace")
				}
			}
		}
		if keep {
			newDecls = append(newDecls, Decl)
		}
	}

	// remove the declarations that start with "Replace"
	f.Decls = newDecls
	replaceIdentifier(f, "ReplaceKey", keyType)
	replaceIdentifier(f, "ReplaceValue", valueType)

	// rename ACache and aWrapper with the type specified
	typeName := strings.ToUpper(valueType[:1]) + valueType[1:]
	replaceIdentifier(f, "ACache", typeName+"Cache")
	replaceIdentifier(f, "NewACache", "New"+typeName+"Cache")
	wrapperName := strings.ToLower(valueType[:1]) + valueType[1:]
	replaceIdentifier(f, "aCache", wrapperName+"Cache")
	replaceIdentifier(f, "aWrapper", wrapperName+"Wrapper")
	replaceIdentifier(f, "stopACacheCleanup", "stop"+typeName+"CacheCleanup")

	// output the new file
	outputFileName := strings.ToLower(valueType) + "_cache.go"
	outputFile(fset, f, outputFileName)

	log.Printf("Wrote %q", outputFileName)
}

// replaceIdentifier replaces the old string with the new string
// in the given go file
func replaceIdentifier(f *ast.File, old, new string) {
	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.Ident:
			if x.Name == old {
				x.Name = new
			}
		}
		return true
	})
}

// outputFile creates, writes and formats a go file
func outputFile(fset *token.FileSet, f *ast.File, path string) {
	fd, err := os.Create(path)
	if err != nil {
		log.Fatalf("Error opening %q: %s", path, err)
	}
	if err := format.Node(fd, fset, f); err != nil {
		log.Fatalf("Error formatting %q: %s", path, err)
	}
	if err := fd.Close(); err != nil {
		log.Fatalf("Error closing %q: %s", path, err)
	}
}

// parseFile gets a Fileset and a ast.File for a particular go file
func parseFile(path string) (*token.FileSet, *ast.File) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		log.Fatalf("Error parsing file: %s", err)
	}
	return fset, f
}
