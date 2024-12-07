# DB Graph
node: {{range .NodeTypes}}![](https://via.placeholder.com/16/{{.Color}}/FFFFFF/?text=%20) `{{.Label}}` {{end}}

edge: {{range .EdgeTypes}}![](https://via.placeholder.com/16/{{.Color}}/FFFFFF/?text=%20) `{{.Label}}` {{end}}
```mermaid
{{.MermaidData}}
```
