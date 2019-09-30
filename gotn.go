package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

func main() {
	file := flag.String("f", "", "go test file")
	pos := flag.Int("p", 0, "position in a file")

	flag.Usage = func() {
		fmt.Print(`'gotn' determines a go test case name by an offset (-p) in a test file (-f).

Usage:
`)
		flag.PrintDefaults()
		fmt.Print(`
Examples.
1. Get test case name at offsett 350 in a gotn_test.go file:
gotn -f gotn_test.go -p 350

2. Use in conjunction with 'go test':
go test -v -run ^$(gotn -f gotn_test.go -p 350)$

`)
	}
	flag.Parse()

	if len(*file) == 0 {
		log.Fatalf("-f flag is required")
	}

	if !strings.HasSuffix(*file, "_test.go") {
		log.Fatalf("not a _test.go file")
	}

	src, err := ioutil.ReadFile(*file)
	if err != nil {
		log.Fatalf("cannot read file %q: %v", *file, err)
	}

	res := getTestNameAtPos(*file, src, *pos)
	if res == "" {
		log.Fatalf("no test function found")
	}
	fmt.Print(res)
}

func getTestNameAtPos(filename string, src []byte, pos int) string {
	fileSet := token.NewFileSet()

	f, err := parser.ParseFile(fileSet, filename, src, 0)
	if err != nil {
		log.Fatalf("cannot parse file %q: %v", filename, err)
	}
	tc := findTestCase(f, pos)

	return strings.Join(tc, "/")
}

func defaultImportPathToName(path, srcDir string) (string, error) {
	pkg, err := build.Default.Import(path, srcDir, 0)
	return pkg.Name, err
}

func findTestCase(f *ast.File, searchpos int) []string {
	var res []string

	visit := func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.CallExpr:
			if searchpos < int(n.Pos()) || int(n.End()) < searchpos {
				return false
			}

			caseName, ok := isRunTestCase(n)
			if !ok {
				return false
			}

			if len(caseName) == 0 {
				//failed to determine test case name - stop going deeper
				return false
			}
			res = append(res, rewrite(caseName))
			return true
		case *ast.FuncDecl:
			if searchpos < int(n.Pos()) || int(n.End()) < searchpos {
				return false
			}

			if !isTestFunc(n.Type) {
				return false
			}
			//log.Printf("func %s at [%d, %d]", n.Name.String(), n.Pos(), n.End())
			res = append(res, n.Name.String())

			return true
		default:
			return true
		}
	}

	ast.Walk(FVisitor(visit), f)

	return res
}

type FVisitor func(n ast.Node) bool

func (f FVisitor) Visit(n ast.Node) ast.Visitor {
	if f(n) {
		return f
	}
	return nil
}

func isRunTestCase(c *ast.CallExpr) (name string, found bool) {
	if len(c.Args) != 2 {
		return
	}

	sel, ok := c.Fun.(*ast.SelectorExpr)
	if !ok {
		return
	}

	if sel.Sel.Name != "Run" {
		return
	}

	f, ok := c.Args[1].(*ast.FuncLit)
	if !ok {
		return
	}

	if !isTestFunc(f.Type) {
		return
	}

	found = true

	bl, ok := c.Args[0].(*ast.BasicLit)
	if !ok {
		return
	}

	if bl.Kind != token.STRING {
		return
	}

	name = bl.Value
	//strip quotes
	name = string(name[1 : len(name)-1])

	return
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

// rewrite rewrites a subname to having only printable characters and no white space.
func rewrite(s string) string {
	b := []byte{}
	for _, r := range s {
		switch {
		case isSpace(r):
			b = append(b, '_')
		case !strconv.IsPrint(r):
			s := strconv.QuoteRune(r)
			b = append(b, s[1:len(s)-1]...)
		default:
			b = append(b, string(r)...)
		}
	}
	return string(b)
}

func isSpace(r rune) bool {
	if r < 0x2000 {
		switch r {
		// Note: not the same as Unicode Z class.
		case '\t', '\n', '\v', '\f', '\r', ' ', 0x85, 0xA0, 0x1680:
			return true
		}
	} else {
		if r <= 0x200a {
			return true
		}
		switch r {
		case 0x2028, 0x2029, 0x202f, 0x205f, 0x3000:
			return true
		}
	}
	return false
}
