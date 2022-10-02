#!/bin/bash

echo -e "package fun\n" > funs.go
echo -e "import \"github.com/goghcrow/yae/val\"\n" >> funs.go
echo "var Funs = []*val.Val{" >> funs.go
egrep -r -h -o --include=*.go "[A-Z][A-Z_]+ = " .| awk -F " " '{print "\t" $1 ","}' >> funs.go
echo '}' >> funs.go

