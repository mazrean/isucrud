package ui

import (
	_ "embed"
	"fmt"
	htmlTemplate "html/template"
	"io"
	"os"
	"strings"
	"text/template"

	"github.com/mazrean/isucrud/dbdoc"
	"github.com/mazrean/isucrud/internal/pkg/list"
)

var (
	//go:embed asset/template.md
	templateMarkdown string
	//go:embed asset/template.html
	templateHTML string
)

type TemplateParam struct {
	IsFiltered  bool
	BasePath    string
	NodeTypes   []NodeType
	EdgeTypes   []EdgeType
	Nodes       []*dbdoc.Node
	MermaidData string
}

func RenderMarkdown(dest string, nodes []*dbdoc.Node) error {
	f, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("failed to make directory: %w", err)
	}
	defer f.Close()

	sb := &strings.Builder{}
	err = RenderMermaid(
		sb,
		nodes,
		RenderMermaidOption{
			IsHttp: false,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to write mermaid: %w", err)
	}

	tmpl, err := template.New("markdown").Parse(templateMarkdown)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	err = tmpl.Execute(f, TemplateParam{
		IsFiltered:  false,
		NodeTypes:   nodeTypes[1:],
		EdgeTypes:   edgeTypes[1:],
		Nodes:       nodes,
		MermaidData: sb.String(),
	})
	if err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}

func RenderHTML(w io.Writer, nodes []*dbdoc.Node, targetNodeID string, basePath string) error {
	filtered := false
	filteredNodes := nodes
	if targetNodeID != "" {
		filtered = true
		filteredNodes = filterNodes(targetNodeID, nodes)
	}

	sb := &strings.Builder{}
	err := RenderMermaid(
		sb,
		filteredNodes,
		RenderMermaidOption{
			IsHttp:   true,
			BasePath: basePath,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to write mermaid: %w", err)
	}

	tmpl, err := htmlTemplate.New("html").Parse(templateHTML)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	err = tmpl.Execute(w, TemplateParam{
		IsFiltered:  filtered,
		BasePath:    basePath,
		NodeTypes:   nodeTypes[1:],
		EdgeTypes:   edgeTypes[1:],
		Nodes:       nodes,
		MermaidData: sb.String(),
	})
	if err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}

func filterNodes(targetNodeID string, nodes []*dbdoc.Node) []*dbdoc.Node {
	nodeMap := make(map[string]*dbdoc.Node)
	for _, node := range nodes {
		nodeMap[node.ID] = node
	}

	targetNode, exists := nodeMap[targetNodeID]
	if !exists {
		return []*dbdoc.Node{}
	}

	parentMap := make(map[string][]*dbdoc.Node)
	for _, node := range nodes {
		for _, edge := range node.Edges {
			if edge.Node != nil {
				parentMap[edge.Node.ID] = append(parentMap[edge.Node.ID], node)
			}
		}
	}

	visited := make(map[string]struct{}, len(nodes))
	relatedNodeMap := make(map[string]*dbdoc.Node)

	childNodeQueue := list.NewQueue[*dbdoc.Node]()
	childNodeQueue.Push(targetNode)
	for current := range childNodeQueue.Iter {
		// Skip if already visited.
		if _, ok := visited[current.ID]; ok {
			continue
		}
		visited[current.ID] = struct{}{}

		relatedNodeMap[current.ID] = current

		for _, edge := range current.Edges {
			childNodeQueue.Push(edge.Node)
		}
	}

	visited = make(map[string]struct{}, len(nodes))
	parentNodeQueue := list.NewQueue[*dbdoc.Node]()
	parentNodeQueue.Push(targetNode)
	for current := range parentNodeQueue.Iter {
		// Skip if already visited.
		if _, ok := visited[current.ID]; ok {
			continue
		}
		visited[current.ID] = struct{}{}

		relatedNodeMap[current.ID] = current

		for _, parent := range parentMap[current.ID] {
			parentNodeQueue.Push(parent)
		}
	}

	relatedNodes := make([]*dbdoc.Node, 0, len(relatedNodeMap))
	for _, node := range relatedNodeMap {
		newEdges := []dbdoc.Edge{}
		for _, edge := range node.Edges {
			if _, ok := relatedNodeMap[edge.Node.ID]; ok {
				newEdges = append(newEdges, edge)
			}
		}

		relatedNodes = append(relatedNodes, &dbdoc.Node{
			ID:       node.ID,
			Label:    node.Label,
			Edges:    newEdges,
			NodeType: node.NodeType,
		})
	}

	return relatedNodes
}
