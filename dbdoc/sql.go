package dbdoc

import (
	"fmt"
	"go/token"
	"log"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	"github.com/mazrean/isucrud/internal/pkg/list"
)

var (
	tableRe        = regexp.MustCompile("^\\s*[\\[\"'`]?(?P<Table>\\w+)[\\]\"'`]?\\s*")
	insertRe       = regexp.MustCompile("^insert\\s+(ignore\\s+)?(into\\s+)?[\\[\"'`]?(?P<Table>\\w+)[\\]\"'`]?\\s*")
	deleteRe       = regexp.MustCompile("^delete\\s+from\\s+[\\[\"'`]?(?P<Table>\\w+)[\\]\"'`]?\\s*")
	selectKeywords = []string{" where ", " group by ", " having ", " window ", " order by ", "limit ", " for "}
)

func AnalyzeSQL(ctx *Context, sql stringLiteral) []Query {
	sqlValue := strings.ToLower(sql.value)

	strQueries := extractSubQueries(ctx, sqlValue)

	var queries []Query
	for _, sqlValue := range strQueries {
		newQueries := analyzeSQLWithoutSubQuery(ctx, sqlValue, sql.pos)
		for _, query := range newQueries {
			fmt.Println(query)
		}
		queries = append(queries, newQueries...)
	}

	return queries
}

var (
	subQueryPrefixRe = regexp.MustCompile(`^\s*\(\s*select\s+`)
)

func extractSubQueries(_ *Context, sql string) []string {
	var subQueries []string

	type subQuery struct {
		query        string
		bracketCount uint
	}

	rootQuery := ""
	subQueryStack := list.NewStack[*subQuery]()
	for i := 0; i < len(sql); i++ {
		r := sql[i]
		switch r {
		case '(':
			if subQuery, ok := subQueryStack.Peek(); ok {
				subQuery.bracketCount++
				subQuery.query += string(r)
			} else {
				rootQuery += string(r)
			}

			match := subQueryPrefixRe.FindString(sql[i:])
			if len(match) != 0 {
				subQueryStack.Push(&subQuery{
					query:        match,
					bracketCount: 0,
				})
				i += len(match)
				continue
			}
		case ')':
			if subQuery, ok := subQueryStack.Peek(); ok && subQuery.bracketCount == 0 {
				subQueries = append(subQueries, subQuery.query)
				subQueryStack.Pop()
			}

			if subQuery, ok := subQueryStack.Peek(); ok {
				subQuery.bracketCount--
				subQuery.query += string(r)
			} else {
				rootQuery += string(r)
			}
		default:
			if subQuery, ok := subQueryStack.Peek(); ok {
				subQuery.query += string(r)
			} else {
				rootQuery += string(r)
			}
		}
	}

	for subQuery, ok := subQueryStack.Pop(); ok; subQuery, ok = subQueryStack.Pop() {
		subQueries = append(subQueries, subQuery.query)
	}

	if rootQuery != "" {
		subQueries = append(subQueries, rootQuery)
	}

	return subQueries
}

func analyzeSQLWithoutSubQuery(ctx *Context, sqlValue string, pos token.Pos) []Query {
	sqlValue = strings.TrimLeftFunc(sqlValue, unicode.IsSpace)
	sqlValue = replaceMultipleWhitespace(sqlValue)

	var queries []Query
	switch {
	case strings.HasPrefix(sqlValue, "select"):
		_, after, found := strings.Cut(sqlValue, " from ")
		if !found {
			tableNames := tableForm(ctx, sqlValue, pos)

			for _, tableName := range tableNames {
				queries = append(queries, Query{
					QueryType: QueryTypeSelect,
					Table:     tableName,
					Pos:       pos,
					Raw:       sqlValue,
				})
			}
			break
		}

		tmpTableNames := strings.Split(after, ",")
		var tableNames []string
	TABLE_LOOP:
		for _, tableName := range tmpTableNames {
			tableNames = append(tableNames, strings.Split(tableName, " join ")...)

			for _, keyword := range selectKeywords {
				if strings.Contains(tableName, keyword) {
					break TABLE_LOOP
				}
			}
		}

		for _, tableName := range tableNames {
			matches := tableRe.FindStringSubmatch(tableName)
			if len(matches) == 0 {
				continue
			}

			for i, name := range tableRe.SubexpNames() {
				if name == "Table" {
					queries = append(queries, Query{
						QueryType: QueryTypeSelect,
						Table:     matches[i],
						Pos:       pos,
						Raw:       sqlValue,
					})
				}
			}
		}
	case strings.HasPrefix(sqlValue, "insert"):
		matches := insertRe.FindStringSubmatch(sqlValue)

		for i, name := range insertRe.SubexpNames() {
			if name == "Table" {
				queries = append(queries, Query{
					QueryType: QueryTypeInsert,
					Table:     matches[i],
					Pos:       pos,
					Raw:       sqlValue,
				})
			}
		}
	case strings.HasPrefix(sqlValue, "update"):
		afterUpdate := strings.TrimPrefix(sqlValue, "update ")
		before, _, found := strings.Cut(afterUpdate, " set ")
		if !found {
			before = afterUpdate
		}

		tmpTableNames := strings.Split(before, ",")
		var tableNames []string
		for _, tableName := range tmpTableNames {
			tableNames = append(tableNames, strings.Split(tableName, " join ")...)
		}

		for _, tableName := range tableNames {
			matches := tableRe.FindStringSubmatch(tableName)
			if len(matches) == 0 {
				continue
			}

			for i, name := range tableRe.SubexpNames() {
				if name == "Table" {
					queries = append(queries, Query{
						QueryType: QueryTypeUpdate,
						Table:     matches[i],
						Pos:       pos,
						Raw:       sqlValue,
					})
				}
			}
		}
	case strings.HasPrefix(sqlValue, "delete"):
		matches := deleteRe.FindStringSubmatch(sqlValue)

		for i, name := range deleteRe.SubexpNames() {
			if name == "Table" {
				queries = append(queries, Query{
					QueryType: QueryTypeDelete,
					Table:     matches[i],
					Pos:       pos,
					Raw:       sqlValue,
				})
			}
		}
	}

	return queries
}

func replaceMultipleWhitespace(s string) string {
	fields := strings.Fields(s)
	return strings.Join(fields, " ")
}

func tableForm(ctx *Context, sqlValue string, pos token.Pos) []string {
	position := ctx.FileSet.Position(pos)
	filename, err := filepath.Rel(ctx.WorkDir, position.Filename)
	if err != nil {
		log.Printf("failed to get relative path: %v", err)
		return nil
	}

	fmt.Printf("query:%s\n", sqlValue)
	fmt.Printf("position: %s:%d:%d\n", filename, position.Line, position.Column)
	fmt.Print("table name?: ")
	var input string
	_, err = fmt.Scanln(&input)
	if err != nil {
		return nil
	}

	if input == "" {
		return nil
	}

	tableNames := strings.Split(input, ",")

	return tableNames
}
