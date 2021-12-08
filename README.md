# bump

An efficient buffering library for Go (Golang).

[![](https://img.shields.io/badge/status-1.0.0-ff00bb.svg?style=flat-square)](https://github.com/surrealdb/bump) [![](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/surrealdb/bump) [![](https://goreportcard.com/badge/github.com/surrealdb/bump?style=flat-square)](https://goreportcard.com/report/github.com/surrealdb/bump) [![](https://img.shields.io/badge/license-Apache_License_2.0-00bfff.svg?style=flat-square)](https://github.com/surrealdb/bump) 

#### Features

- Simple and efficient buffering
- Reuse readers and writers repeatedly
- Write to io.Writer, or directly to a byte slice
- Read from io.Reader, or directly from a byte slice
- Reading directly from byte slice requires no allocations

#### Installation

```bash
go get github.com/surrealdb/bump
```
