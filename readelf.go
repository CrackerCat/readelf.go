package main

import (
	"os"
	"strings"
	"debug/elf"
	"log"
	"fmt"
)

var (
	elfPath                 string
	isShowHeader            bool
	isShowSection           bool
	isShowProgram           bool
	isShowSymbols           bool
	isShowDynamicSymbols    bool
	isShowAll               bool
	isShowImportedLibraries bool
	isShowImportedSymbols   bool
)

func init() {
	log.SetOutput(os.Stdout)
}

func main() {
	for i, arg := range os.Args {
		if i == 0 {
			continue
		}
		if strings.HasPrefix(arg, "-") { // flags
			if strings.Contains(arg, "a") {
				isShowAll = true
			}
			if strings.Contains(arg, "h") {
				isShowHeader = true
			}
			if strings.Contains(arg, "S") {
				isShowSection = true
			}
			if strings.Contains(arg, "l") {
				isShowProgram = true
			}
			if strings.Contains(arg, "s") {
				isShowSymbols = true
			}
			if strings.Contains(arg, "d") {
				isShowDynamicSymbols = true
			}
			if strings.Contains(arg, "I") {
				isShowImportedLibraries = true
			}
			if strings.Contains(arg, "i") {
				isShowImportedSymbols = true
			}
		} else { // elf path
			elfPath = arg
		}
	}
	if elfPath == "" {
		showUsage()
		os.Exit(1)
	}
	elfFile, err := elf.Open(elfPath)
	if err != nil {
		log.Fatal(err)
	}
	defer elfFile.Close()
	if isShowHeader || isShowAll {
		showHeader(elfFile)
	}
	if isShowSection || isShowAll {
		showSections(elfFile.Sections)
	}
	if isShowProgram || isShowAll {
		showProgram(elfFile.Progs)
	}
	if isShowSymbols || isShowAll {
		symbols, err := elfFile.Symbols()
		if err != nil {
			log.Fatal(err)
		}
		showSymbols(symbols, "Symbols")
	}
	if isShowDynamicSymbols || isShowAll {
		symbols, err := elfFile.DynamicSymbols()
		if err != nil {
			log.Fatal(err)
		}
		showSymbols(symbols, "DynamicSymbols")
	}
	if isShowImportedLibraries || isShowAll {
		libraries, err := elfFile.ImportedLibraries()
		if err != nil {
			log.Fatal(err)
		}
		showImportedLibraries(libraries)
	}
	if isShowImportedSymbols || isShowAll {
		symbols, err := elfFile.ImportedSymbols()
		if err != nil {
			log.Fatal(err)
		}
		showImportedSymbols(symbols)
	}
}

func showUsage() {
	fmt.Printf("Usage: readelf [-hSlsdaIi] elfPath\n")
	fmt.Printf("  -h: ELF File Header\n")
	fmt.Printf("  -S: ELF Sections\n")
	fmt.Printf("  -l: ELF Programs\n")
	fmt.Printf("  -s: ELF Symbols\n")
	fmt.Printf("  -d: ELF DynamicSymbols\n")
	fmt.Printf("  -I: ELF ImportedLibraries\n")
	fmt.Printf("  -i: ELF ImportedSymbols\n")
	fmt.Printf("  -a: All above\n")
}

func showHeader(file *elf.File) {
	fmt.Printf("ELF File Header:\n")
	fmt.Printf("  Class:      %s\n", file.Class)
	fmt.Printf("  Version:    %s\n", file.Version)
	fmt.Printf("  Data:       %s\n", file.Data)
	fmt.Printf("  OSABI:      %s\n", file.OSABI)
	fmt.Printf("  ABIVersion: %d\n", file.ABIVersion)
	fmt.Printf("  ByteOrder:  %s\n", file.ByteOrder)
	fmt.Printf("  Type:       %s\n", file.Type)
	fmt.Printf("  Machine:    %s\n", file.Machine)
	fmt.Printf("  Entry:      %d\n", file.Entry)
}

func showSections(sections []*elf.Section) {
	fmt.Printf("ELF Sections:\n")
	fmt.Printf("  [%2s] %-25s %-15s %-8s %-8s %-8s %2s %3s %3s %3s %4s\n", "Nr",
		"Name",
		"Type",
		"Addr",
		"Off",
		"Size",
		"ES",
		"Flg",
		"Lk",
		"Inf",
		"Al")
	for i, section := range sections {
		fmt.Printf("  [%2d] %-25s %-15s %08x %08x %08x %02x %3s %3d %3d %4x\n", i,
			section.Name,
			strings.Replace(section.Type.String(), "SHT_", "", -1),
			section.Addr,
			section.Offset,
			section.Size,
			section.Entsize,
			strings.Replace(section.Flags.String(), "SHF_", "", -1),
			section.Link,
			section.Info,
			section.Addralign)
	}
}

func showProgram(progs []*elf.Prog) {
	fmt.Printf("ELF Programs:\n")
	fmt.Printf("  [%2s] %-20s %-8s %-8s %-8s %-8s %-8s %3s %5s\n", "Nr",
		"Type",
		"Offset",
		"VirtAddr",
		"PhysAddr",
		"FileSiz",
		"MemSiz",
		"Flg",
		"Align")
	for i, prog := range progs {
		fmt.Printf("  [%2d] %-20s %08x %08x %08x %08x %08x %3s %5d\n", i,
			strings.Replace(prog.Type.String(), "PT_", "", -1),
			prog.Off,
			prog.Vaddr,
			prog.Paddr,
			prog.Filesz,
			prog.Memsz,
			strings.Replace(prog.Flags.String(), "PF_", "", -1),
			prog.Align)
	}
}

func showSymbols(symbols []elf.Symbol, title string) {
	fmt.Printf("ELF %s:\n", title)
	fmt.Printf("  [%2s] %-30s %-10s %-8s %-8s %3s %3s\n", "Nr",
		"Name",
		"Section",
		"Value",
		"Size",
		"Inf",
		"Oth")
	for i, symbol := range symbols {
		fmt.Printf("  [%2d] %-30s %-10s %08x %08x %3d %3d\n", i,
			symbol.Name,
			strings.Replace(symbol.Section.String(), "SHN_", "", -1),
			symbol.Value,
			symbol.Size,
			symbol.Info,
			symbol.Other)
	}
}

func showImportedLibraries(libraries []string) {
	fmt.Printf("ELF ImportedLibraries:\n")
	fmt.Printf("  [%2s] %s\n", "Nr", "Name")
	for i, library := range libraries {
		fmt.Printf("  [%2d] %-30s\n", i, library)
	}
}

func showImportedSymbols(symbols []elf.ImportedSymbol) {
	fmt.Printf("ELF ImportedSymbols:\n")
	fmt.Printf("  [%2s] %-30s %-10s %s\n", "Nr",
		"Name",
		"Version",
		"Library")
	for i, symbol := range symbols {
		fmt.Printf("  [%2d] %-30s %-10s %s\n", i,
			symbol.Name,
			symbol.Version,
			symbol.Library)
	}
}
