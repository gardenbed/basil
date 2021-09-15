package factory

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	s := String()
	assert.NotEmpty(t, s)
	assert.Len(t, s, stringLen)
}

func TestStringPtr(t *testing.T) {
	s := StringPtr()
	assert.NotNil(t, s)
}

func TestBool(t *testing.T) {
	b := Bool()
	assert.NotZero(t, b)
}

func TestBoolPtr(t *testing.T) {
	b := BoolPtr()
	assert.NotNil(t, b)
}

func TestByte(t *testing.T) {
	b := Byte()
	assert.NotZero(t, b)
}

func TestBytePtr(t *testing.T) {
	b := BytePtr()
	assert.NotNil(t, b)
}

func TestRune(t *testing.T) {
	r := Rune()
	assert.NotZero(t, r)
}

func TestRunePtr(t *testing.T) {
	r := RunePtr()
	assert.NotNil(t, r)
}

func TestInt(t *testing.T) {
	i := Int()
	assert.NotZero(t, i)
}

func TestIntPtr(t *testing.T) {
	i := IntPtr()
	assert.NotNil(t, i)
}

func TestInt8(t *testing.T) {
	i := Int8()
	assert.NotZero(t, i)
}

func TestInt8Ptr(t *testing.T) {
	i := Int8Ptr()
	assert.NotNil(t, i)
}

func TestInt16(t *testing.T) {
	i := Int16()
	assert.NotZero(t, i)
}

func TestInt16Ptr(t *testing.T) {
	i := Int16Ptr()
	assert.NotNil(t, i)
}

func TestInt32(t *testing.T) {
	i := Int32()
	assert.NotZero(t, i)
}

func TestInt32Ptr(t *testing.T) {
	i := Int32Ptr()
	assert.NotNil(t, i)
}

func TestInt64(t *testing.T) {
	i := Int64()
	assert.NotZero(t, i)
}

func TestInt64Ptr(t *testing.T) {
	i := Int64Ptr()
	assert.NotNil(t, i)
}

func TestUint(t *testing.T) {
	u := Uint()
	assert.NotZero(t, u)
}

func TestUintPtr(t *testing.T) {
	u := UintPtr()
	assert.NotNil(t, u)
}

func TestUint8(t *testing.T) {
	u := Uint8()
	assert.NotZero(t, u)
}

func TestUint8Ptr(t *testing.T) {
	u := Uint8Ptr()
	assert.NotNil(t, u)
}

func TestUint16(t *testing.T) {
	u := Uint16()
	assert.NotZero(t, u)
}

func TestUint16Ptr(t *testing.T) {
	u := Uint16Ptr()
	assert.NotNil(t, u)
}

func TestUint32(t *testing.T) {
	u := Uint32()
	assert.NotZero(t, u)
}

func TestUint32Ptr(t *testing.T) {
	u := Uint32Ptr()
	assert.NotNil(t, u)
}

func TestUint64(t *testing.T) {
	u := Uint64()
	assert.NotZero(t, u)
}

func TestUint64Ptr(t *testing.T) {
	u := Uint64Ptr()
	assert.NotNil(t, u)
}

func TestUintptr(t *testing.T) {
	u := Uintptr()
	assert.NotZero(t, u)
}

func TestUintptrPtr(t *testing.T) {
	u := UintptrPtr()
	assert.NotNil(t, u)
}

func TestFloat32(t *testing.T) {
	f := Float32()
	assert.NotZero(t, f)
}

func TestFloat32Ptr(t *testing.T) {
	f := Float32Ptr()
	assert.NotNil(t, f)
}

func TestFloat64(t *testing.T) {
	f := Float64()
	assert.NotZero(t, f)
}

func TestFloat64Ptr(t *testing.T) {
	f := Float64Ptr()
	assert.NotNil(t, f)
}

func TestComplex64(t *testing.T) {
	c := Complex64()
	assert.NotZero(t, c)
}

func TestComplex64Ptr(t *testing.T) {
	c := Complex64Ptr()
	assert.NotNil(t, c)
}

func TestComplex128(t *testing.T) {
	c := Complex128()
	assert.NotZero(t, c)
}

func TestComplex128Ptr(t *testing.T) {
	c := Complex128Ptr()
	assert.NotNil(t, c)
}

func TestError(t *testing.T) {
	e := Error()
	assert.Error(t, e)
	assert.NotEmpty(t, e.Error())
}

func TestDuration(t *testing.T) {
	d := Duration()
	assert.NotZero(t, d)
}

func TestDurationPtr(t *testing.T) {
	d := DurationPtr()
	assert.NotNil(t, d)
}
