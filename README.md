# elfquery

An ELF file analysis tool written in Golang.

This tool parses the symbolic content of an ELF file and allows the data
to be analysed via:

- SQL queries (`elfquery sql`)
- A web interface (`elfquery http`)

Additional commands are also available, as described in the `--help` menu.

## Installation

ToDo

## Usage

### SQL Queries (`sql`)

Once parsed, the ELF file can be queried in the REPL, or by sending a SQL query
with the `-q` parameter:

```bash
$ elfquery sql samples/lpc55s69_zephyr.elf -q \
  "SELECT printf('0x%X', Value) AS Address, Name, Binding, Size FROM symbols \
  WHERE Section LIKE 'bss' ORDER BY Size DESC LIMIT 10"

+------------+--------------------------+---------+------+
| ADDRESS    | NAME                     | BINDING | SIZE |
+------------+--------------------------+---------+------+
| 0x30000110 | z_main_thread            | global  | 128  |
| 0x30000090 | z_idle_threads           | global  | 128  |
| 0x300001C0 | gpio_mcux_lpc_port0_data | local   | 80   |
| 0x30000210 | gpio_mcux_lpc_port1_data | local   | 80   |
| 0x30000298 | _kernel                  | global  | 48   |
| 0x30000270 | s_pintCallback           | local   | 32   |
| 0x300001AC | dyn_reg_info             | local   | 20   |
| 0x30000290 | s_secpintCallback        | local   | 8    |
| 0x30000190 | curr_tick                | local   | 8    |
| 0x30000260 | mcux_flexcomm_0_data     | local   | 8    |
+------------+--------------------------+---------+------+
```

#### Output Formatting

The following output options (`-o`) are supported:

- `text`: ASCII table (**default**)
- `pretty`: Unicode table
- `color`: Color unicode table
- `csv`: Comma-separated value table
- `md`: Markdown table
- `html`: HTML table
- `json`: JSON data

#### Table Definitions

Two tables are available in the SQLite database:

- `symbols`
```
  ID            Integer  Internal autoincrementing counter for symbols
  Value         Integer  Value associated with the symbol
  Size          Integer  Size in bytes
  Type          Text     Symbol type (data, code, etc.)
  Binding       Text     Symbol binding type (local, global, weak, etc.)
  Visibility    Text     Symbol visiblity (default, hidden, etc.)
  SectionIndex  Integer  Section index
  Name          Text     Symbol name
  Section       Text     Section name
```

 - `sections`

```
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
```

#### SQL Examples

To list all sections in the ELF file:

```SQL
SELECT Name, printf('0x%X', Address) AS Address, Size FROM sections
```

To list every symbol in order:

```SQL
SELECT * FROM symbols ORDER BY ID ASC
```

To select the name and size of each symbol in the 'bss' section:

```SQL
SELECT Name, Size FROM symbols WHERE Section = 'bss'
```

To do the same query but restrict it to the 10 largest symbols:

```SQL
SELECT Name, Size FROM symbols WHERE Section = 'bss' ORDER BY Size DESC LIMIT 10
```

To show 'Weak' symbols implemented in the ELF file:

```SQL
SELECT * FROM symbols WHERE Binding LIKE 'weak'
```

Any SQL query supported by SQLite3 can used!

### HTTP

You can analyse the contents of the ELF file in any web browser via the
`http` command.

```bash
$ elfquery http samples/lpc55s69_zephyr.elf
Starting HTTP server on port http://localhost:1443
```

TODO: Animated GIF

### Command Line

```bash
$ elfquery info samples/lpc55s69_zephyr.elf
Machine: ARM
ELF Class: 32 bits
ELF Type: ET_EXEC
ELF Data: ELFDATA2LSB
OS ABI: ELFOSABI_NONE
OS ABI Version: 0x0
Entry Point: 0x10000C15
```
