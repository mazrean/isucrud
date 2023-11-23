package dbdoc

import (
	"fmt"
	"go/token"
	"os"

	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

type Config struct {
	WorkDir             string
	BuildArgs           []string
	IgnoreFuncs         []string
	IgnoreFuncPrefixes  []string
	DestinationFilePath string
}

func Run(conf Config) error {
	ctx := &context{
		fileSet: token.NewFileSet(),
		workDir: conf.WorkDir,
	}

	ssaProgram, pkgs, err := buildSSA(ctx, conf.BuildArgs)
	if err != nil {
		return fmt.Errorf("failed to build ssa: %w", err)
	}

	funcs, err := buildFuncs(ctx, pkgs, ssaProgram)
	if err != nil {
		return fmt.Errorf("failed to build funcs: %w", err)
	}

	nodes := buildGraph(funcs, conf.IgnoreFuncs, conf.IgnoreFuncPrefixes)

	f, err := os.Create(conf.DestinationFilePath)
	if err != nil {
		return fmt.Errorf("failed to make directory: %w", err)
	}
	defer f.Close()

	err = writeMermaid(f, nodes)
	if err != nil {
		return fmt.Errorf("failed to write mermaid: %w", err)
	}

	return nil
}

func buildSSA(ctx *context, args []string) (*ssa.Program, []*packages.Package, error) {
	pkgs, err := packages.Load(&packages.Config{
		Fset: ctx.fileSet,
		Mode: packages.NeedFiles | packages.NeedSyntax | packages.NeedTypes | packages.NeedImports | packages.NeedTypesInfo | packages.NeedName | packages.NeedModule,
	}, args...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load packages: %w", err)
	}

	ssaProgram, _ := ssautil.AllPackages(pkgs, ssa.BareInits)
	ssaProgram.Build()

	return ssaProgram, pkgs, nil
}
