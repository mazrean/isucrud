package dbdoc

import (
	"go/token"
)

type Context struct {
	FileSet *token.FileSet
	WorkDir string
}

type LoopRange struct {
	start token.Pos
	end   token.Pos
}
type LoopRanges []LoopRange
type LoopRangeMap map[string]LoopRanges

func (lr LoopRanges) Search(fset *token.FileSet, pos token.Pos) bool {
	position := fset.Position(pos)
	for _, r := range lr {
		start := fset.Position(r.start)
		end := fset.Position(r.end)
		if position.Filename != start.Filename || position.Filename != end.Filename ||
			position.Line < start.Line || position.Line > end.Line ||
			(position.Line == start.Line && position.Column < start.Column) ||
			(position.Line == end.Line && position.Column > end.Column) {
			continue
		}

		return true
	}

	return false
}

type inLoop[T any] struct {
	value  T
	inLoop bool
}

type function struct {
	id      string
	name    string
	queries []inLoop[query]
	calls   []inLoop[string]
}

type stringLiteral struct {
	value string
	pos   token.Pos
}

type query struct {
	queryType queryType
	table     string
	pos       token.Pos
}

type queryType uint8

const (
	queryTypeSelect queryType = iota + 1
	queryTypeInsert
	queryTypeUpdate
	queryTypeDelete
)

func (qt queryType) String() string {
	switch qt {
	case queryTypeSelect:
		return "select"
	case queryTypeInsert:
		return "insert"
	case queryTypeUpdate:
		return "update"
	case queryTypeDelete:
		return "delete"
	}

	return ""
}

type node struct {
	id       string
	label    string
	nodeType nodeType
	edges    []edge
}

type nodeType uint8

const (
	nodeTypeUnknown nodeType = iota
	nodeTypeTable
	nodeTypeFunction
)

type edge struct {
	label    string
	node     *node
	edgeType edgeType
	inLoop   bool
}

type edgeType uint8

const (
	edgeTypeUnknown edgeType = iota
	edgeTypeInsert
	edgeTypeUpdate
	edgeTypeDelete
	edgeTypeSelect
	edgeTypeCall
)
