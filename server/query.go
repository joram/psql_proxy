package server

import (
	"database/sql"
	"strings"
)

type QueryType int

const (
	SELECT QueryType = iota
	INSERT
	UPDATE
	DELETE
	DROP
	CREATE
	DESCRIBE
)

func getQueryType(query string) QueryType {
	if isSelect(query) {
		return SELECT
	}
	if isInsert(query) {
		return INSERT
	}
	if isUpdate(query) {
		return UPDATE
	}
	if isDelete(query) {
		return DELETE
	}
	if isCreate(query) {
		return CREATE
	}
	if isDescribe(query) {
		return DESCRIBE
	}
	if isDrop(query) {
		return DROP
	}
	return -1
}

func isSelect(query string) bool {
	return isQueryType(query, "select")
}

func isInsert(query string) bool {
	return isQueryType(query, "insert")
}

func isUpdate(query string) bool {
	return isQueryType(query, "update")
}

func isDelete(query string) bool {
	return isQueryType(query, "delete")
}

func isCreate(query string) bool {
	return isQueryType(query, "create")
}

func isDescribe(query string) bool {
	return isQueryType(query, "\\d")
}

func isDrop(query string) bool {
	return isQueryType(query, "drop")
}

func isQueryType(query string, typeStr string) bool {
	return strings.HasPrefix(strings.ToLower(query), typeStr)
}

func RowsToMaps(rows *sql.Rows) ([]map[string]string, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	columnCount := len(columns)

	cursor := make([]interface{}, columnCount)
	for i := 0; i < columnCount; i++ {
		var columnValue string
		cursor[i] = &columnValue
	}

	var resultMaps []map[string]string
	for rows.Next() {
		err := rows.Scan(cursor...)
		if err != nil {
			return resultMaps, err
		}
		rowMap := make(map[string]string, columnCount)
		for i, columnPtr := range cursor {
			key := columns[i]
			var columnStr string
			if columnStrPtr := columnPtr.(*string); columnStrPtr != nil {
				columnStr = *columnStrPtr
			}
			rowMap[key] = columnStr
		}
		resultMaps = append(resultMaps, rowMap)
	}
	if err := rows.Err(); err != nil {
		return resultMaps, err
	}
	return resultMaps, nil
}
