package dbdoc

import (
	"fmt"
	"go/ast"
	"go/parser"
)

func BuildLoopRangeMap(ctx *Context) (LoopRangeMap, error) {
	astPkgs, err := parser.ParseDir(ctx.FileSet, ctx.WorkDir, nil, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to parse dir: %w", err)
	}

	lrm := make(LoopRangeMap)
	for _, astPkg := range astPkgs {
		for _, astFile := range astPkg.Files {
			for _, decl := range astFile.Decls {
				if f, ok := decl.(*ast.FuncDecl); ok {
					v := &loopRangeVisitor{
						lr: nil,
					}
					ast.Walk(v, f)
					lrm[f.Name.Name] = v.lr
				}
			}
		}
	}

	return lrm, nil
}

type loopRangeVisitor struct {
	lr LoopRanges
}

func (v *loopRangeVisitor) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.ForStmt:
		v.lr = append(v.lr, LoopRange{
			start: n.Body.Lbrace,
			end:   n.Body.Rbrace,
		})
		return nil
	case *ast.RangeStmt:
		v.lr = append(v.lr, LoopRange{
			start: n.Body.Lbrace,
			end:   n.Body.Rbrace,
		})
		return nil
	}

	return v
}
