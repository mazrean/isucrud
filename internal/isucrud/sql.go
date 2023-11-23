package isucrud

import (
	"fmt"
	"go/token"
	"log"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/mazrean/isucrud/internal/pkg/list"
)

var (
	tableRe        = regexp.MustCompile("^\\s*[\\[\"'`]?(?P<Table>\\w+)[\\]\"'`]?\\s*")
	insertRe       = regexp.MustCompile("^insert\\s+(ignore\\s+)?(into\\s+)?[\\[\"'`]?(?P<Table>\\w+)[\\]\"'`]?\\s*")
	deleteRe       = regexp.MustCompile("^delete\\s+from\\s+[\\[\"'`]?(?P<Table>\\w+)[\\]\"'`]?\\s*")
	selectKeywords = []string{" where ", " group by ", " having ", " window ", " order by ", "limit ", " for "}
)

func AnalyzeSQL(ctx *context, sql stringLiteral) []query {
	sqlValue := strings.ToLower(sql.value)

	strQueries := extractSubQueries(ctx, sqlValue)

	var queries []query
	for _, sqlValue := range strQueries {
		newQueries := analyzeSQLWithoutSubQuery(ctx, sqlValue, sql.pos)
		for _, query := range newQueries {
			fmt.Printf("%s(%s): %s\n", query.queryType, query.table, sqlValue)
		}
		queries = append(queries, newQueries...)
	}

	return queries
}

var (
	subQueryPrefixRe = regexp.MustCompile(`^\s*\(\s*select\s+`)
)

func extractSubQueries(ctx *context, sql string) []string {
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

func analyzeSQLWithoutSubQuery(ctx *context, sqlValue string, pos token.Pos) []query {
	var queries []query
	switch {
	case strings.HasPrefix(sqlValue, "select"):
		_, after, found := strings.Cut(sqlValue, " from ")
		if !found {
			tableNames := tableForm(ctx, sqlValue, pos)

			for _, tableName := range tableNames {
				queries = append(queries, query{
					queryType: queryTypeSelect,
					table:     tableName,
					pos:       pos,
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
					queries = append(queries, query{
						queryType: queryTypeSelect,
						table:     matches[i],
						pos:       pos,
					})
				}
			}
		}
	case strings.HasPrefix(sqlValue, "insert"):
		matches := insertRe.FindStringSubmatch(sqlValue)

		for i, name := range insertRe.SubexpNames() {
			if name == "Table" {
				queries = append(queries, query{
					queryType: queryTypeInsert,
					table:     matches[i],
					pos:       pos,
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
					queries = append(queries, query{
						queryType: queryTypeUpdate,
						table:     matches[i],
						pos:       pos,
					})
				}
			}
		}
	case strings.HasPrefix(sqlValue, "delete"):
		matches := deleteRe.FindStringSubmatch(sqlValue)

		for i, name := range deleteRe.SubexpNames() {
			if name == "Table" {
				queries = append(queries, query{
					queryType: queryTypeDelete,
					table:     matches[i],
					pos:       pos,
				})
			}
		}
	}

	return queries
}

func tableForm(ctx *context, sqlValue string, pos token.Pos) []string {
	position := ctx.fileSet.Position(pos)
	filename, err := filepath.Rel(ctx.workDir, position.Filename)
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
