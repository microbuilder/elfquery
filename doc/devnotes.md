# Development Notes

## Module Init

Create a `go.mod` file via:

```bash
$ go mod init github.com/microbuilder/elfquery
```

Point to the local repo, rather than the remote one, in `go.mod`:

```bash
$ go mod edit -replace github.com/microbuilder/elfquery=/Users/kevin/zendnode/ELF/elfquery
```

## Packages

```bash
$ go get github.com/spf13/cobra/cobra
$ go get github.com/yalue/elf_reader
$ go get github.com/mattn/go-sqlite3
$ go get github.com/microbuilder/elfquery
```

## Cobra

Initialise cobra in the app via:

```bash
cobra -a "Kevin Townsend <kevin@ktownsend.com>" init --pkg-name elfquery
```

Add commands via:

```bash
cobra -a "Kevin Townsend <kevin@ktownsend.com>" add <commandName>
```

## Redirect modules to local code

