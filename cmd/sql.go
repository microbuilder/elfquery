package cmd

import (
	"fmt"

	"github.com/microbuilder/elfquery/elf2sql"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// sqlCmd represents the sql command
var sqlCmd = &cobra.Command{
	Use:   "sql filename",
	Short: "Run SQL queries against the ELF file",
	Long: `Reads all symbolic information from the ELF file and adds it to an
in-memory SQLite database, which can be queried in the REPL or via a SQL
query string (-q).

Two tables are available in the SQLite database:

  symbols

  ID            Integer  Internal autoincrementing counter for symbols
  Value         Integer  Value associated with the symbol
  Size          Integer  Size in bytes
  Type          Text     Symbol type (data, code, etc.)
  Binding       Text     Symbol binding type (local, global, weak, etc.)
  Visibility    Text     Symbol visiblity (default, hidden, etc.)
  SectionIndex  Integer  Section index
  Name          Text     Symbol name
  Section       Text     Section name

  sections

  ID            Integer   Numeric ID to distinction sections
  Name          Text      Section name
  Type          Text      Section type
  Flags         Text      Section attribute flags
  Address       Integer   Address of the first byte of this section
  Offset        Integer   Offset from the start of file
  Size          Integer   Section size in bytes
  LinkedIndex   Integer   Section header table index
  Info          Integer   Extra information (usage varies)
  Alignment     Integer   Address alignment constraints
  EntrySize     Integer   Size in bytes of each fixed-size entry

To list all sections in the ELF file ('sections' alias):

  SELECT Name, printf('0x%X', Address) AS Address, Size FROM sections

To list every symbol in order ('symbols' alias):

  SELECT * FROM symbols ORDER BY ID ASC

To select the name and size of each symbol in the 'bss' section ('bss' alias):

  SELECT Name, Size FROM symbols WHERE Section = 'bss'

To do the same query but restrict it to the 10 largest symbols ('bss10' alias):

  SELECT Name, Size FROM symbols WHERE Section = 'bss' ORDER BY Size DESC LIMIT 10

To show 'Weak' symbols implemented in the ELF file ('weak' alias):

  SELECT * FROM symbols WHERE Binding LIKE 'weak'
`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Check display format
		output, _ := cmd.Flags().GetString("output")
		outputDict := map[string]elf2sql.DisplayFormat{
			"text":   elf2sql.DFText,
			"pretty": elf2sql.DFPretty,
			"color":  elf2sql.DFPrettyCol,
			"csv":    elf2sql.DFCSV,
			"md":     elf2sql.DFMarkdown,
			"html":   elf2sql.DFHtml,
			"json":   elf2sql.DFJson,
		}
		df, ok := outputDict[output]
		if !ok {
			fmt.Printf("invalid output flag: %s\n", output)
			return
		}

		// Populate the database with the ELF data
		e := elf2sql.InitDB(args[0])
		defer elf2sql.CloseDB()
		if e != nil {
			fmt.Printf("unable to initialise the SQLite3 database in memory\n")
			return
		}

		// Check for SQL aliases
		alias, _ := cmd.Flags().GetString("alias")
		if alias != "" {
			aliases := viper.GetStringMapString("sqlaliases")

			query, exists := aliases[alias]
			if !exists {
				fmt.Printf("Invalid SQL alias: %s\n", alias)
				if viper.ConfigFileUsed() == "" {
					fmt.Printf("No elfquery.toml found to parse sqlaliases!\n")
					return
				}
				fmt.Printf("%s defines the following 'sqlaliases':\n", viper.ConfigFileUsed())
				for a := range aliases {
					fmt.Printf(" - %s\n", a)
				}
				return
			}

			// Request and display the alias query results
			s, e := elf2sql.RunQuery(query, df)
			if e != nil {
				fmt.Printf("invalid query: %s\n", query)
				return
			}
			fmt.Printf(s)
			return
		}

		// Parse SQL query if no alias is provided
		query, _ := cmd.Flags().GetString("query")
		if query == "" {
			fmt.Printf("TODO: REPL mode\n")
			return
		}

		// Request and display the alias query results
		s, e := elf2sql.RunQuery(query, df)
		if e != nil {
			fmt.Printf("invalid query: %s\n", query)
			return
		}
		fmt.Printf(s)
	},
}

func init() {
	rootCmd.AddCommand(sqlCmd)

	sqlCmd.Flags().StringP("query", "q", "", "SQL query to execute")
	sqlCmd.Flags().StringP("alias", "a", "", "SQL alias to execute (see .elfquery.toml)")
	sqlCmd.Flags().StringP("output", "o", "text", "output format (text, pretty, color, csv, md, html, json)")
}
