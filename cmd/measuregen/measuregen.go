package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"os"
	"strconv"

	"golang.org/x/tools/go/ast/inspector"
)

const (
	importPATH = "github.com/najeira/measure"
)

func main() {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "main.go", os.Stdin, parser.ParseComments)
	if err != nil {
		log.Fatal("Error:", err)
	}

	closures := map[*ast.FuncDecl]int{}
	filter := []ast.Node{
		new(ast.FuncDecl),
		new(ast.FuncLit),
	}
	inspector.New([]*ast.File{f}).WithStack(filter, func(n ast.Node, push bool, stack []ast.Node) bool {
		if !push {
			return true
		}

		switch n := n.(type) {
		case *ast.FuncDecl:
			expr, err := parser.ParseExpr(fmt.Sprintf(`measure.Start("%s").Stop()`, n.Name.Name))
			if err != nil {
				log.Fatal("Error:", err)
			}
			deferStmt := &ast.DeferStmt{Call: expr.(*ast.CallExpr)}

			if hasMeasure(n.Body.List, deferStmt) {
				return true
			}

			n.Body.List = append([]ast.Stmt{deferStmt}, n.Body.List...)
		case *ast.FuncLit:
			name := "NONAME"
			if parent := findParent(stack); parent != nil {
				closures[parent]++
				name = fmt.Sprintf("%s-%d", parent.Name.Name, closures[parent])
			}
			expr, err := parser.ParseExpr(fmt.Sprintf(`measure.Start("%s").Stop()`, name))
			if err != nil {
				log.Fatal("Error:", err)
			}

			deferStmt := &ast.DeferStmt{Call: expr.(*ast.CallExpr)}

			if hasMeasure(n.Body.List, deferStmt) {
				return true
			}

			n.Body.List = append([]ast.Stmt{deferStmt}, n.Body.List...)
		}

		return true
	})

	addImport(f)

	format.Node(os.Stdout, fset, f)
}

func findParent(stack []ast.Node) *ast.FuncDecl {
	for i := range stack {
		if funcdecl, ok := stack[i].(*ast.FuncDecl); ok {
			return funcdecl
		}
	}
	return nil
}

func addImport(f *ast.File) {

	for _, im := range f.Imports {
		path, err := strconv.Unquote(im.Path.Value)
		if err != nil {
			continue
		}

		// already imported
		if path == importPATH {
			return
		}
	}

	importSpec := &ast.ImportSpec{
		Name: ast.NewIdent("measure"),
		Path: &ast.BasicLit{
			Kind:  token.STRING,
			Value: strconv.Quote(importPATH),
		},
	}
	f.Imports = append(f.Imports, importSpec)

	// Find last imnport group
	var lastGenDecl *ast.GenDecl
	for _, decl := range f.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.IMPORT {
			continue
		}
		lastGenDecl = genDecl
	}

	// The file already have import statements
	if lastGenDecl != nil {
		lastGenDecl.Specs = append(lastGenDecl.Specs, importSpec)
		return
	}

	f.Decls = append([]ast.Decl{&ast.GenDecl{
		Tok:   token.IMPORT,
		Specs: []ast.Spec{importSpec},
	}}, f.Decls...)
}

func nodeStr(node ast.Node) (string, error) {
	fset := token.NewFileSet()
	var buf bytes.Buffer
	if err := printer.Fprint(&buf, fset, node); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func hasMeasure(stmts []ast.Stmt, deferStmt *ast.DeferStmt) bool {
	if len(stmts) == 0 {
		return false
	}

	firstStmtStr, err := nodeStr(stmts[0])
	if err != nil {
		return false
	}

	deferStmtStr, err := nodeStr(deferStmt)
	if err != nil {
		return false
	}

	return firstStmtStr == deferStmtStr
}
