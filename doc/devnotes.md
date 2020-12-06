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
$ go get github.com/goreleaser/goreleaser
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

## goreleaser

### Initialisation

Init goreleaser once:

```bash
$ goreleaser init
   • Generating .goreleaser.yml file
   • config created; please edit accordingly to your needs file=.goreleaser.yml
```

### Tagging a release

Create a new tag using semantic versioning:

```bash
$ git tag -a v0.1.0 -m "First release"
$ git push origin v0.1.0
```

### Build a release

Test build via:

```bash
$ goreleaser --snapshot --skip-publish --rm-dist
```

Submit a new release to Github's release page:

```bash
$ goreleaser
```
