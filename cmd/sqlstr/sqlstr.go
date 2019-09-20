package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/constant"
	"go/parser"
	"go/printer"
	"go/token"
	"go/types"
	"log"
	"os"
	"strconv"
	"strings"

	"golang.org/x/tools/go/ast/inspector"
)

func main() {

	fset := token.NewFileSet()

	name := "stdin.go"
	var src interface{} = os.Stdin
	if len(os.Args) >= 2 {
		name = os.Args[1]
		src = nil
	}

	f, err := parser.ParseFile(fset, name, src, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	run(fset, f)
}

func run(fset *token.FileSet, f *ast.File) (rerr error) {
	inspector.New([]*ast.File{f}).WithStack(nil, func(n ast.Node, push bool, stack []ast.Node) bool {

		if !push {
			return true
		}

		expr, ok := n.(ast.Expr)
		if !ok {
			return true
		}

		tv, err := eval(expr)
		if err != nil || tv.Value == nil || tv.Value.Kind() != constant.String {
			return true
		}
		s, err := strconv.Unquote(tv.Value.ExactString())
		if err != nil {
			rerr = err
			return false
		}

		// SQL文っぽい
		us := strings.ToUpper(strings.TrimSpace(s))
		if strings.HasPrefix(us, "SELECT") ||
			strings.HasPrefix(us, "INSERT") ||
			strings.HasPrefix(us, "UPDATE") ||
			strings.HasPrefix(us, "DELETE") {

			funcDecl := findFunc(stack)
			if funcDecl != nil {
				fmt.Printf("%s in %s\n", fset.Position(expr.Pos()), funcDecl.Name.Name)
			} else {
				fmt.Println(fset.Position(expr.Pos()))
			}
			fmt.Printf("\t%s\n\n", s)
		}

		return true
	})

	return
}

func eval(expr ast.Expr) (types.TypeAndValue, error) {
	fset := token.NewFileSet()
	var buf bytes.Buffer
	if err := printer.Fprint(&buf, fset, expr); err != nil {
		return types.TypeAndValue{}, err
	}
	pkg := types.NewPackage("main", "main")
	return types.Eval(fset, pkg, token.NoPos, buf.String())
}

func findFunc(stack []ast.Node) *ast.FuncDecl {
	for i := range stack {
		if funcdecl, ok := stack[i].(*ast.FuncDecl); ok {
			return funcdecl
		}
	}
	return nil
}
