# Beignet

Converts .dylib files to MacOS shellcode.

## CLI

Build:

`go build -o beignet ./cli`

Convert a dylib to a raw shellcode buffer:

`./beignet --out payload.bin ./payload.dylib`

Print the embedded loader C source:

`./beignet dump-loader-c`

## Regenerating the embedded loader (darwin/arm64)

`GOCACHE=/tmp/go-build go generate ./internal/stager`
