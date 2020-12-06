# goelf

An ELF file analysis utility written in Golang that allows you to analyse
ELF files via:

- a web interface (`goelf http`)
- SQL commands (`goelf sql`)
- the command line (`goelf info`)

## Installation

ToDo

## Usage

### HTTP

You can analyse the contents of the ELF file in any web browser via the
`http` command.

```bash
$ goelf http samples/lpc55s69_zephyr.elf
```

TODO: Animated GIF

#### Key Generation

The HTTP server requires a private key for TLS, which can be generated via:

```bash
$ openssl ecparam -name secp256r1 -genkey -out SERVER.key
```

You can then generate a self-signed X.509 certificate via:

```bash
$ openssl req -new -x509 -sha256 -days 3650 -key SERVER.key -out SERVER.crt \
        -subj "/O=Linaro, LTD/CN=localhost"
```

This certificate should be available on any device(s) connecting to the HTTP
server to verify that we are communicating with the intended server.

### SQL

The symbolic contents of the ELF file will be parsed into a memory-based
SQLite3 database, which can be queried in the REPL, or by sending a SQL query
via the `-q` parameter.

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

#### Sending Queries 

One-off queries can be executed as follows:

```bash
$ goelf sql samples/lpc55s69_zephyr.elf \
  -q "SELECT printf('0x%X', Value) AS Address, Size, Type, Binding, \
  Visibility, Section, Name FROM symbols ORDER BY Size DESC LIMIT 10"
```

Which, depending on your ELF file, may result in something resembling:

```
Address, Size, Type, Binding, Visibility, Section, Name

0x30000820, 2048, data, global, default, noinit, z_interrupt_stacks
0x300002E0, 1024, data, global, default, noinit, z_main_stack
0x10000161, 750, code, global, hidden, text, __udivmoddi4
0x10000511, 604, code, global, default, text, z_vprintk
0x10003128, 480, data, global, default, sw_isr_table, _sw_isr_table
0x300006E0, 320, data, local, default, noinit, z_idle_stacks
0x10000F81, 316, code, local, default, text, mpu_configure_regions_and_partition.constprop.0
0x10002D99, 300, code, global, default, text, USART_Init
0x10000D31, 280, code, global, default, text, z_arm_fault
0x100021B5, 272, code, global, default, text, z_add_timeout
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
````

To do the same query but restrict it to the 10 largest symbols:

```SQL
SELECT Name, Size  FROM symbols WHERE Section = 'bss' ORDER BY Size DESC LIMIT 10
````

To show 'Weak' symbols implemented in the ELF file:

```SQL
SELECT * FROM symbols WHERE Binding LIKE 'weak'
```

Any SQL query supported by SQLite3 can used!

### Command Line

```bash
$ goelf info samples/lpc55s69_zephyr.elf
Machine: ARM
ELF Class: 32 bits
ELF Type: ET_EXEC
ELF Data: ELFDATA2LSB
OS ABI: ELFOSABI_NONE
OS ABI Version: 0x0
Entry Point: 0x10000C15
```
