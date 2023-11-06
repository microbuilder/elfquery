package cmd

import (
	"debug/elf"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info filename",
	Short: "Basic file details",
	Long: `Lists key information about the specified ELF file, such as the
target machine, ELF file type, sections, etc.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		full, _ := cmd.Flags().GetBool("full")

		f := ioReader(args[0])
		_elf, err := elf.NewFile(f)
		check(err)

		// Read and decode ELF identifier
		var ident [16]uint8
		f.ReadAt(ident[0:], 0)
		check(err)

		if ident[0] != '\x7f' || ident[1] != 'E' || ident[2] != 'L' || ident[3] != 'F' {
			fmt.Printf("Bad magic number at %d\n", ident[0:4])
			os.Exit(1)
		}

		var arch string
		switch _elf.Class.String() {
		case "ELFCLASS64":
			arch = "64 bits"
		case "ELFCLASS32":
			arch = "32 bits"
		}

		var mach string
		switch _elf.Machine.String() {
		case "EM_AARCH64":
			mach = "ARM64"
		case "EM_ARM":
			mach = "ARM"
		case "EM_386":
			mach = "x86"
		case "EM_X86_64":
			mach = "x86_64"
		}

		fmt.Printf("Machine: %s\n", mach)
		fmt.Printf("ELF Class: %s\n", arch)
		fmt.Printf("ELF Type: %s\n", _elf.Type)
		fmt.Printf("ELF Data: %s\n", _elf.Data)
		fmt.Printf("OS ABI: %s\n", _elf.OSABI)
		fmt.Printf("OS ABI Version: 0x%X\n", _elf.ABIVersion)
		fmt.Printf("Entry Point: 0x%08X\n", _elf.Entry)

		// Calculate size data across each section
		sztext, szdata, szbss := 0, 0, 0
		for _, s := range _elf.Sections {
			_text, _data, _bss := sectionSize(s.SectionHeader)
			sztext += _text
			szdata += _data
			szbss += _bss
		}
		fmt.Printf("Text size: %d\n", sztext)
		fmt.Printf("Data size: %d\n", szdata)
		fmt.Printf("BSS size: %d\n", szbss)
		fmt.Printf("Total size: %d\n", sztext+szdata+szbss)

		if full {
			// Display individual sections
			fmt.Printf("Sections:\n")
			for _, s := range _elf.Sections {
				fmt.Printf("  0x%08X\t%d\t%s\n", s.Addr, s.Size, s.Name)
			}
		}
	},
}

// sectionSize determine the text, data and bss size for the supplied ELF
// section header using the same algorithm as GNU binutils 'size' tool.
func sectionSize(sec elf.SectionHeader) (int, int, int) {
	text, data, bss := 0, 0, 0

	// Only count allocated memory
	if strings.Contains(sec.Flags.String(), "SHF_ALLOC") {
		// Text consists of executable instructions, or not writable
		if strings.Contains(sec.Flags.String(), "SHF_EXECINSTR") ||
			!strings.Contains(sec.Flags.String(), "SHF_WRITE") {
			text = int(sec.Size)
		} else {
			// No data means bss
			if strings.Contains(sec.Type.String(), "SHT_NOBITS") {
				bss = int(sec.Size)
			} else {
				// Otherwise, count as data
				data = int(sec.Size)
			}
		}
	}

	return text, data, bss
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

func init() {
	rootCmd.AddCommand(infoCmd)

	infoCmd.Flags().BoolP("full", "f", false, "Display full result set")
}
