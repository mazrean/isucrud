package dbdoc

import (
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
)

const (
	mermaidHeader = "# DB Graph\n" +
		"```mermaid\n" +
		"graph LR\n" +
		"  classDef func fill:" + funcNodeColor + ",fill-opacity:0.5\n" +
		"  classDef table fill:" + tableNodeColor + ",fill-opacity:0.5\n"
	mermaidFooter = "```"

	funcNodeColor   = "#1976D2"
	tableNodeColor  = "#795548"
	insertLinkColor = "#CDDC39"
	deleteLinkColor = "#F44336"
	selectLinkColor = "#78909C"
	updateLinkColor = "#FF9800"
	callLinkColor   = "#BBDEFB"
)

func writeMermaid(w io.StringWriter, nodes []*node) error {
	_, err := w.WriteString(mermaidHeader)
	if err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	edgeID := 0
	var insertLinks, deleteLinks, selectLinks, updateLinks, callLinks []string
	for _, node := range nodes {
		var src string
		switch node.nodeType {
		case nodeTypeTable:
			src = fmt.Sprintf("%s[%s]:::table", node.id, node.label)
		case nodeTypeFunction:
			src = fmt.Sprintf("%s[%s]:::func", node.id, node.label)
		default:
			log.Printf("unknown node type: %v\n", node.nodeType)
			src = fmt.Sprintf("%s[%s]", node.id, node.label)
		}

		for _, edge := range node.edges {
			var dst, line string
			switch edge.node.nodeType {
			case nodeTypeTable:
				dst = fmt.Sprintf("%s[%s]:::table", edge.node.id, edge.node.label)
			case nodeTypeFunction:
				dst = fmt.Sprintf("%s[%s]:::func", edge.node.id, edge.node.label)
			default:
				log.Printf("unknown node type: %v\n", edge.node.nodeType)
				dst = fmt.Sprintf("%s[%s]", edge.node.id, edge.node.label)
			}

			line = "--"

			if edge.label == "" {
				_, err = w.WriteString(fmt.Sprintf("  %s %s> %s\n", src, line, dst))
				if err != nil {
					return fmt.Errorf("failed to write edge: %w\n", err)
				}
			} else {
				_, err = w.WriteString(fmt.Sprintf("  %s %s %s %s> %s\n", src, line, edge.label, line, dst))
				if err != nil {
					return fmt.Errorf("failed to write edge: %w\n", err)
				}
			}

			switch edge.edgeType {
			case edgeTypeInsert:
				insertLinks = append(insertLinks, strconv.Itoa(edgeID))
			case edgeTypeDelete:
				deleteLinks = append(deleteLinks, strconv.Itoa(edgeID))
			case edgeTypeSelect:
				selectLinks = append(selectLinks, strconv.Itoa(edgeID))
			case edgeTypeUpdate:
				updateLinks = append(updateLinks, strconv.Itoa(edgeID))
			case edgeTypeCall:
				callLinks = append(callLinks, strconv.Itoa(edgeID))
			default:
				log.Printf("unknown edge type: %v\n", edge.edgeType)
			}

			edgeID++
		}
	}

	if len(insertLinks) > 0 {
		_, err = w.WriteString(fmt.Sprintf("  linkStyle %s stroke:%s,stroke-width:2px\n", strings.Join(insertLinks, ","), insertLinkColor))
		if err != nil {
			return fmt.Errorf("failed to write link style: %w\n", err)
		}
	}
	if len(deleteLinks) > 0 {
		_, err = w.WriteString(fmt.Sprintf("  linkStyle %s stroke:%s,stroke-width:2px\n", strings.Join(deleteLinks, ","), deleteLinkColor))
		if err != nil {
			return fmt.Errorf("failed to write link style: %w\n", err)
		}
	}
	if len(selectLinks) > 0 {
		_, err = w.WriteString(fmt.Sprintf("  linkStyle %s stroke:%s,stroke-width:2px\n", strings.Join(selectLinks, ","), selectLinkColor))
		if err != nil {
			return fmt.Errorf("failed to write link style: %w\n", err)
		}
	}
	if len(updateLinks) > 0 {
		_, err = w.WriteString(fmt.Sprintf("  linkStyle %s stroke:%s,stroke-width:2px\n", strings.Join(updateLinks, ","), updateLinkColor))
		if err != nil {
			return fmt.Errorf("failed to write link style: %w\n", err)
		}
	}
	if len(callLinks) > 0 {
		_, err = w.WriteString(fmt.Sprintf("  linkStyle %s stroke:%s,stroke-width:2px\n", strings.Join(callLinks, ","), callLinkColor))
		if err != nil {
			return fmt.Errorf("failed to write link style: %w\n", err)
		}
	}

	_, err = w.WriteString(mermaidFooter)
	if err != nil {
		return fmt.Errorf("failed to write footer: %w", err)
	}

	return nil
}
