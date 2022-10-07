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
			off := i - 1
			c, w := b.readConst(i)
			i += w
			buf.WriteString(fmt.Sprintf("[%d] OP_CONST %v\n", off, c))

		case OP_LOAD:
			off := i - 1
			c, w := b.readConst(i)
			i += w
			buf.WriteString(fmt.Sprintf("[%d] OP_LOAD %s\n", off, c.(string)))

		case OP_ADD_NUM:
			off := i - 1
			buf.WriteString(fmt.Sprintf("[%d] OP_ADD_NUM\n", off))

		case OP_ADD_NUM_NUM:
			off := i - 1
			buf.WriteString(fmt.Sprintf("[%d] OP_ADD_NUM_NUM\n", off))

		case OP_SUB_NUM:
			off := i - 1
			buf.WriteString(fmt.Sprintf("[%d] OP_SUB_NUM\n", off))

		case OP_SUB_NUM_NUM:
			off := i - 1
			buf.WriteString(fmt.Sprintf("[%d] OP_SUB_NUM_NUM\n", off))

		case OP_MUL_NUM_NUM:
			off := i - 1
			buf.WriteString(fmt.Sprintf("[%d] OP_MUL_NUM_NUM\n", off))

		case OP_DIV_NUM_NUM:
			off := i - 1
			buf.WriteString(fmt.Sprintf("[%d] OP_DIV_NUM_NUM\n", off))

		case OP_MOD_NUM_NUM:
			off := i - 1
			buf.WriteString(fmt.Sprintf("[%d] OP_MOD_NUM_NUM\n", off))

		case OP_EXP_NUM_NUM:
			off := i - 1
			buf.WriteString(fmt.Sprintf("[%d] OP_EXP_NUM_NUM\n", off))

		case OP_IF:
			off := i - 1
			tOff, sz := b.readMediumInt(i)
			i += sz
			fOff, sz := b.readMediumInt(i)
			i += sz
			buf.WriteString(fmt.Sprintf("[%d] OP_IF %d, %d\n", off, tOff, fOff))

		case OP_JUMP:
			off := i - 1
			target, sz := b.readMediumInt(i)
			i += sz
			buf.WriteString(fmt.Sprintf("[%d] OP_JUMP %d\n", off, target))

		case OP_NEW_LIST:
			off := i - 1
			kd, w := b.readConst(i)
			i += w
			sz, w := b.readMediumInt(i)
			i += w
			buf.WriteString(fmt.Sprintf("[%d] OP_NEW_LIST %s, %d\n", off, kd, sz))

		case OP_NEW_MAP:
			off := i - 1
			kd, w := b.readConst(i)
			i += w
			sz, w := b.readMediumInt(i)
			i += w
			buf.WriteString(fmt.Sprintf("[%d] OP_NEW_MAP %s, %d\n", off, kd, sz))

		case OP_NEW_OBJ:
			off := i - 1
			kd, w := b.readConst(i)
			i += w
			buf.WriteString(fmt.Sprintf("[%d] OP_NEW_OBJ %s\n", off, kd))

		case OP_LIST_LOAD:
			off := i - 1
			buf.WriteString(fmt.Sprintf("[%d] OP_LIST_LOAD\n", off))

		case OP_MAP_LOAD:
			off := i - 1
			buf.WriteString(fmt.Sprintf("[%d] OP_MAP_LOAD\n", off))

		case OP_OBJ_LOAD:
			off := i - 1
			sz, w := b.readMediumInt(i)
			i += w
			buf.WriteString(fmt.Sprintf("[%d] OP_OBJ_LOAD %d\n", off, sz))

		case OP_INVOKE_STATIC:
			off := i - 1
			f, w := b.readConst(i)
			i += w
			argc := b.readUint8(i)
			i += 1
			buf.WriteString(fmt.Sprintf("[%d] OP_INVOKE_STATIC %s, %d\n", off, f, argc))

		case OP_INVOKE_STATIC_LAZY:
			off := i - 1
			f, w := b.readConst(i)
			i += w
			argc := b.readUint8(i)
			i += 1
			buf.WriteString(fmt.Sprintf("[%d] OP_INVOKE_STATIC_LAZY %s, %d\n", off, f, argc))

		case OP_INVOKE_DYNAMIC:
			off := i - 1
			argc := b.readUint8(i)
			i += 1
			buf.WriteString(fmt.Sprintf("[%d] OP_INVOKE_DYNAMIC %d\n", off, argc))

		case OP_NOP:
			off := i - 1
			buf.WriteString(fmt.Sprintf("[%d] OP_NOP\n", off))
		case OP_RETURN:
			off := i - 1
			buf.WriteString(fmt.Sprintf("[%d] OP_RETURN\n", off))

		default:
			off := i - 1
			panic(fmt.Errorf("[%d] unsupported opcode %d", off, opcode))
		}

	}
	return buf.String()
}
