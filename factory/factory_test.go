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

func TestStringSlice(t *testing.T) {
	s := StringSlice()
	assert.NotNil(t, s)
}

func TestBool(t *testing.T) {
	_ = Bool()
}

func TestBoolPtr(t *testing.T) {
	b := BoolPtr()
	assert.NotNil(t, b)
}

func TestBoolSlice(t *testing.T) {
	b := BoolSlice()
	assert.NotNil(t, b)
}

func TestInt(t *testing.T) {
	i := Int()
	assert.NotZero(t, i)
}

func TestIntPtr(t *testing.T) {
	i := IntPtr()
	assert.NotNil(t, i)
}

func TestIntSlice(t *testing.T) {
	i := IntSlice()
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

func TestInt8Slice(t *testing.T) {
	i := Int8Slice()
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

func TestInt16Slice(t *testing.T) {
	i := Int16Slice()
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

func TestInt32Slice(t *testing.T) {
	i := Int32Slice()
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

func TestInt64Slice(t *testing.T) {
	i := Int64Slice()
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

func TestUintSlice(t *testing.T) {
	u := UintSlice()
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

func TestUint8Slice(t *testing.T) {
	u := Uint8Slice()
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

func TestUint16Slice(t *testing.T) {
	u := Uint16Slice()
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

func TestUint32Slice(t *testing.T) {
	u := Uint32Slice()
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

func TestUint64Slice(t *testing.T) {
	u := Uint64Slice()
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

func TestFloat32Slice(t *testing.T) {
	f := Float32Slice()
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

func TestFloat64Slice(t *testing.T) {
	f := Float64Slice()
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

func TestComplex64Slice(t *testing.T) {
	c := Complex64Slice()
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

func TestComplex128Slice(t *testing.T) {
	c := Complex128Slice()
	assert.NotNil(t, c)
}

func TestByte(t *testing.T) {
	b := Byte()
	assert.NotZero(t, b)
}

func TestBytePtr(t *testing.T) {
	b := BytePtr()
	assert.NotNil(t, b)
}

func TestByteSlice(t *testing.T) {
	b := ByteSlice()
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

func TestRuneSlice(t *testing.T) {
	r := RuneSlice()
	assert.NotNil(t, r)
}

func TestDuration(t *testing.T) {
	d := Duration()
	assert.NotZero(t, d)
}

func TestDurationPtr(t *testing.T) {
	d := DurationPtr()
	assert.NotNil(t, d)
}

func TestDurationSlice(t *testing.T) {
	d := DurationSlice()
	assert.NotNil(t, d)
}
