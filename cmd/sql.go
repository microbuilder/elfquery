package cmd

import (
	"fmt"

	"github.com/microbuilder/goelf/elf2sql"
	"github.com/spf13/cobra"
)

// sqlCmd represents the sql command
var sqlCmd = &cobra.Command{
	Use:   "sql [filename]",
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

To list all sections in the ELF file:

  SELECT Name, printf('0x%X', Address) AS Address, Size FROM sections

To list every symbol in order:

  SELECT * FROM symbols ORDER BY ID ASC

To select the name and size of each symbol in the 'bss' section:

  SELECT Name, Size FROM symbols WHERE Section = 'bss'

To do the same query but restrict it to the 10 largest symbols:

  SELECT Name, Size  FROM symbols WHERE Section = 'bss' ORDER BY Size DESC LIMIT 10

To show 'Weak' symbols implemented in the ELF file:

  SELECT * FROM symbols WHERE Binding LIKE 'weak'
`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		query, _ := cmd.Flags().GetString("query")
		if query == "" {
			fmt.Printf("TODO: REPL mode\n")
		} else {
			elf2sql.RunQuery(args[0], query)
		}
	},
}

func init() {
	rootCmd.AddCommand(sqlCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// sqlCmd.PersistentFlags().String("foo", "", "A help for foo")

	sqlCmd.Flags().StringP("query", "q", "", "SQL query to execute")
	sqlCmd.Flags().StringP("output", "o", "", "output format")
}
