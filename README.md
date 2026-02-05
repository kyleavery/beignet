# Beignet

Converts .dylib files to MacOS shellcode.

## CLI

Build:

`make`

Convert a dylib to a raw shellcode buffer:

`./beignet --out payload.bin ./payload.dylib`

## Regenerating the embedded loader (darwin/arm64)

`go generate ./internal/stager`
