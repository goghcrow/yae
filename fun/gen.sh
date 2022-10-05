#!/bin/bash

echo -e "// Generated Code; DO NOT EDIT.\n" > funs.go
echo -e "package fun\n" >> funs.go
echo -e "import \"github.com/goghcrow/yae/val\"\n" >> funs.go
echo "var funs = []*val.Val{" >> funs.go
egrep -r -h -o --include=*.go --exclude=const.go "[A-Z][A-Z_]+ = " .| awk -F " " '{print "\t" $1 ","}' >> funs.go
echo -e "}\n" >> funs.go

echo -e "func BuildIn() []*val.Val {" >> funs.go
echo -e "\treturn funs" >> funs.go
echo -e "}\n" >> funs.go


#egrep -r -h -o --include=*.go "// [A-Z][A-Z_]+ .+" . | awk '{$1=""}1' | awk '{$1=""}1' |awk '{$1=$1}1' | sort > signature.txt