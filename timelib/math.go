package timelib

import (
	"C"
	m "math"
)

// 解决链接 libm.a 报错 undefined symbol in cgo

//export fabs
func fabs(x C.double) C.double {
	return C.double(m.Abs(float64(x)))
}

//export acos
func acos(x C.double) C.double {
	return C.double(m.Acos(float64(x)))
}

////export acosf
//func acosf(x C.double) C.double {
//	return C.double(m.Acos(float64(x)))
//}

//export acosh
func acosh(x C.double) C.double {
	return C.double(m.Acosh(float64(x)))
}

//export asin
func asin(x C.double) C.double {
	return C.double(m.Asin(float64(x)))
}

//export asinh
func asinh(x C.double) C.double {
	return C.double(m.Asinh(float64(x)))
}

//export atan
func atan(x C.double) C.double {
	return C.double(m.Atan(float64(x)))
}

//export atanh
func atanh(x C.double) C.double {
	return C.double(m.Atanh(float64(x)))
}

func atan2(y C.double, x C.double) C.double {
	return C.double(m.Atan2(float64(y), float64(x)))
}

//export cbrt
func cbrt(x C.double) C.double {
	return C.double(m.Cbrt(float64(x)))
}

//export ceil
func ceil(x C.double) C.double {
	return C.double(m.Ceil(float64(x)))
}

//export copysign
func copysign(x C.double, y C.double) C.double {
	return C.double(m.Copysign(float64(x), float64(y)))
}

//export cos
func cos(x C.double) C.double {
	return C.double(m.Cos(float64(x)))
}

//export cosh
func cosh(x C.double) C.double {
	return C.double(m.Cosh(float64(x)))
}

//export fdim
func fdim(x C.double, y C.double) C.double {
	return C.double(m.Dim(float64(x), float64(y)))
}

//export erf
func erf(x C.double) C.double {
	return C.double(m.Erf(float64(x)))
}

//export erfc
func erfc(x C.double) C.double {
	return C.double(m.Erfc(float64(x)))
}

//export exp
func exp(x C.double) C.double {
	return C.double(m.Exp(float64(x)))
}

//export exp2
func exp2(x C.double) C.double {
	return C.double(m.Exp2(float64(x)))
}

//export expm1
func expm1(x C.double) C.double {
	return C.double(m.Expm1(float64(x)))
}

//export floor
func floor(x C.double) C.double {
	return C.double(m.Floor(float64(x)))
}

func frexp(f C.double, expptr *C.int) C.double {
	frac, exp := m.Frexp(float64(f))
	*expptr = C.int(exp)
	return C.double(frac)
}

//export log
func log(x C.double) C.double {
	return C.double(m.Log(float64(x)))
}

//export log10
func log10(x C.double) C.double {
	return C.double(m.Log10(float64(x)))
}

//export fmax
func fmax(x C.double, y C.double) C.double {
	return C.double(m.Max(float64(x), float64(y)))
}

//export fmin
func fmin(x C.double, y C.double) C.double {
	return C.double(m.Min(float64(x), float64(y)))
}

//export pow
func pow(x C.double, y C.double) C.double {
	return C.double(m.Pow(float64(x), float64(y)))
}

//export round
func round(x C.double) C.double {
	return C.double(m.Round(float64(x)))
}

//export sqrt
func sqrt(x C.double) C.double {
	return C.double(m.Sqrt(float64(x)))
}

//export sin
func sin(x C.double) C.double {
	return C.double(m.Sin(float64(x)))
}

//export sinh
func sinh(x C.double) C.double {
	return C.double(m.Sinh(float64(x)))
}

//export tan
func tan(x C.double) C.double {
	return C.double(m.Tan(float64(x)))
}

//export tanh
func tanh(x C.double) C.double {
	return C.double(m.Tanh(float64(x)))
}
