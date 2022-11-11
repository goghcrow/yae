package vm

import (
	"bytes"
	"fmt"
)

//goland:noinspection GoUnhandledErrorResult
func (b *bytecode) String() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "{\n")
	fmt.Fprintf(&buf, "\tdata: %v\n", b.data)
	fmt.Fprintf(&buf, "\tcode:\n")
	disassemble(b, &buf)
	fmt.Fprintf(&buf, "}")
	return buf.String()
}

func Disassemble(b *bytecode) string {
	var buf bytes.Buffer
	disassemble(b, &buf)
	return buf.String()
}

//goland:noinspection GoUnhandledErrorResult
func disassemble(b *bytecode, buf *bytes.Buffer) {
	i := 0
	for {
		if i >= len(b.code) {
			break
		}

		op := opcode(b.code[i])
		i += 1

		switch op {
		case OP_CONST:
			fallthrough
		case OP_LOAD:
			off := i - 1
			c, w := b.readConst(i)
			i += w
			fmt.Fprintf(buf, "[%d] %s %v\n", off, op, c)
		case OP_IF_TRUE:
			fallthrough
		case OP_JUMP:
			off := i - 1
			target, sz := b.readMediumInt(i)
			i += sz
			fmt.Fprintf(buf, "[%d] %s %d\n", off, op, target)

		case OP_NEW_LIST:
			fallthrough
		case OP_NEW_MAP:
			off := i - 1
			ty, w := b.readConst(i)
			i += w
			sz, w := b.readMediumInt(i)
			i += w
			fmt.Fprintf(buf, "[%d] %s %s, %d\n", off, op, ty, sz)

		case OP_NEW_OBJ:
			off := i - 1
			ty, w := b.readConst(i)
			i += w
			fmt.Fprintf(buf, "[%d] %s %s\n", off, op, ty)

		case OP_OBJ_LOAD:
			off := i - 1
			sz, w := b.readMediumInt(i)
			i += w
			fmt.Fprintf(buf, "[%d] %s %d\n", off, op, sz)

		case OP_CALL_BY_VALUE:
			fallthrough
		case OP_CALL_BY_NEED:
			off := i - 1
			f, w := b.readConst(i)
			i += w
			argc := b.readUint8(i)
			i += 1
			fmt.Fprintf(buf, "[%d] %s %s, %d\n", off, op, f, argc)

		case OP_DYNAMIC_CALL:
			off := i - 1
			argc := b.readUint8(i)
			i += 1
			fmt.Fprintf(buf, "[%d] %s %d\n", off, op, argc)

		default:
			if op < _END_ {
				fmt.Fprintf(buf, "[%d] %s\n", i-1, op)
			} else {
				panic(fmt.Errorf("[%d] unsupported opcode %d", i-1, op))
			}
		}
	}
}
