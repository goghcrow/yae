#!/bin/bash

echo "package fun" > gen_test.go
echo "import \"github.com/goghcrow/yae/val\"" >> gen_test.go
echo "var funs = []*val.Val{" >> gen_test.go
egrep -r -h -o --include=*.go "[A-Z][A-Z_]+ = " .| awk -F " " '{print $1 ","}' >> gen_test.go
echo '}' >> gen_test.go

