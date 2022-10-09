package val

import "math"

const epsilon = 1e-9

// 👇🏻实现的 EQ 是不正确的, 正确的实现参考 https://floating-point-gui.de/errors/comparison/

func NumEQ(x, y *NumVal) bool { return math.Abs(x.V-y.V) < epsilon }
func NumNE(x, y *NumVal) bool { return math.Abs(x.V-y.V) >= epsilon }
func NumLT(x, y *NumVal) bool { return x.V < y.V && NumNE(x, y) }
func NumLE(x, y *NumVal) bool { return x.V <= y.V || NumEQ(x, y) }
func NumGT(x, y *NumVal) bool { return x.V > y.V && NumNE(x, y) }
func NumGE(x, y *NumVal) bool { return x.V >= y.V || NumEQ(x, y) }
