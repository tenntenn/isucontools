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
)

func main() {

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "app.go", os.Stdin, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	run(f)
}

func run(f *ast.File) {
	var buf bytes.Buffer
	ast.Inspect(f, func(n ast.Node) bool {
		if expr, ok := n.(ast.Expr); ok {
			tv, err := eval(expr)
			if err == nil && tv.Value != nil && tv.Value.Kind() == constant.String {
				s, err := strconv.Unquote(tv.Value.ExactString())
				if err != nil {
					log.Fatal(nil)
				}

				// SQL文っぽい
				us := strings.ToUpper(s)
				if strings.HasPrefix(us, "SELECT") ||
					strings.HasPrefix(us, "INSERT") ||
					strings.HasPrefix(us, "UPDATE") ||
					strings.HasPrefix(us, "DELETE") {
					fmt.Fprint(&buf, s)
					return true
				}
			}
		}

		// その他のものが来た場合はそれまでの内容を出力
		if buf.Len() > 0 {
			fmt.Println(buf.String())
			buf.Reset()
		}
		return true
	})
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
