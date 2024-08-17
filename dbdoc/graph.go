package dbdoc

import (
	"container/list"
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/mazrean/isucrud/internal/pkg/analyze"
)

func BuildGraph(funcs []Function, ignoreFuncs, ignoreFuncPrefixes []string, ignoreMain, ignoreInitialize bool) []*Node {
	type tmpEdge struct {
		label    string
		edgeType EdgeType
		childID  string
		inLoop   bool
	}
	type tmpNode struct {
		*Node
		edges []tmpEdge
	}
	tmpNodeMap := make(map[string]tmpNode, len(funcs))
FUNC_LOOP:
	for _, f := range funcs {
		if (ignoreMain && f.Name == "main") ||
			(ignoreInitialize && analyze.IsInitializeFuncName(f.Name)) {
			continue
		}

		for _, ignore := range ignoreFuncs {
			if f.Name == ignore {
				continue FUNC_LOOP
			}
		}

		for _, ignorePrefix := range ignoreFuncPrefixes {
			if strings.HasPrefix(f.Name, ignorePrefix) {
				continue FUNC_LOOP
			}
		}

		var edges []tmpEdge
		for _, q := range f.Queries {
			id := tableID(q.Value.Table)
			tmpNodeMap[id] = tmpNode{
				Node: &Node{
					ID:       id,
					Label:    q.Value.Table,
					NodeType: NodeTypeTable,
				},
			}

			var edgeType EdgeType
			switch q.Value.QueryType {
			case QueryTypeSelect:
				edgeType = EdgeTypeSelect
			case QueryTypeInsert:
				edgeType = EdgeTypeInsert
			case QueryTypeUpdate:
				edgeType = EdgeTypeUpdate
			case QueryTypeDelete:
				edgeType = EdgeTypeDelete
			default:
				log.Printf("unknown query type: %v\n", q.Value.QueryType)
				continue
			}

			edges = append(edges, tmpEdge{
				label:    "",
				edgeType: edgeType,
				childID:  tableID(q.Value.Table),
				inLoop:   q.InLoop,
			})
		}

		for _, c := range f.Calls {
			id := funcID(c.Value.FunctionID)
			edges = append(edges, tmpEdge{
				label:    "",
				edgeType: EdgeTypeCall,
				childID:  id,
				inLoop:   c.InLoop,
			})
		}

		slices.SortFunc(edges, func(a, b tmpEdge) int {
			switch {
			case a.childID < b.childID:
				return -1
			case a.childID > b.childID:
				return 1
			default:
				return 0
			}
		})
		edges = slices.Compact(edges)

		id := funcID(f.ID)
		tmpNodeMap[id] = tmpNode{
			Node: &Node{
				ID:       id,
				Label:    f.Name,
				NodeType: NodeTypeFunction,
			},
			edges: edges,
		}
	}

	type revEdge struct {
		label    string
		edgeType EdgeType
		parentID string
		inLoop   bool
	}
	revEdgeMap := make(map[string][]revEdge)
	for _, tmpNode := range tmpNodeMap {
		for _, tmpEdge := range tmpNode.edges {
			revEdgeMap[tmpEdge.childID] = append(revEdgeMap[tmpEdge.childID], revEdge{
				label:    tmpEdge.label,
				edgeType: tmpEdge.edgeType,
				parentID: tmpNode.ID,
				inLoop:   tmpEdge.inLoop,
			})
		}
	}

	newNodeMap := make(map[string]tmpNode, len(tmpNodeMap))
	nodeQueue := list.New()
	for id, node := range tmpNodeMap {
		if node.NodeType == NodeTypeTable {
			newNodeMap[id] = node
			nodeQueue.PushBack(node)
			delete(tmpNodeMap, id)
			continue
		}
	}

	for element := nodeQueue.Front(); element != nil; element = nodeQueue.Front() {
		nodeQueue.Remove(element)

		node := element.Value.(tmpNode)
		for _, edge := range revEdgeMap[node.ID] {
			parent := tmpNodeMap[edge.parentID]
			newNodeMap[edge.parentID] = parent
			nodeQueue.PushBack(parent)
		}
		delete(revEdgeMap, node.ID)
	}

	var nodes []*Node
	for _, tmpNode := range newNodeMap {
		node := tmpNode.Node
		for _, tmpEdge := range tmpNode.edges {
			child, ok := newNodeMap[tmpEdge.childID]
			if !ok {
				continue
			}

			node.Edges = append(node.Edges, Edge{
				Label:    tmpEdge.label,
				Node:     child.Node,
				EdgeType: tmpEdge.edgeType,
				InLoop:   tmpEdge.inLoop,
			})
		}
		nodes = append(nodes, node)
	}

	return nodes
}

func funcID(functionID string) string {
	functionID = strings.Replace(functionID, "(", "", -1)
	functionID = strings.Replace(functionID, ")", "", -1)
	functionID = strings.Replace(functionID, "[", "", -1)
	functionID = strings.Replace(functionID, "]", "", -1)

	return fmt.Sprintf("func:%s", functionID)
}

func tableID(table string) string {
	table = strings.Replace(table, "(", "", -1)
	table = strings.Replace(table, ")", "", -1)
	table = strings.Replace(table, "[", "", -1)
	table = strings.Replace(table, "]", "", -1)

	return fmt.Sprintf("table:%s", table)
}
