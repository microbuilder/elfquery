package elf2sql

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

// Renders rows in plain text
func renderText(rows *sql.Rows) string {
	var sb strings.Builder

	// Display the column names
	cols, _ := rows.Columns()
	for _, col := range cols {
		sb.WriteString(fmt.Sprintf("%s, ", col))
	}
	sb.WriteString(fmt.Sprintf("\n\n"))

	// Iterate over each row
	for rows.Next() {
		// Use reflection to parse each column for its data type
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		// Scan the row, dumping the values into columnPointers
		if err := rows.Scan(columnPointers...); err != nil {
			return ""
		}

		// User reflection to determine each row's value type
		for i := range cols {
			val := columnPointers[i].(*interface{})
			if *val != nil {
				switch reflect.Indirect(reflect.ValueOf(val)).Elem().Kind() {
				case reflect.String:
					sb.WriteString(fmt.Sprintf("%s, ", *val))
				case reflect.Int64:
					sb.WriteString(fmt.Sprintf("%d, ", *val))
				default:
					sb.WriteString(fmt.Sprintf("%s, ", *val))
				}
			}
		}
		sb.WriteString(fmt.Sprintf("\n"))
	}

	return sb.String()
}
