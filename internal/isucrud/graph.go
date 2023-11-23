package isucrud

import (
	"container/list"
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/mazrean/isucrud/internal/pkg/analyze"
)

func buildGraph(funcs []function, ignoreFuncs, ignoreFuncPrefixes []string) []*node {
	type tmpEdge struct {
		label    string
		edgeType edgeType
		childID  string
	}
	type tmpNode struct {
		*node
		edges []tmpEdge
	}
	tmpNodeMap := make(map[string]tmpNode, len(funcs))
FUNC_LOOP:
	for _, f := range funcs {
		if f.name == "main" || analyze.IsInitializeFuncName(f.name) {
			continue
		}

		for _, ignore := range ignoreFuncs {
			if f.name == ignore {
				continue FUNC_LOOP
			}
		}

		for _, ignorePrefix := range ignoreFuncPrefixes {
			if strings.HasPrefix(f.name, ignorePrefix) {
				continue FUNC_LOOP
			}
		}

		var edges []tmpEdge
		for _, q := range f.queries {
			id := tableID(q.table)
			tmpNodeMap[id] = tmpNode{
				node: &node{
					id:       id,
					label:    q.table,
					nodeType: nodeTypeTable,
				},
			}

			var edgeType edgeType
			switch q.queryType {
			case queryTypeSelect:
				edgeType = edgeTypeSelect
			case queryTypeInsert:
				edgeType = edgeTypeInsert
			case queryTypeUpdate:
				edgeType = edgeTypeUpdate
			case queryTypeDelete:
				edgeType = edgeTypeDelete
			default:
				log.Printf("unknown query type: %v\n", q.queryType)
				continue
			}

			edges = append(edges, tmpEdge{
				label:    "",
				edgeType: edgeType,
				childID:  tableID(q.table),
			})
		}

		for _, c := range f.calls {
			id := funcID(c)
			edges = append(edges, tmpEdge{
				label:    "",
				edgeType: edgeTypeCall,
				childID:  id,
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

		id := funcID(f.id)
		tmpNodeMap[id] = tmpNode{
			node: &node{
				id:       id,
				label:    f.name,
				nodeType: nodeTypeFunction,
			},
			edges: edges,
		}
	}

	type revEdge struct {
		label    string
		edgeType edgeType
		parentID string
	}
	revEdgeMap := make(map[string][]revEdge)
	for _, tmpNode := range tmpNodeMap {
		for _, tmpEdge := range tmpNode.edges {
			revEdgeMap[tmpEdge.childID] = append(revEdgeMap[tmpEdge.childID], revEdge{
				label:    tmpEdge.label,
				edgeType: tmpEdge.edgeType,
				parentID: tmpNode.id,
			})
		}
	}

	newNodeMap := make(map[string]tmpNode, len(tmpNodeMap))
	nodeQueue := list.New()
	for id, node := range tmpNodeMap {
		if node.nodeType == nodeTypeTable {
			newNodeMap[id] = node
			nodeQueue.PushBack(node)
			delete(tmpNodeMap, id)
			continue
		}
	}

	for {
		element := nodeQueue.Front()
		if element == nil {
			break
		}
		nodeQueue.Remove(element)

		node := element.Value.(tmpNode)
		for _, edge := range revEdgeMap[node.id] {
			parent := tmpNodeMap[edge.parentID]
			newNodeMap[edge.parentID] = parent
			nodeQueue.PushBack(parent)
		}
		delete(revEdgeMap, node.id)
	}

	var nodes []*node
	for _, tmpNode := range newNodeMap {
		node := tmpNode.node
		for _, tmpEdge := range tmpNode.edges {
			child, ok := newNodeMap[tmpEdge.childID]
			if !ok {
				continue
			}

			node.edges = append(node.edges, edge{
				label:    tmpEdge.label,
				node:     child.node,
				edgeType: tmpEdge.edgeType,
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
