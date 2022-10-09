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
	buf := &strings.Builder{}
	i := 0
	for {
		if i >= len(b.code) {
			break
		}

		opcode := op(b.code[i])
		i += 1

		switch opcode {
		case OP_CONST:
			fallthrough
		case OP_LOAD:
			off := i - 1
			c, w := b.readConst(i)
			i += w
			buf.WriteString(fmt.Sprintf("[%d] %s %v\n", off, opcode, c))

		case OP_IF_TRUE:
			fallthrough
		case OP_JUMP:
			off := i - 1
			target, sz := b.readMediumInt(i)
			i += sz
			buf.WriteString(fmt.Sprintf("[%d] %s %d\n", off, opcode, target))

		case OP_NEW_LIST:
			fallthrough
		case OP_NEW_MAP:
			off := i - 1
			kd, w := b.readConst(i)
			i += w
			sz, w := b.readMediumInt(i)
			i += w
			buf.WriteString(fmt.Sprintf("[%d] %s %s, %d\n", off, opcode, kd, sz))

		case OP_NEW_OBJ:
			off := i - 1
			kd, w := b.readConst(i)
			i += w
			buf.WriteString(fmt.Sprintf("[%d] %s %s\n", off, opcode, kd))

		case OP_OBJ_LOAD:
			off := i - 1
			sz, w := b.readMediumInt(i)
			i += w
			buf.WriteString(fmt.Sprintf("[%d] %s %d\n", off, opcode, sz))

		case OP_INVOKE_STATIC:
			fallthrough
		case OP_INVOKE_STATIC_LAZY:
			off := i - 1
			f, w := b.readConst(i)
			i += w
			argc := b.readUint8(i)
			i += 1
			buf.WriteString(fmt.Sprintf("[%d] %s %s, %d\n", off, opcode, f, argc))

		case OP_INVOKE_DYNAMIC:
			off := i - 1
			argc := b.readUint8(i)
			i += 1
			buf.WriteString(fmt.Sprintf("[%d] %s %d\n", off, opcode, argc))

		default:
			if opcode < _END_ {
				buf.WriteString(fmt.Sprintf("[%d] %s\n", i-1, opcode))
			} else {
				panic(fmt.Errorf("[%d] unsupported opcode %d", i-1, opcode))
			}
		}

	}
	return buf.String()
}
