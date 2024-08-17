package dbdoc

import (
	"go/constant"
	"go/token"
	"go/types"

	"github.com/mazrean/isucrud/internal/pkg/list"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"
)

func BuildFuncs(ctx *Context, pkgs []*packages.Package, ssaProgram *ssa.Program, loopRangeMap LoopRangeMap) ([]Function, error) {
	var funcs []Function
	for _, pkg := range pkgs {
		for _, def := range pkg.TypesInfo.Defs {
			if def == nil {
				continue
			}

			switch def := def.(type) {
			case *types.Func:
				ssaFunc := ssaProgram.FuncValue(def)
				if ssaFunc == nil {
					continue
				}

				stringLiterals, calls := analyzeFuncBody(ctx, ssaFunc.Blocks, getPos(ssaFunc.Pos(), def.Pos()))

				anonFuncQueue := list.NewQueue[*ssa.Function]()
				for _, anonFunc := range ssaFunc.AnonFuncs {
					anonFuncQueue.Push(anonFunc)
				}

				for anonFunc, ok := anonFuncQueue.Pop(); ok; anonFunc, ok = anonFuncQueue.Pop() {
					anonQueries, anonCalls := analyzeFuncBody(ctx, anonFunc.Blocks, getPos(anonFunc.Pos(), ssaFunc.Pos(), def.Pos()))
					stringLiterals = append(stringLiterals, anonQueries...)
					calls = append(calls, anonCalls...)

					for _, anonFunc := range anonFunc.AnonFuncs {
						anonFuncQueue.Push(anonFunc)
					}
				}

				if len(stringLiterals) == 0 && len(calls) == 0 {
					continue
				}

				queries := make([]Query, 0, len(stringLiterals))
				for _, strLiteral := range stringLiterals {
					newQueries := AnalyzeSQL(ctx, strLiteral)
					queries = append(queries, newQueries...)
				}

				loopRanges := loopRangeMap[ssaFunc.Name()]
				queriesInLoop := make([]InLoop[Query], 0, len(queries))
				for _, q := range queries {
					queriesInLoop = append(queriesInLoop, InLoop[Query]{
						Value:  q,
						InLoop: loopRanges.Search(ctx.FileSet, q.Pos),
					})
				}

				callsInLoop := make([]InLoop[Call], 0, len(calls))
				for _, call := range calls {
					callsInLoop = append(callsInLoop, InLoop[Call]{
						Value: Call{
							FunctionID: call.id,
							Pos:        call.pos,
						},
						InLoop: loopRanges.Search(ctx.FileSet, call.pos),
					})
				}

				funcs = append(funcs, Function{
					ID:      def.Id(),
					Name:    def.Name(),
					Pos:     def.Pos(),
					Queries: queriesInLoop,
					Calls:   callsInLoop,
				})
			}
		}
	}

	return funcs, nil
}

type funcCall struct {
	id  string
	pos token.Pos
}

func analyzeFuncBody(ctx *Context, blocks []*ssa.BasicBlock, pos token.Pos) ([]stringLiteral, []funcCall) {
	type ssaValue struct {
		value ssa.Value
		pos   token.Pos
	}
	var ssaValues []ssaValue
	var calls []funcCall
	for _, block := range blocks {
		for _, instr := range block.Instrs {
			switch instr := instr.(type) {
			case *ssa.BinOp:
				if instr.X != nil {
					ssaValues = append(ssaValues, ssaValue{
						value: instr.X,
						pos:   getPos(instr.X.Pos(), instr.Pos(), pos),
					})
				}

				if instr.Y != nil {
					ssaValues = append(ssaValues, ssaValue{
						value: instr.Y,
						pos:   getPos(instr.Y.Pos(), instr.Pos(), pos),
					})
				}
			case *ssa.ChangeType:
				if instr.X != nil {
					ssaValues = append(ssaValues, ssaValue{
						value: instr.X,
						pos:   getPos(instr.X.Pos(), instr.Pos(), pos),
					})
				}
			case *ssa.Convert:
				if instr.X != nil {
					ssaValues = append(ssaValues, ssaValue{
						value: instr.X,
						pos:   getPos(instr.X.Pos(), instr.Pos(), pos),
					})
				}
			case *ssa.MakeClosure:
				for _, bind := range instr.Bindings {
					if bind != nil {
						ssaValues = append(ssaValues, ssaValue{
							value: bind,
							pos:   getPos(bind.Pos(), instr.Pos(), pos),
						})
					}
				}
			case *ssa.MultiConvert:
				if instr.X != nil {
					ssaValues = append(ssaValues, ssaValue{
						value: instr.X,
						pos:   getPos(instr.X.Pos(), instr.Pos(), pos),
					})
				}
			case *ssa.Store:
				if instr.Val != nil {
					ssaValues = append(ssaValues, ssaValue{
						value: instr.Val,
						pos:   getPos(instr.Val.Pos(), instr.Pos(), pos),
					})
				}
			case *ssa.Call:
				if f, ok := instr.Call.Value.(*ssa.Function); ok {
					if f.Object() == nil {
						continue
					}
					calls = append(calls, funcCall{
						id:  f.Object().Id(),
						pos: getPos(f.Pos(), instr.Pos(), pos),
					})
				}

				for _, arg := range instr.Call.Args {
					if arg != nil {
						ssaValues = append(ssaValues, ssaValue{
							value: arg,
							pos:   getPos(arg.Pos(), instr.Pos(), pos),
						})
					}
				}
			case *ssa.Defer:
				if f, ok := instr.Call.Value.(*ssa.Function); ok {
					if f.Object() == nil {
						continue
					}
					calls = append(calls, funcCall{
						id:  f.Object().Id(),
						pos: getPos(instr.Call.Pos(), instr.Pos(), pos),
					})
				}

				for _, arg := range instr.Call.Args {
					if arg != nil {
						ssaValues = append(ssaValues, ssaValue{
							value: arg,
							pos:   getPos(arg.Pos(), instr.Pos(), pos),
						})
					}
				}
			case *ssa.Go:
				if f, ok := instr.Call.Value.(*ssa.Function); ok {
					if f.Object() == nil {
						continue
					}
					calls = append(calls, funcCall{
						id:  f.Object().Id(),
						pos: getPos(instr.Call.Pos(), instr.Pos(), pos),
					})
				}

				for _, arg := range instr.Call.Args {
					if arg != nil {
						ssaValues = append(ssaValues, ssaValue{
							value: arg,
							pos:   getPos(arg.Pos(), instr.Pos(), pos),
						})
					}
				}
			}
		}
	}

	queries := make([]stringLiteral, 0, len(ssaValues))
	for _, ssaValue := range ssaValues {
		strValue, ok := checkValue(ctx, ssaValue.value)
		if ok {
			queries = append(queries, stringLiteral{
				value: strValue,
				pos:   ssaValue.pos,
			})
		}
	}

	return queries, calls
}

func getPos(posList ...token.Pos) token.Pos {
	for _, pos := range posList {
		if pos.IsValid() {
			return pos
		}
	}

	return token.NoPos
}

func checkValue(_ *Context, v ssa.Value) (string, bool) {
	constValue, ok := v.(*ssa.Const)
	if !ok || constValue == nil || constValue.Value == nil {
		return "", false
	}

	if constValue.Value.Kind() != constant.String {
		return "", false
	}

	return constant.StringVal(constValue.Value), true
}
