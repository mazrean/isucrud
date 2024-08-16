package dbdoc

import (
	"go/token"
)

type context struct {
	fileSet *token.FileSet
	workDir string
}

type function struct {
	id      string
	name    string
	queries []query
	calls   []string
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
