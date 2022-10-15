package vm

import (
	"fmt"
	"strings"
)

func (b *bytecode) String() string {
	buf := "{\n"
	buf += fmt.Sprintf("\tdata: %v\n", b.data)
	buf += fmt.Sprintf("\tcode:\n%s", Disassemble(b))
	buf += "}"
	return buf
}

func Disassemble(b *bytecode) string {
	var buf strings.Builder
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
			buf.WriteString(fmt.Sprintf("[%d] %s %v\n", off, op, c))

		case OP_IF_TRUE:
			fallthrough
		case OP_JUMP:
			off := i - 1
			target, sz := b.readMediumInt(i)
			i += sz
			buf.WriteString(fmt.Sprintf("[%d] %s %d\n", off, op, target))

		case OP_NEW_LIST:
			fallthrough
		case OP_NEW_MAP:
			off := i - 1
			kd, w := b.readConst(i)
			i += w
			sz, w := b.readMediumInt(i)
			i += w
			buf.WriteString(fmt.Sprintf("[%d] %s %s, %d\n", off, op, kd, sz))

		case OP_NEW_OBJ:
			off := i - 1
			kd, w := b.readConst(i)
			i += w
			buf.WriteString(fmt.Sprintf("[%d] %s %s\n", off, op, kd))

		case OP_OBJ_LOAD:
			off := i - 1
			sz, w := b.readMediumInt(i)
			i += w
			buf.WriteString(fmt.Sprintf("[%d] %s %d\n", off, op, sz))

		case OP_CALL_BY_VALUE:
			fallthrough
		case OP_CALL_BY_NEED:
			off := i - 1
			f, w := b.readConst(i)
			i += w
			argc := b.readUint8(i)
			i += 1
			buf.WriteString(fmt.Sprintf("[%d] %s %s, %d\n", off, op, f, argc))

		case OP_DYNAMIC_CALL:
			off := i - 1
			argc := b.readUint8(i)
			i += 1
			buf.WriteString(fmt.Sprintf("[%d] %s %d\n", off, op, argc))

		default:
			if op < _END_ {
				buf.WriteString(fmt.Sprintf("[%d] %s\n", i-1, op))
			} else {
				panic(fmt.Errorf("[%d] unsupported opcode %d", i-1, op))
			}
		}

	}
	return buf.String()
}
