package protocol

import (
	"encoding/binary"
	"math"
)

type ByteOrder interface {
	Uint16([]byte) uint16
	Uint24([]byte) uint32
	Uint32([]byte) uint32
	Uint40([]byte) uint64
	Uint48([]byte) uint64
	Uint56([]byte) uint64
	Uint64([]byte) uint64
	Float32([]byte) float32
	Float64([]byte) float64
	PutUint16([]byte, uint16)
	PutUint24([]byte, uint32)
	PutUint32([]byte, uint32)
	PutUint40([]byte, uint64)
	PutUint48([]byte, uint64)
	PutUint56([]byte, uint64)
	PutUint64([]byte, uint64)
	PutFloat32([]byte, float32)
	PutFloat64([]byte, float64)
}

// LittleEndian is the little-endian implementation of ByteOrder.
var LittleEndian littleEndian

// BigEndian is the big-endian implementation of ByteOrder.
var BigEndian bigEndian

type bigEndian struct{}

func (bigEndian) Uint16(b []byte) uint16 {
	return binary.BigEndian.Uint16(b)
}

func (bigEndian) Uint24(b []byte) uint32 {
	return uint32(b[2]) | uint32(b[1])<<8 | uint32(b[0])<<16
}

func (bigEndian) Uint32(b []byte) uint32 {
	return binary.BigEndian.Uint32(b)
}

func (bigEndian) Uint40(b []byte) uint64 {
	return uint64(b[4]) | uint64(b[3])<<8 | uint64(b[2])<<16 | uint64(b[1])<<24 |
		uint64(b[0])<<32
}

func (bigEndian) Uint48(b []byte) uint64 {
	return uint64(b[5]) | uint64(b[4])<<8 | uint64(b[3])<<16 | uint64(b[2])<<24 |
		uint64(b[1])<<32 | uint64(b[0])<<40
}

func (bigEndian) Uint56(b []byte) uint64 {
	return uint64(b[6]) | uint64(b[5])<<8 | uint64(b[4])<<16 | uint64(b[3])<<24 |
		uint64(b[2])<<32 | uint64(b[1])<<40 | uint64(b[0])<<48
}

func (bigEndian) Uint64(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
}

func (e bigEndian) Float32(b []byte) float32 {
	return math.Float32frombits(e.Uint32(b))
}

func (e bigEndian) Float64(b []byte) float64 {
	return math.Float64frombits(e.Uint64(b))
}

func (bigEndian) PutUint16(b []byte, v uint16) {
	binary.BigEndian.PutUint16(b, v)
}

func (bigEndian) PutUint24(b []byte, v uint32) {
	b[0] = byte(v >> 16)
	b[1] = byte(v >> 8)
	b[2] = byte(v)

}

func (bigEndian) PutUint32(b []byte, v uint32) {
	binary.BigEndian.PutUint32(b, v)
}

func (bigEndian) PutUint40(b []byte, v uint64) {
	b[0] = byte(v >> 32)
	b[1] = byte(v >> 24)
	b[2] = byte(v >> 16)
	b[3] = byte(v >> 8)
	b[4] = byte(v)
}

func (bigEndian) PutUint48(b []byte, v uint64) {
	b[0] = byte(v >> 40)
	b[1] = byte(v >> 32)
	b[2] = byte(v >> 24)
	b[3] = byte(v >> 16)
	b[4] = byte(v >> 8)
	b[5] = byte(v)
}

func (bigEndian) PutUint56(b []byte, v uint64) {
	b[0] = byte(v >> 48)
	b[1] = byte(v >> 40)
	b[2] = byte(v >> 32)
	b[3] = byte(v >> 24)
	b[4] = byte(v >> 16)
	b[5] = byte(v >> 8)
	b[6] = byte(v)
}

func (bigEndian) PutUint64(b []byte, v uint64) {
	binary.BigEndian.PutUint64(b, v)
}

func (e bigEndian) PutFloat32(b []byte, v float32) {
	e.PutUint32(b, math.Float32bits(v))
}

func (e bigEndian) PutFloat64(b []byte, v float64) {
	e.PutUint64(b, math.Float64bits(v))
}

type littleEndian struct{}

func (littleEndian) Uint16(b []byte) uint16 {
	return binary.LittleEndian.Uint16(b)
}

func (littleEndian) Uint24(b []byte) uint32 {
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16
}

func (littleEndian) Uint32(b []byte) uint32 {
	return binary.LittleEndian.Uint32(b)
}

func (littleEndian) Uint40(b []byte) uint64 {
	return uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 |
		uint64(b[4])<<32
}

func (littleEndian) Uint48(b []byte) uint64 {
	return uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 |
		uint64(b[4])<<32 | uint64(b[5])<<40
}

func (littleEndian) Uint56(b []byte) uint64 {
	return uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 |
		uint64(b[4])<<32 | uint64(b[5])<<40 | uint64(b[6])<<48
}

func (littleEndian) Uint64(b []byte) uint64 {
	return binary.LittleEndian.Uint64(b)
}

func (e littleEndian) Float32(b []byte) float32 {
	return math.Float32frombits(e.Uint32(b))
}

func (e littleEndian) Float64(b []byte) float64 {
	return math.Float64frombits(e.Uint64(b))
}

func (littleEndian) PutUint16(b []byte, v uint16) {
	binary.LittleEndian.PutUint16(b, v)
}

func (littleEndian) PutUint24(b []byte, v uint32) {
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
}

func (littleEndian) PutUint32(b []byte, v uint32) {
	binary.LittleEndian.PutUint32(b, v)
}

func (littleEndian) PutUint40(b []byte, v uint64) {
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
	b[3] = byte(v >> 24)
	b[4] = byte(v >> 32)
}

func (littleEndian) PutUint48(b []byte, v uint64) {
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
	b[3] = byte(v >> 24)
	b[4] = byte(v >> 32)
	b[5] = byte(v >> 40)
}

func (littleEndian) PutUint56(b []byte, v uint64) {
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
	b[3] = byte(v >> 24)
	b[4] = byte(v >> 32)
	b[5] = byte(v >> 40)
	b[6] = byte(v >> 48)
}

func (littleEndian) PutUint64(b []byte, v uint64) {
	binary.LittleEndian.PutUint64(b, v)
}

func (e littleEndian) PutFloat32(b []byte, v float32) {
	e.PutUint32(b, math.Float32bits(v))
}

func (e littleEndian) PutFloat64(b []byte, v float64) {
	e.PutUint64(b, math.Float64bits(v))
}
