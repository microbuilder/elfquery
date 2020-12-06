package elf2sql

import (
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"reflect"

	"github.com/yalue/elf_reader"

	// Use sqlite3 for the SQL database
	_ "github.com/mattn/go-sqlite3"
)

// SymBinding represents symbolic table binding values
type SymBinding uint8

// Symbolic table binding value definitions
const (
	StbLocal  SymBinding = 0  // Local symbol
	StbGlobal            = 1  // Global symbol
	StbWeak              = 2  // Weak global symbol
	StbLoOs              = 10 // Reserved for OS-specific semantics
	StbHiOs              = 12 // Reserved for OS-specific semantics
	StbLoProc            = 13 // Reserved for processor-specific semantics
	StbHiProc            = 15 // Reserved for processor-specific semantics
)

// String map for SymType values
var symBindingStrings = map[SymBinding]string{
	StbLocal:  "local",
	StbGlobal: "global",
	StbWeak:   "weak",
	StbLoOs:   "loos",
	StbHiOs:   "hios",
	StbLoProc: "loproc",
	StbHiProc: "hiproc",
}

// SymType represents ELF symbol types
type SymType uint8

// ELF
const (
	SttNoType  SymType = 0  // Not specified
	SttObject          = 1  // Data object (variable, array, etc.)
	SttFunc            = 2  // Function or executable code
	SttSection         = 3  // Section
	SttFile            = 4  // Name of source file
	SttCommon          = 5  // Uninitalised common block
	SttTLS             = 6  // Thread-local storage entity
	SttLoOs            = 10 // Reserved for OS-specific semantics
	SttHiOs            = 12 // Reserved for OS-specific semantics
	SttLoProc          = 13 // Reserved for processor-specific semantics
	SttHiProc          = 15 // Reserved for processor-specific semantics
)

// String map for SymType values
var symTypeStrings = map[SymType]string{
	SttNoType:  "none",
	SttObject:  "data",
	SttFunc:    "code",
	SttSection: "section",
	SttFile:    "filename",
	SttCommon:  "common",
	SttTLS:     "tls",
	SttLoOs:    "loos",
	SttHiOs:    "hios",
	SttLoProc:  "loproc",
	SttHiProc:  "hiproc",
}

// SymVisibility represents ELF symbol visibility values
type SymVisibility uint8

// ELF
const (
	SymVisDefault   SymVisibility = 0 // Default
	SymVisInternal                = 1 // Reserved value
	SymVisHidden                  = 2 // Not visible to other components
	SymVisProtected               = 3 // Visible but can't be preeùpted
	SymVisExported                = 4 // Global
	SymVisSingleton               = 5 // Global, singleton
	SymVisEliminate               = 6 // Extends hidden
)

// String map for SymVisibility values
var symVisStrings = map[SymVisibility]string{
	SymVisDefault:   "default",
	SymVisInternal:  "internal",
	SymVisHidden:    "hidden",
	SymVisProtected: "protected",
	SymVisExported:  "exported",
	SymVisSingleton: "singleton",
	SymVisEliminate: "eliminate",
}

// Section encapsulates a section entry in the DB
type Section struct {
	id          int
	name        string
	stype       string
	flags       string
	address     uint64
	offset      uint64
	size        uint64
	linkedindex uint32
	info        uint32
	alignment   uint64
	entrysize   uint64
}

// Symbol encapsulates a symbol entry in the DB
type Symbol struct {
	id           int
	value        uint64
	size         uint64
	symboltype   SymType
	binding      SymBinding
	visibility   SymVisibility
	sectionindex uint16
	name         string
	section      string
}

const createSectionTable string = `CREATE TABLE sections (
	ID          integer primary key,
	Name        text,
	Type        text,
	Flags       text,
	Address     integer,
	Offset      integer,
	Size        integer,
	LinkedIndex integer,
	Info        integer,
	Alignment   integer,
	EntrySize   integer
	)`

const createSymbolTable string = `CREATE TABLE symbols (
	ID           integer primary key autoincrement,
	Value        integer,
	Size         integer,
	Type         text,
	Binding      text,
	Visibility   text,
	SectionIndex integer,
	Name         text,
	Section      text
	)`

// RunQuery loads the specified ELF file into a memory-based SQLite database,
// and runs the specific SQL query against the database.
// The database contains two tables: 'sections' and 'symbols'.
func RunQuery(filename string, query string) {
	f, e := ioutil.ReadFile(filename)
	_elf, e := elf_reader.ParseELFFile(f)
	check(e)

	if query == "" {
		fmt.Printf("No query string provided")
		return
	}

	// Open a new SQLite database in memory
	db, e := sql.Open("sqlite3", ":memory:")
	check(e)
	defer db.Close()

	// Create sections table
	_, e = db.Exec(createSectionTable)
	check(e)

	// Create symbols table
	_, e = db.Exec(createSymbolTable)
	check(e)

	// Iterate over sections to populate the database
	count := _elf.GetSectionCount()
	for i := uint16(0); i < count; i++ {
		// Get section name
		_name, e := _elf.GetSectionName(i)
		if e != nil {
			_name = "<NULL>"
		}

		// Get section header
		header, e := _elf.GetSectionHeader(i)
		if e != nil {
			fmt.Printf("Error getting section %d header: %s", i, e)
		}
		_sec := Section{
			id:          int(i),
			name:        _name,
			stype:       fmt.Sprint(header.GetType()),
			flags:       fmt.Sprint(header.GetFlags()),
			address:     header.GetVirtualAddress(),
			offset:      header.GetFileOffset(),
			size:        header.GetSize(),
			linkedindex: header.GetLinkedIndex(),
			info:        header.GetInfo(),
			alignment:   header.GetAlignment(),
			entrysize:   header.GetEntrySize(),
		}

		// Insert the section into the DB
		tx, e := db.Begin()
		check(e)
		stmt, e := tx.Prepare(`INSERT INTO sections VALUES (?,?,?,?,?,?,?,?,?,?,?)`)
		check(e)
		defer stmt.Close()
		_, e = stmt.Exec(_sec.id, _sec.name, _sec.stype, _sec.flags,
			_sec.address, _sec.offset, _sec.size, _sec.linkedindex, _sec.info,
			_sec.alignment, _sec.entrysize)
		check(e)
		tx.Commit()

		// Get Symbols
		symbols, names, e := _elf.GetSymbols(i)
		if e == nil {
			for j := range symbols {
				// Assign symbol values
				_sym := Symbol{
					id:           int(j),
					value:        symbols[j].GetValue(),
					size:         symbols[j].GetSize(),
					symboltype:   SymType(symbols[j].GetInfo().SymbolType()),
					binding:      SymBinding(symbols[j].GetInfo().Binding()),
					visibility:   SymVisibility(symbols[j].GetOther()),
					sectionindex: symbols[j].GetSectionIndex(),
					name:         names[j],
				}

				// Lookup the matching section name
				if _sym.sectionindex >= 0xFF00 {
					_sym.section = "<Unknown>"
				} else {
					_sname, err := _elf.GetSectionName(_sym.sectionindex)
					if err != nil {
						_sname = "<NULL>"
					}
					_sym.section = _sname
				}

				// Insert symbol into table
				tx, e := db.Begin()
				check(e)
				stmt, e := tx.Prepare(`INSERT INTO symbols VALUES (NULL,?,?,?,?,?,?,?,?)`)
				check(e)
				defer stmt.Close()
				_, e = stmt.Exec(_sym.value, _sym.size,
					symTypeStrings[_sym.symboltype],
					symBindingStrings[_sym.binding],
					symVisStrings[_sym.visibility],
					_sym.sectionindex, _sym.name, _sym.section)
				check(e)
				tx.Commit()
			}
		}
	}

	// Execute the provided query
	rows, e := db.Query(query)
	check(e)
	defer rows.Close()

	// Display the column names
	cols, _ := rows.Columns()
	for _, col := range cols {
		fmt.Printf("%s, ", col)
	}
	fmt.Printf("\n\n")

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
			return
		}

		// User reflection to determine each row's value type
		for i := range cols {
			val := columnPointers[i].(*interface{})
			if *val != nil {
				switch reflect.Indirect(reflect.ValueOf(val)).Elem().Kind() {
				case reflect.String:
					fmt.Printf("%s, ", *val)
				case reflect.Int64:
					fmt.Printf("%d, ", *val)
				default:
					fmt.Printf("%s, ", *val)
				}
			}
		}
		fmt.Printf("\n")
	}

	// // Query section data
	// rows, e := db.Query("SELECT * FROM sections")
	// check(e)
	// defer rows.Close()
	// for rows.Next() {
	// 	var _sec Section
	// 	e := rows.Scan(&_sec.id, &_sec.name, &_sec.stype, &_sec.flags,
	// 		&_sec.address, &_sec.offset, &_sec.size, &_sec.linkedindex,
	// 		&_sec.info, &_sec.alignment, &_sec.entrysize)
	// 	check(e)
	// 	fmt.Printf("[%3d]: %s\n", _sec.id, _sec.name)
	// 	// fmt.Printf("[%3d]: %s\n", _id, _name, _type, _flags, _address, _offset, _size, _linkedindex, _info, _alignment, _entrysize)
	// }

	// // Query symbol data
	// rows, e := db.Query("SELECT * FROM symbols WHERE section = 'text' AND size > 0 ORDER BY Size DESC")
	// check(e)
	// defer rows.Close()
	// fmt.Printf("Symbols in 'text' > 0 bytes, sorted by size desc:\n")
	// for rows.Next() {
	// 	var _sym Symbol
	// 	e := rows.Scan(&_sym.id, &_sym.value, &_sym.size, &_sym.symboltype,
	// 		&_sym.binding, &_sym.visibility, &_sym.sectionindex,
	// 		&_sym.name, &_sym.section)
	// 	check(e)
	// 	fmt.Printf("%s, %d bytes, %s\n", _sym.name, _sym.size, _sym.section)
	// }
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func ioReader(file string) io.ReaderAt {
	r, err := os.Open(file)
	check(err)
	return r
}
