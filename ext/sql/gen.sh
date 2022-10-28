#!/bin/bash

{
  echo -e "// Generated Code; DO NOT EDIT.\n"
  echo -e "package sql\n"
  echo -e "import \"github.com/goghcrow/yae/val\"\n"
  echo -e "//go:generate /bin/bash gen.sh\n"
  echo "var funs = []*val.Val{"
  egrep -r -h -o --include=*.go --exclude=const.go "[A-Z][A-Z_]+ = " .| sort | awk -F " " '{print "\t" $1 ","}'
  echo -e "}\n"
  echo -e "func BuiltIn() []*val.Val {"
  echo -e "\treturn funs"
  echo -e "}\n"
} > gen.go