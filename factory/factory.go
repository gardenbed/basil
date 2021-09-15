// Package factory is used for generating random values for testing purposes.
package factory

import (
	"errors"
	"math/rand"
	"strings"
	"time"
)

const (
	stringLen = 64
	indexBits = 6                // 6 bits to represent a character index
	indexMax  = 63 / indexBits   // 10 character indices fits in 63 bits
	indexMask = 1<<indexBits - 1 // 111111 All 1-bits as many as indexBits
)

var (
	randSource = rand.NewSource(time.Now().UnixNano())
	chars      = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

// String generates a random string value.
func String() string {
	n := stringLen
	b := new(strings.Builder)
	b.Grow(n)

	// val is a randomly generated value.
	// count keeps track of how many more character indicies left in the current val.
	// l keeps track of the current length of the generated string.
	for val, count, ln := randSource.Int63(), indexMax, 0; ln < n; {
		if count == 0 {
			val, count = randSource.Int63(), indexMax
		}

		if i := int(val & indexMask); i < len(chars) {
			b.WriteByte(chars[i])
			ln++
		}

		val >>= indexBits
		count--
	}

	return b.String()
}

// StringPtr generates a random string value and returns a pointer to it.
func StringPtr() *string {
	s := String()
	return &s
}

// Bool generates a random bool value.
func Bool() bool {
	return rand.Intn(2) == 1
}

// BoolPtr generates a random bool value and returns a pointer to it.
func BoolPtr() *bool {
	b := Bool()
	return &b
}

// Byte generates a random byte value.
func Byte() byte {
	return byte(rand.Int63())
}

// BytePtr generates a random byte value and returns a pointer to it.
func BytePtr() *byte {
	b := Byte()
	return &b
}

// Rune generates a random rune value.
func Rune() rune {
	return rune(rand.Int63())
}

// RunePtr generates a random rune value and returns a pointer to it.
func RunePtr() *rune {
	r := Rune()
	return &r
}

// Int generates a random int value.
func Int() int {
	return int(rand.Int63())
}

// IntPtr generates a random int value and returns a pointer to it.
func IntPtr() *int {
	i := Int()
	return &i
}

// Int8 generates a random int8 value.
func Int8() int8 {
	return int8(rand.Int63())
}

// Int8Ptr generates a random int8 value and returns a pointer to it.
func Int8Ptr() *int8 {
	i := Int8()
	return &i
}

// Int16 generates a random int16 value.
func Int16() int16 {
	return int16(rand.Int63())
}

// Int16Ptr generates a random int16 value and returns a pointer to it.
func Int16Ptr() *int16 {
	i := Int16()
	return &i
}

// Int32 generates a random int32 value.
func Int32() int32 {
	return int32(rand.Int63())
}

// Int32Ptr generates a random int32 value and returns a pointer to it.
func Int32Ptr() *int32 {
	i := Int32()
	return &i
}

// Int64 generates a random int64 value.
func Int64() int64 {
	return int64(rand.Int63())
}

// Int64Ptr generates a random int64 value and returns a pointer to it.
func Int64Ptr() *int64 {
	i := Int64()
	return &i
}

// Uint generates a random uint value.
func Uint() uint {
	return uint(rand.Uint64())
}

// UintPtr generates a random uint value and returns a pointer to it.
func UintPtr() *uint {
	u := Uint()
	return &u
}

// Uint8 generates a random uint8 value.
func Uint8() uint8 {
	return uint8(rand.Uint64())
}

// Uint8Ptr generates a random uint8 value and returns a pointer to it.
func Uint8Ptr() *uint8 {
	u := Uint8()
	return &u
}

// Uint16 generates a random uint16 value.
func Uint16() uint16 {
	return uint16(rand.Uint64())
}

// Uint16Ptr generates a random uint16 value and returns a pointer to it.
func Uint16Ptr() *uint16 {
	u := Uint16()
	return &u
}

// Uint32 generates a random uint32 value.
func Uint32() uint32 {
	return rand.Uint32()
}

// Uint32Ptr generates a random uint32 value and returns a pointer to it.
func Uint32Ptr() *uint32 {
	u := Uint32()
	return &u
}

// Uint64 generates a random uint64 value.
func Uint64() uint64 {
	return rand.Uint64()
}

// Uint64Ptr generates a random uint64 value and returns a pointer to it.
func Uint64Ptr() *uint64 {
	u := Uint64()
	return &u
}

// Uintptr generates a random uintptr value.
func Uintptr() uintptr {
	return uintptr(rand.Uint64())
}

// UintptrPtr generates a random uintptr value and returns a pointer to it.
func UintptrPtr() *uintptr {
	u := Uintptr()
	return &u
}

// Float32 generates a random float32 value.
func Float32() float32 {
	return rand.Float32()
}

// Float32Ptr generates a random float32 value and returns a pointer to it.
func Float32Ptr() *float32 {
	f := Float32()
	return &f
}

// Float64 generates a random float64 value.
func Float64() float64 {
	return rand.Float64()
}

// Float64Ptr generates a random float64 value and returns a pointer to it.
func Float64Ptr() *float64 {
	f := Float64()
	return &f
}

// Complex64 generates a random complex64 value.
func Complex64() complex64 {
	r := rand.Float32()
	i := rand.Float32()
	return complex(r, i)
}

// Complex64Ptr generates a random complex64 value and returns a pointer to it.
func Complex64Ptr() *complex64 {
	c := Complex64()
	return &c
}

// Complex128 generates a random complex128 value.
func Complex128() complex128 {
	r := rand.Float64()
	i := rand.Float64()
	return complex(r, i)
}

// Complex128Ptr generates a random complex128 value and returns a pointer to it.
func Complex128Ptr() *complex128 {
	c := Complex128()
	return &c
}

// Error generates a random error value.
func Error() error {
	return errors.New(String())
}

// Duration generates a random time.Duration value.
func Duration() time.Duration {
	return time.Duration(rand.Int63())
}

// DurationPtr generates a random time.Duration value and returns a pointer to it.
func DurationPtr() *time.Duration {
	d := Duration()
	return &d
}
