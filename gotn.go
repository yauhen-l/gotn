package main

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"runtime"
)

var FileSet = token.NewFileSet()

func main() {
	filename := "/home/yauhen/ws/golang/src/github.com/yauhenl/gotn/gotn_test.go"
	src, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("cannot read file %q: %v", filename, err)
	}
	f, err := parser.ParseFile(FileSet, filename, src, parser.Trace)
	if err != nil {
		log.Fatalf("cannot parse file %q: %v", filename, err)
	}
	o := findIdentifier(f, 65)
	log.Print(o)
}

func defaultImportPathToName(path, srcDir string) (string, error) {
	pkg, err := build.Default.Import(path, srcDir, 0)
	return pkg.Name, err
}

func findIdentifier(f *ast.File, searchpos int) ast.Node {
	ec := make(chan ast.Node)

	found := func(startPos, endPos token.Pos) bool {
		start := FileSet.Position(startPos).Offset
		end := start + int(endPos-startPos)
		return start <= searchpos && searchpos <= end
	}

	go func() {
		var curTestFunc string

		visit := func(n ast.Node) bool {
			var startPos token.Pos

			switch n := n.(type) {
			case *ast.CallExpr:
				caseName, ok := isRunTestCase(n)
				if !ok {
					return false
				}
			case *ast.FuncLit:
				if searchpos < int(n.Pos()) || int(n.End()) < searchpos {
					return false
				}

				if !isTestFunc(n.Type) {
					return false
				}

				return true
			case *ast.FuncDecl:
				if searchpos < int(n.Pos()) || int(n.End()) < searchpos {
					return false
				}

				if !isTestFunc(n.Type) {
					return false
				}
				log.Printf("func %s at [%d, %d]", n.Name.String(), n.Pos(), n.End())
				curTestFunc = n.Name.String()
				log.Printf(curTestFunc)
				return true
			default:
				return true
			}
			if found(startPos, n.End()) {
				ec <- n
				runtime.Goexit()
			}
			return true
		}
		ast.Walk(FVisitor(visit), f)
		ec <- nil
	}()
	ev := <-ec
	if ev == nil {
		log.Fatal("no identifier found")
	}
	return ev
}

func isRunTestCase(c *ast.CallExpr) (string, bool) {
	if len(c.Args) != 2 {
		return "", false
	}

	sel, ok := c.Fun.(*ast.SelectorExpr)
	if !ok {
		return "", false
	}

	if sel.Sel.Name != "Run" {
		return "", false
	}

	f, ok := c.Args[1].(*ast.FuncLit)
	if !ok {
		return "", false
	}

	if !isTestFunc(f.Type) {
		return "", false
	}

	name, ok := c.Args[0].(*ast.BasicLit)
	if !ok {
		return "", true
	}

	return name.Value, true
}

func isTestFunc(ft *ast.FuncType) bool {
	if len(ft.Params.List) != 1 {
		return false
	}
	star, ok := ft.Params.List[0].Type.(*ast.StarExpr)
	if !ok {
		return false
	}
	sel, ok := star.X.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	return fmt.Sprintf("%s.%s", sel.X, sel.Sel.Name) == "testing.T"
}

func isTestingTExpr(expr string) ast.Expr {
	n, err := parser.ParseExpr(expr)
	if err != nil {
		log.Fatalf("cannot parse expression %q: %v", expr, err)
	}
	switch n := n.(type) {
	case *ast.Ident, *ast.SelectorExpr:
		return n
	}
	log.Fatalf("no identifier found in expression %q", expr)
	return nil
}

type FVisitor func(n ast.Node) bool

func (f FVisitor) Visit(n ast.Node) ast.Visitor {
	if f(n) {
		return f
	}
	return nil
}
