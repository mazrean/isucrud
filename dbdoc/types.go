package dbdoc

import (
	"go/token"
)

type Context struct {
	FileSet *token.FileSet
	WorkDir string
}

type LoopRange struct {
	Start token.Pos
	End   token.Pos
}
type LoopRanges []LoopRange
type LoopRangeMap map[string]LoopRanges

func (lr LoopRanges) Search(fset *token.FileSet, pos token.Pos) bool {
	position := fset.Position(pos)
	for _, r := range lr {
		start := fset.Position(r.Start)
		end := fset.Position(r.End)
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

type InLoop[T any] struct {
	Value  T
	InLoop bool
}

type Function struct {
	ID      string
	Name    string
	Pos     token.Pos
	Queries []InLoop[Query]
	Calls   []InLoop[Call]
}

type stringLiteral struct {
	value string
	pos   token.Pos
}

type Call struct {
	FunctionID string
	Pos        token.Pos
}

type Query struct {
	QueryType QueryType
	Table     string
	Pos       token.Pos
}

type QueryType uint8

const (
	QueryTypeSelect QueryType = iota + 1
	QueryTypeInsert
	QueryTypeUpdate
	QueryTypeDelete
)

func (qt QueryType) String() string {
	switch qt {
	case QueryTypeSelect:
		return "select"
	case QueryTypeInsert:
		return "insert"
	case QueryTypeUpdate:
		return "update"
	case QueryTypeDelete:
		return "delete"
	}

	return ""
}

type Node struct {
	ID       string
	Label    string
	NodeType NodeType
	Edges    []Edge
}

type NodeType uint8

const (
	NodeTypeUnknown NodeType = iota
	NodeTypeTable
	NodeTypeFunction
)

type Edge struct {
	Label    string
	Node     *Node
	EdgeType EdgeType
	InLoop   bool
}

type EdgeType uint8

const (
	EdgeTypeUnknown EdgeType = iota
	EdgeTypeInsert
	EdgeTypeUpdate
	EdgeTypeDelete
	EdgeTypeSelect
	EdgeTypeCall
)
