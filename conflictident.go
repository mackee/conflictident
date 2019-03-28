package conflictident

import (
	"go/ast"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer is analysis.Analyzer settings for conflictident
var Analyzer = &analysis.Analyzer{
	Name: "conflictident",
	Doc:  Doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

// Doc is description for the CLI
const Doc = "conflictident is linter that discover conflicts ident"

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	//inspect.Preorder(nil, func(n ast.Node) {
	//	fmt.Printf("%T %+v\n", n, n)
	//})

	nodeFilter := []ast.Node{
		(*ast.File)(nil),
		(*ast.TypeSpec)(nil),
		(*ast.ValueSpec)(nil),
		(*ast.FuncDecl)(nil),
		(*ast.AssignStmt)(nil),
		(*ast.Ident)(nil),
	}

	imports := map[string]string{}
	for _, imp := range pass.Pkg.Imports() {
		imports[imp.Path()] = imp.Name()
	}

	pkgScopeIdent := map[string]ast.Node{}
	scopeIdent := map[token.Pos]map[string]ast.Node{}
	inspect.WithStack(nodeFilter, func(n ast.Node, push bool, stack []ast.Node) bool {
		switch n := n.(type) {
		case *ast.File:
			for _, imp := range n.Imports {
				if imp.Name != nil && imp.Name.Name != "_" {
					pkgScopeIdent[imp.Name.Name] = imp
				} else {
					v, _ := strconv.Unquote(imp.Path.Value)
					if v != "_" {
						pkgScopeIdent[imports[v]] = imp
					}
				}
			}
			return true
		case *ast.FuncDecl:
			if _, ok := scopeIdent[n.Pos()]; !ok {
				scopeIdent[n.Pos()] = map[string]ast.Node{}
			}
			return true
		case *ast.TypeSpec:
			return true
		case *ast.ValueSpec:
			return true
		case *ast.AssignStmt:
			return n.Tok == token.DEFINE
		case *ast.Ident:
		IOUTER:
			for i := len(stack) - 1; i >= 0; i-- {
				pn := stack[i]
				switch pnt := pn.(type) {
				case *ast.File:
					pkgScopeIdent[n.Name] = n
				case *ast.FuncDecl:
					scopeIdent[pnt.Pos()][n.Name] = n
					break IOUTER
				}
			}
			return false
		}
		return false
	})

	dir, err := os.Getwd()
	if err != nil {
		dir = ""
	}
	for _, idents := range scopeIdent {
		for name, ident := range idents {
			if pkgIdent, ok := pkgScopeIdent[name]; ok {
				pos := pass.Fset.Position(pkgIdent.Pos())
				pos.Filename = strings.TrimPrefix(pos.Filename, dir+string(filepath.Separator))
				pass.Reportf(
					ident.Pos(),
					"conflict identifier name of '%s' by %s.",
					name,
					pos,
				)
			}
		}
	}

	return nil, nil
}
