// Package ptr is used for getting pointer to values.
// It eliminates the need for defining a new variable.
package ptr

import "time"

// String returns the pointer to a string value.
func String(s string) *string {
	return &s
}

// Bool returns the pointer to a bool value.
func Bool(b bool) *bool {
	return &b
}

// Int returns the pointer to an int value.
func Int(i int) *int {
	return &i
}

// Int8 returns the pointer to an int8 value.
func Int8(i int8) *int8 {
	return &i
}

// Int16 returns the pointer to an int16 value.
func Int16(i int16) *int16 {
	return &i
}

// Int32 returns the pointer to an int32 value.
func Int32(i int32) *int32 {
	return &i
}

// Int64 returns the pointer to an int64 value.
func Int64(i int64) *int64 {
	return &i
}

// Uint returns the pointer to an uint value.
func Uint(u uint) *uint {
	return &u
}

// Uint8 returns the pointer to an uint8 value.
func Uint8(u uint8) *uint8 {
	return &u
}

// Uint16 returns the pointer to an uint16 value.
func Uint16(u uint16) *uint16 {
	return &u
}

// Uint32 returns the pointer to an uint32 value.
func Uint32(u uint32) *uint32 {
	return &u
}

// Uint64 returns the pointer to an uint64 value.
func Uint64(u uint64) *uint64 {
	return &u
}

// Float32 returns the pointer to a float32 value.
func Float32(f float32) *float32 {
	return &f
}

// Float64 returns the pointer to a float64 value.
func Float64(f float64) *float64 {
	return &f
}

// Complex64 returns the pointer to a complex64 value.
func Complex64(c complex64) *complex64 {
	return &c
}

// Complex128 returns the pointer to a complex128 value.
func Complex128(c complex128) *complex128 {
	return &c
}

// Byte returns the pointer to a byte value.
func Byte(b byte) *byte {
	return &b
}

// Rune returns the pointer to a rune value.
func Rune(r rune) *rune {
	return &r
}

// Duration returns the pointer to a Duration value.
func Duration(d time.Duration) *time.Duration {
	return &d
}
