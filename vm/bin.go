package vm

import (
	"encoding/binary"
	"math"
	"unsafe"
)

func uint16ToByte(i uint16) []byte {
	var buf [2]byte
	binary.BigEndian.PutUint16(buf[:], i)
	return buf[:]
}

func byteToUInt16(buf []byte) uint16 {
	return binary.BigEndian.Uint16(buf)
}

func uint8ToByte(i uint8) []byte {
	return []byte{i}
}

func byteToUInt8(buf []byte) uint8 {
	return buf[0]
}

func int64ToByte(i int64) []byte {
	var buf [8]byte
	ui64 := *(*uint64)(unsafe.Pointer(&i))
	binary.BigEndian.PutUint64(buf[:], ui64)
	return buf[:]
}

func byteToInt64(buf []byte) int64 {
	ui64 := binary.BigEndian.Uint64(buf)
	return *(*int64)(unsafe.Pointer(&ui64))
}

func float64ToByte(f float64) []byte {
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], math.Float64bits(f))
	return buf[:]
}

func byteToFloat64(buf []byte) float64 {
	ui64 := binary.BigEndian.Uint64(buf)
	return *(*float64)(unsafe.Pointer(&ui64))
}
