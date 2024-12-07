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
)

var (
	//go:embed asset/template.md
	templateMarkdown string
	//go:embed asset/template.html
	templateHTML string
)

type TemplateParam struct {
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
		nil,
		false,
	)
	if err != nil {
		return fmt.Errorf("failed to write mermaid: %w", err)
	}

	tmpl, err := template.New("markdown").Parse(templateMarkdown)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	err = tmpl.Execute(f, TemplateParam{
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

func RenderHTML(w io.Writer, nodes []*dbdoc.Node) error {

	sb := &strings.Builder{}
	err := RenderMermaid(
		sb,
		nodes,
		nil,
		true,
	)
	if err != nil {
		return fmt.Errorf("failed to write mermaid: %w", err)
	}

	tmpl, err := htmlTemplate.New("html").Parse(templateHTML)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	err = tmpl.Execute(w, TemplateParam{
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
