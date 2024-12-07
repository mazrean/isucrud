package ui

import (
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"

	"github.com/mazrean/isucrud/dbdoc"
)

const (
	funcNodeColor   = "1976D2"
	tableNodeColor  = "795548"
	insertLinkColor = "CDDC39"
	deleteLinkColor = "F44336"
	selectLinkColor = "78909C"
	updateLinkColor = "FF9800"
	callLinkColor   = "BBDEFB"
)

type NodeType struct {
	name  string
	Label string
	Color string
	valid bool
}

type EdgeType struct {
	Label string
	Color string
	valid bool
}

var (
	nodeTypes = []NodeType{
		dbdoc.NodeTypeTable:    {"table", "テーブル", tableNodeColor, true},
		dbdoc.NodeTypeFunction: {"func", "関数", funcNodeColor, true},
	}
	edgeTypes = []EdgeType{
		dbdoc.EdgeTypeInsert: {"INSERT", insertLinkColor, true},
		dbdoc.EdgeTypeUpdate: {"UPDATE", updateLinkColor, true},
		dbdoc.EdgeTypeDelete: {"DELETE", deleteLinkColor, true},
		dbdoc.EdgeTypeSelect: {"SELECT", selectLinkColor, true},
		dbdoc.EdgeTypeCall:   {"関数呼び出し", callLinkColor, true},
	}
)

func RenderMermaid(
	w io.StringWriter,
	nodes []*dbdoc.Node,
	isHttp bool,
) error {
	_, err := w.WriteString("graph LR\n")
	if err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	for _, nodeType := range nodeTypes {
		if nodeType.valid {
			_, err = w.WriteString(fmt.Sprintf("  classDef %s fill:#%s,fill-opacity:0.5\n", nodeType.name, nodeType.Color))
			if err != nil {
				return fmt.Errorf("failed to write class def: %w", err)
			}
		}
	}

	edgeLinksMap := map[dbdoc.EdgeType][]string{}
	edgeID := 0
	for _, node := range nodes {
		var src string
		if nodeType := nodeTypes[node.NodeType]; nodeType.valid {
			src = fmt.Sprintf("%s[%s]:::%s", node.ID, node.Label, nodeType.name)
		} else {
			log.Printf("unknown node type: %v\n", node.NodeType)
			continue
		}

		for _, edge := range node.Edges {
			var dst string
			if nodeType := nodeTypes[edge.Node.NodeType]; nodeType.valid {
				dst = fmt.Sprintf("%s[%s]:::%s", edge.Node.ID, edge.Node.Label, nodeType.name)
			} else {
				log.Printf("unknown node type: %v\n", node.NodeType)
				continue
			}

			line := "--"
			if edge.InLoop {
				line = "=="
			}

			var edgeExpr string
			if edge.Label == "" {
				edgeExpr = fmt.Sprintf("%s>", line)
			} else {
				edgeExpr = fmt.Sprintf("%s %s %s>", line, edge.Label, line)
			}
			_, err = w.WriteString(fmt.Sprintf("  %s %s %s\n", src, edgeExpr, dst))
			if err != nil {
				return fmt.Errorf("failed to write edge: %w", err)
			}

			edgeLinksMap[edge.EdgeType] = append(edgeLinksMap[edge.EdgeType], strconv.Itoa(edgeID))

			edgeID++
		}
	}

	if isHttp {
		for _, node := range nodes {
			_, err = w.WriteString(fmt.Sprintf("  click %s \"/?node=%s\"\n", node.ID, node.ID))
			if err != nil {
				return fmt.Errorf("failed to write click event: %w", err)
			}
		}
	}

	for edgeType, links := range edgeLinksMap {
		if len(links) == 0 {
			continue
		}
		if info := edgeTypes[edgeType]; info.valid {
			_, err = w.WriteString(fmt.Sprintf("  linkStyle %s stroke:#%s,stroke-width:2px\n", strings.Join(links, ","), info.Color))
			if err != nil {
				return fmt.Errorf("failed to write link style: %w", err)
			}
		} else {
			log.Printf("unknown edge type: %v\n", edgeType)
		}
	}

	return nil
}
