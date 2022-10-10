#!/bin/bash

echo -e "// Generated Code; DO NOT EDIT.\n" > gen.go
echo -e "package fun\n" >> gen.go
echo -e "import \"github.com/goghcrow/yae/val\"\n" >> gen.go
echo -e "//go:generate /bin/bash gen.sh\n" >> gen.go
echo "var funs = []*val.Val{" >> gen.go
egrep -r -h -o --include=*.go --exclude=const.go "[A-Z][A-Z_]+ = " .| sort | awk -F " " '{print "\t" $1 ","}' >> gen.go
echo -e "}\n" >> gen.go

echo -e "func BuildIn() []*val.Val {" >> gen.go
echo -e "\treturn funs" >> gen.go
echo -e "}\n" >> gen.go


#egrep -r -h -o --include=*.go "// [A-Z][A-Z_]+ .+" . | awk '{$1=""}1' | awk '{$1=""}1' |awk '{$1=$1}1' | sort > signature.txt