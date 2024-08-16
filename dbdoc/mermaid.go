package dbdoc

import (
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
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

var (
	nodeTypes = []struct {
		name  string
		label string
		color string
		valid bool
	}{
		nodeTypeTable:    {"table", "テーブル", tableNodeColor, true},
		nodeTypeFunction: {"func", "関数", funcNodeColor, true},
	}
	edgeTypes = []struct {
		label string
		color string
		valid bool
	}{
		edgeTypeInsert: {"INSERT", insertLinkColor, true},
		edgeTypeUpdate: {"UPDATE", updateLinkColor, true},
		edgeTypeDelete: {"DELETE", deleteLinkColor, true},
		edgeTypeSelect: {"SELECT", selectLinkColor, true},
		edgeTypeCall:   {"関数呼び出し", callLinkColor, true},
	}
)

func writeMermaid(w io.StringWriter, nodes []*node) error {
	_, err := w.WriteString("# DB Graph\n")
	if err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	_, err = w.WriteString("node: ")
	if err != nil {
		return fmt.Errorf("failed to write node description start: %w", err)
	}

	for _, nodeType := range nodeTypes {
		if nodeType.valid {
			_, err = w.WriteString(fmt.Sprintf("![](https://via.placeholder.com/16/%s/FFFFFF/?text=%%20) `%s` ", nodeType.color, nodeType.label))
			if err != nil {
				return fmt.Errorf("failed to write node description: %w", err)
			}
		}
	}

	_, err = w.WriteString("\n\n")
	if err != nil {
		return fmt.Errorf("failed to write node description end: %w", err)
	}

	_, err = w.WriteString("edge: ")
	if err != nil {
		return fmt.Errorf("failed to write edge description start: %w", err)
	}

	for _, edgeType := range edgeTypes {
		if edgeType.valid {
			_, err = w.WriteString(fmt.Sprintf("![](https://via.placeholder.com/16/%s/FFFFFF/?text=%%20) `%s` ", edgeType.color, edgeType.label))
			if err != nil {
				return fmt.Errorf("failed to write edge description: %w", err)
			}
		}
	}

	_, err = w.WriteString("\n")
	if err != nil {
		return fmt.Errorf("failed to write edge description end: %w", err)
	}

	_, err = w.WriteString("```mermaid\n" +
		"graph LR\n")
	if err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	for _, nodeType := range nodeTypes {
		if nodeType.valid {
			_, err = w.WriteString(fmt.Sprintf("  classDef %s fill:#%s,fill-opacity:0.5\n", nodeType.name, nodeType.color))
			if err != nil {
				return fmt.Errorf("failed to write class def: %w", err)
			}
		}
	}

	edgeLinksMap := map[edgeType][]string{}
	edgeID := 0
	for _, node := range nodes {
		var src string
		if nodeType := nodeTypes[node.nodeType]; nodeType.valid {
			src = fmt.Sprintf("%s[%s]:::%s", node.id, node.label, nodeType.name)
		} else {
			log.Printf("unknown node type: %v\n", node.nodeType)
			continue
		}

		for _, edge := range node.edges {
			var dst string
			if nodeType := nodeTypes[edge.node.nodeType]; nodeType.valid {
				dst = fmt.Sprintf("%s[%s]:::%s", edge.node.id, edge.node.label, nodeType.name)
			} else {
				log.Printf("unknown node type: %v\n", node.nodeType)
				continue
			}

			line := "--"

			var edgeExpr string
			if edge.label == "" {
				edgeExpr = fmt.Sprintf("%s>", line)
			} else {
				edgeExpr = fmt.Sprintf("%s %s %s>", line, edge.label, line)
			}
			_, err = w.WriteString(fmt.Sprintf("  %s %s %s\n", src, edgeExpr, dst))
			if err != nil {
				return fmt.Errorf("failed to write edge: %w\n", err)
			}

			edgeLinksMap[edge.edgeType] = append(edgeLinksMap[edge.edgeType], strconv.Itoa(edgeID))

			edgeID++
		}
	}

	for edgeType, links := range edgeLinksMap {
		if len(links) == 0 {
			continue
		}
		if info := edgeTypes[edgeType]; info.valid {
			_, err = w.WriteString(fmt.Sprintf("  linkStyle %s stroke:#%s,stroke-width:2px\n", strings.Join(links, ","), info.color))
			if err != nil {
				return fmt.Errorf("failed to write link style: %w\n", err)
			}
		} else {
			log.Printf("unknown edge type: %v\n", edgeType)
		}
	}

	_, err = w.WriteString("```")
	if err != nil {
		return fmt.Errorf("failed to write footer: %w", err)
	}

	return nil
}
