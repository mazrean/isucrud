package dbdoc

import (
	"fmt"
	"go/token"

	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

type Config struct {
	WorkDir                      string
	BuildArgs                    []string
	IgnoreFuncs                  []string
	IgnoreFuncPrefixes           []string
	IgnoreMain, IgnoreInitialize bool
	DestinationFilePath          string
}

func Run(conf Config) ([]*Node, error) {
	ctx := &Context{
		FileSet: token.NewFileSet(),
		WorkDir: conf.WorkDir,
	}

	ssaProgram, pkgs, err := BuildSSA(ctx, conf.BuildArgs)
	if err != nil {
		return nil, fmt.Errorf("failed to build ssa: %w", err)
	}

	loopRangeMap, err := BuildLoopRangeMap(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to build loop range map: %w", err)
	}

	funcs, err := BuildFuncs(ctx, pkgs, ssaProgram, loopRangeMap)
	if err != nil {
		return nil, fmt.Errorf("failed to build funcs: %w", err)
	}

	nodes := BuildGraph(
		funcs,
		conf.IgnoreFuncs, conf.IgnoreFuncPrefixes,
		conf.IgnoreMain, conf.IgnoreInitialize,
	)

	return nodes, nil
}

func BuildSSA(ctx *Context, args []string) (*ssa.Program, []*packages.Package, error) {
	pkgs, err := packages.Load(&packages.Config{
		Fset: ctx.FileSet,
		Mode: packages.NeedFiles | packages.NeedSyntax | packages.NeedTypes | packages.NeedImports | packages.NeedTypesInfo | packages.NeedName | packages.NeedModule,
	}, args...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load packages: %w", err)
	}

	ssaProgram, _ := ssautil.AllPackages(pkgs, ssa.BareInits)
	ssaProgram.Build()

	return ssaProgram, pkgs, nil
}
