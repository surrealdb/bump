# bump

An efficient buffering library for Go (Golang).

[![](https://img.shields.io/circleci/token/40ded060b367ddedc5e98f14d0553e6e380af543/project/abcum/bump/master.svg?style=flat-square)](https://circleci.com/gh/abcum/bump) [![](https://img.shields.io/badge/status-alpha-ff00bb.svg?style=flat-square)](https://github.com/abcum/bump) [![](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/abcum/bump) [![](https://goreportcard.com/badge/github.com/abcum/bump?style=flat-square)](https://goreportcard.com/report/github.com/abcum/bump) [![](https://img.shields.io/coveralls/abcum/bump/master.svg?style=flat-square)](https://coveralls.io/github/abcum/bump?branch=master) [![](https://img.shields.io/badge/license-Apache_License_2.0-00bfff.svg?style=flat-square)](https://github.com/abcum/bump) 

#### Features

- Simple and efficient buffering
- Reuse readers and writers repeatedly
- Write to io.Writer, or directly to a byte slice
- Read from io.Reader, or directly from a byte slice
- Reading directly from byte slice requires no allocations

#### Installation

```bash
go get github.com/abcum/bump
```
