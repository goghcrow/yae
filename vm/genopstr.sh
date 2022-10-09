#!/bin/bash

echo -e "// Generated Code; DO NOT EDIT.\n" > string.go
echo -e "package vm\n" >> string.go
echo "func (o op) String() string {" >> string.go
echo -e "\treturn [...]string {" >> string.go
egrep -r -h -o "OP_\w+( )?" opcode.go | awk -F " " '{print "\t\t\"" $1 "\","}' >> string.go
echo -e "\t}[o]" >> string.go
echo -e "}\n" >> string.go