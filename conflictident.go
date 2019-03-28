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
		(*ast.ImportSpec)(nil),
		(*ast.SelectorExpr)(nil),
		(*ast.CompositeLit)(nil),
		(*ast.CallExpr)(nil),
		(*ast.ReturnStmt)(nil),
		(*ast.ExprStmt)(nil),
		(*ast.BinaryExpr)(nil),
		(*ast.Field)(nil),
		(*ast.Ident)(nil),
	}

	imports := map[string]string{}
	for _, imp := range pass.Pkg.Imports() {
		imports[imp.Path()] = imp.Name()
	}

	pkgScopeIdent := map[string]ast.Node{}
	fileScopeIdent := map[token.Pos]map[string]ast.Node{}
	scopeIdent := map[token.Pos]map[string]ast.Node{}
	inspect.WithStack(nodeFilter, func(n ast.Node, push bool, stack []ast.Node) bool {
		switch n := n.(type) {
		case *ast.File:
			fileScopeIdent[n.Pos()] = map[string]ast.Node{}
			for _, imp := range n.Imports {
				if imp.Name != nil && imp.Name.Name != "_" {
					fileScopeIdent[n.Pos()][imp.Name.Name] = imp
				} else {
					v, _ := strconv.Unquote(imp.Path.Value)
					if v != "_" {
						fileScopeIdent[n.Pos()][imports[v]] = imp
					}
				}
			}
			return true
		case *ast.FuncDecl:
			if _, ok := scopeIdent[n.Pos()]; !ok {
				scopeIdent[n.Pos()] = map[string]ast.Node{}
			}
			// receiver
			if n.Recv != nil {
				for _, field := range n.Recv.List {
					for _, name := range field.Names {
						scopeIdent[n.Pos()][name.Name] = name
					}
				}
			}
			// arguments
			for _, field := range n.Type.Params.List {
				for _, name := range field.Names {
					scopeIdent[n.Pos()][name.Name] = name
				}
			}
			return true
		case *ast.TypeSpec:
			dependIdent(n.Name, stack, pkgScopeIdent, scopeIdent)
			return false
		case *ast.ValueSpec:
			for _, name := range n.Names {
				dependIdent(name, stack, pkgScopeIdent, scopeIdent)
			}
			return false
		case *ast.AssignStmt:
			if n.Tok != token.DEFINE {
				return false
			}
			for _, expr := range n.Lhs {
				if idt, ok := expr.(*ast.Ident); ok {
					dependIdent(idt, stack, pkgScopeIdent, scopeIdent)
				}
			}
		case
			*ast.CallExpr,
			*ast.ImportSpec,
			*ast.SelectorExpr,
			*ast.Field,
			*ast.ReturnStmt,
			*ast.ExprStmt,
			*ast.CompositeLit:
			return false
		case *ast.Ident:
			if !push {
				return false
			}
			switch stack[len(stack)-2].(type) {
			case *ast.CaseClause, *ast.IfStmt, *ast.BinaryExpr:
				return false
			}
			dependIdent(n, stack, pkgScopeIdent, scopeIdent)
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
			var hasConflict bool
			var conflictPos token.Pos
			if pkgIdent, ok := pkgScopeIdent[name]; ok {
				hasConflict = true
				conflictPos = pkgIdent.Pos()
			}
			f := pass.Fset.File(ident.Pos())
			if fileIdent, ok := fileScopeIdent[f.Pos(0)][name]; ok {
				hasConflict = true
				conflictPos = fileIdent.Pos()
			}
			if hasConflict {
				pos := pass.Fset.Position(conflictPos)
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

func dependIdent(n *ast.Ident, stack []ast.Node, pkgScopeIdent map[string]ast.Node, scopeIdent map[token.Pos]map[string]ast.Node) {
	for i := len(stack) - 1; i >= 0; i-- {
		pn := stack[i]
		switch pnt := pn.(type) {
		case *ast.File:
			pkgScopeIdent[n.Name] = n
		case *ast.FuncDecl:
			scopeIdent[pnt.Pos()][n.Name] = n
			return
		}
	}
}
