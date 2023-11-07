#!/bin/bash

rm -rf bin
GOOS=windows GOARCH=amd64 go build -o bin/windows-x86-64/elfquery.exe main.go
GOOS=windows GOARCH=arm64 go build -o bin/windows-aarch64/elfquery.exe main.go
GOOS=darwin GOARCH=amd64 go build -o bin/darwin-x86-64/elfquery main.go
GOOS=darwin GOARCH=arm64 go build -o bin/darwin-aarch64/elfquery main.go
GOOS=linux GOARCH=386 go build -o bin/linux-386/elfquery main.go
GOOS=linux GOARCH=amd64 go build -o bin/linux-x86-64/elfquery main.go
GOOS=linux GOARCH=arm go build -o bin/linux-aarch32/elfquery main.go
GOOS=linux GOARCH=arm64 go build -o bin/linux-aarch64/elfquery main.go
