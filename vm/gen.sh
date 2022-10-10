#!/bin/bash

echo -e "// Generated Code; DO NOT EDIT.\n" > gen.go
echo -e "package vm\n" >> gen.go
echo "func (o opcode) String() string {" >> gen.go
echo -e "\treturn [...]string {" >> gen.go
egrep -r -h -o "OP_\w+( )?" opcode.go | awk -F " " '{print "\t\t\"" $1 "\","}' >> gen.go
echo -e "\t}[o]" >> gen.go
echo -e "}\n" >> gen.go