package vm

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"testing"
)

func TestGenOpcodeString(t *testing.T) {
	cmd := exec.Command("/bin/bash", "gen.sh")
	_, _ = cmd.Output()
}

func TestGenCallThreadingBySwitch(t *testing.T) {
	buf := ""
	file, _ := os.ReadFile("switchthread.go")
	// `(?s)case (OP_\w+):\n(.*?)[\n\t]*(?=case OP_)` // 不支持正向预查
	pattern := regexp.MustCompile(`(?s)case (OP_\w+):\n(.*?)[\n\t]*(?:case OP_|// --------)`)
	str := string(file)

	var ops []string
	for {
		s := pattern.FindStringSubmatch(str)
		if s == nil {
			break
		}
		if s[1] != "OP_RETURN" {
			ops = append(ops, s[1])
			buf += fmt.Sprintf(
				`//goland:noinspection GoSnakeCaseUsage
func %s_Handler(v *VM) {
%s
}

`, s[1], strings.Replace(s[2], "\t\t\t", "\t", -1))
		}

		loc := pattern.FindStringIndex(str)
		if loc == nil {
			break
		}
		str = str[loc[1]-len("case OP_"):]
	}

	buf += `//goland:noinspection GoSnakeCaseUsage
func OP_NOP_Handler(v *VM) {}`

	opBuf := ""
	for _, op := range ops {
		opBuf += fmt.Sprintf("\tinstructions[%s] = %s_Handler\n", op, op)
	}

	f := fmt.Sprintf(`// Generated Code; DO NOT EDIT.

package vm

import (
	"github.com/goghcrow/yae/timelib"
	"github.com/goghcrow/yae/types"
	"github.com/goghcrow/yae/util"
	"github.com/goghcrow/yae/val"
	"math"
	"time"
	"unicode/utf8"
	"unsafe"
)

type OpcodeHandler func(vm *VM)

var instructions [_END_]OpcodeHandler

const limit = 1024

func callThreading(v *VM) *val.Val {
	for i := 0; i < limit; i++ {
		op := opcode(v.code[v.pc])
		v.pc += 1
		if op == OP_RETURN {
			return v.Pop()
		}
		instructions[op](v)
	}
	panic("over exec limit")
}

func init() {
%s
	instructions[OP_NOP] = OP_NOP_Handler
}

%s
`, opBuf, buf)

	_ = os.WriteFile("callthread.go", []byte(f), 0644)
}
