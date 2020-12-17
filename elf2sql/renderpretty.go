package elf2sql

import (
	"database/sql"
	"fmt"
	"reflect"

	"github.com/jedib0t/go-pretty/table"
)

// PrettyFormat represents the various pretty output options
type PrettyFormat uint8

// Pretty format values
const (
	PrettyASCII    PrettyFormat = 0 // ASCII table
	PrettyUnicode               = 1 // Unicode table
	PrettyColor                 = 2 // Unicode color table
	PrettyMarkdown              = 3 // Markdown table
	PrettyHTML                  = 4 // HTML table
	PrettyCSV                   = 5 // CSV table
)

// Renders rows in pretty text
func renderPretty(rows *sql.Rows, format PrettyFormat) string {
	t := table.NewWriter()
	// t.SetOutputMirror(os.Stdout)

	var tr table.Row

	// Display the column names
	cols, _ := rows.Columns()
	for _, col := range cols {
		tr = append(tr, fmt.Sprintf("%s", col))
	}
	t.AppendHeader(tr)

	// Iterate over each row
	for rows.Next() {
		tr = nil
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
					tr = append(tr, fmt.Sprintf("%s", *val))
				case reflect.Int64:
					tr = append(tr, fmt.Sprintf("%d", *val))
				default:
					tr = append(tr, fmt.Sprintf("%s", *val))
				}
			}
		}
		t.AppendRow(tr)
	}

	switch format {
	case PrettyASCII:
		return fmt.Sprintf("%s\n", t.Render())
	case PrettyUnicode:
		t.SetStyle(table.StyleLight)
		return fmt.Sprintf("%s\n", t.Render())
	case PrettyColor:
		t.SetStyle(table.StyleColoredBright)
		return fmt.Sprintf("%s\n", t.Render())
	case PrettyMarkdown:
		return fmt.Sprintf("%s\n", t.RenderMarkdown())
	case PrettyHTML:
		t.SetHTMLCSSClass("table")
		return fmt.Sprintf("%s\n", t.RenderHTML())
	case PrettyCSV:
		return fmt.Sprintf("%s\n", t.RenderCSV())
	default:
		return fmt.Sprintf("%s\n", t.Render())
	}
}
