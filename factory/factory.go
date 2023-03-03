// Package factory is used for generating random values for testing purposes.
package factory

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const (
	minSliceLen = 2
	maxSliceLen = 4
	minMapLen   = 2
	maxMapLen   = 4

	stringLen = 64
	indexBits = 6                // 6 bits to represent a character index
	indexMax  = 63 / indexBits   // 10 character indices fits in 63 bits
	indexMask = 1<<indexBits - 1 // 111111 All 1-bits as many as indexBits
)

var (
	rnd   *rand.Rand
	chars = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

// init function will be only called once in runtime regardless of how many times the package is imported.
func init() {
	src := rand.NewSource(time.Now().UnixNano())
	rnd = rand.New(src)
}

// randInRange returns a non-negative pseudo-random number in the closed interval [min,max] from the default Source.
func randInRange(min, max int) int {
	return min + rnd.Intn(max-min+1)
}

func randPick(args ...string) string {
	i := rnd.Intn(len(args))
	return args[i]
}

// String generates a random string value.
func String() string {
	n := stringLen
	b := new(strings.Builder)
	b.Grow(n)

	// val is a randomly generated value.
	// count keeps track of how many more character indicies left in the current val.
	// l keeps track of the current length of the generated string.
	for val, count, ln := rnd.Int63(), indexMax, 0; ln < n; {
		if count == 0 {
			val, count = rnd.Int63(), indexMax
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

// StringPtr generates a random string value and returns the pointer to it.
func StringPtr() *string {
	s := String()
	return &s
}

// StringSlice generates a random string slice.
func StringSlice() []string {
	x := make([]string, randInRange(minSliceLen, maxSliceLen))
	for i := range x {
		x[i] = String()
	}

	return x
}

// Bool generates a random bool value.
func Bool() bool {
	return rnd.Intn(2) == 1
}

// BoolPtr generates a random bool value and returns the pointer to it.
func BoolPtr() *bool {
	b := Bool()
	return &b
}

// BoolSlice generates a random bool slice.
func BoolSlice() []bool {
	x := make([]bool, randInRange(minSliceLen, maxSliceLen))
	for i := range x {
		x[i] = Bool()
	}

	return x
}

// Int generates a random int value.
func Int() int {
	return int(rnd.Int63())
}

// IntPtr generates a random int value and returns the pointer to it.
func IntPtr() *int {
	i := Int()
	return &i
}

// IntSlice generates a random int slice.
func IntSlice() []int {
	x := make([]int, randInRange(minSliceLen, maxSliceLen))
	for i := range x {
		x[i] = Int()
	}

	return x
}

// Int8 generates a random int8 value.
func Int8() int8 {
	return int8(rnd.Int63())
}

// Int8Ptr generates a random int8 value and returns the pointer to it.
func Int8Ptr() *int8 {
	i := Int8()
	return &i
}

// Int8Slice generates a random int8 slice.
func Int8Slice() []int8 {
	x := make([]int8, randInRange(minSliceLen, maxSliceLen))
	for i := range x {
		x[i] = Int8()
	}

	return x
}

// Int16 generates a random int16 value.
func Int16() int16 {
	return int16(rnd.Int63())
}

// Int16Ptr generates a random int16 value and returns the pointer to it.
func Int16Ptr() *int16 {
	i := Int16()
	return &i
}

// Int16Slice generates a random int16 slice.
func Int16Slice() []int16 {
	x := make([]int16, randInRange(minSliceLen, maxSliceLen))
	for i := range x {
		x[i] = Int16()
	}

	return x
}

// Int32 generates a random int32 value.
func Int32() int32 {
	return int32(rnd.Int63())
}

// Int32Ptr generates a random int32 value and returns the pointer to it.
func Int32Ptr() *int32 {
	i := Int32()
	return &i
}

// Int32Slice generates a random int32 slice.
func Int32Slice() []int32 {
	x := make([]int32, randInRange(minSliceLen, maxSliceLen))
	for i := range x {
		x[i] = Int32()
	}

	return x
}

// Int64 generates a random int64 value.
func Int64() int64 {
	return int64(rnd.Int63())
}

// Int64Ptr generates a random int64 value and returns the pointer to it.
func Int64Ptr() *int64 {
	i := Int64()
	return &i
}

// Int64Slice generates a random int64 slice.
func Int64Slice() []int64 {
	x := make([]int64, randInRange(minSliceLen, maxSliceLen))
	for i := range x {
		x[i] = Int64()
	}

	return x
}

// Uint generates a random uint value.
func Uint() uint {
	return uint(rnd.Uint64())
}

// UintPtr generates a random uint value and returns the pointer to it.
func UintPtr() *uint {
	u := Uint()
	return &u
}

// UintSlice generates a random uint slice.
func UintSlice() []uint {
	x := make([]uint, randInRange(minSliceLen, maxSliceLen))
	for i := range x {
		x[i] = Uint()
	}

	return x
}

// Uint8 generates a random uint8 value.
func Uint8() uint8 {
	return uint8(rnd.Uint64())
}

// Uint8Ptr generates a random uint8 value and returns the pointer to it.
func Uint8Ptr() *uint8 {
	u := Uint8()
	return &u
}

// Uint8Slice generates a random uint8 slice.
func Uint8Slice() []uint8 {
	x := make([]uint8, randInRange(minSliceLen, maxSliceLen))
	for i := range x {
		x[i] = Uint8()
	}

	return x
}

// Uint16 generates a random uint16 value.
func Uint16() uint16 {
	return uint16(rnd.Uint64())
}

// Uint16Ptr generates a random uint16 value and returns the pointer to it.
func Uint16Ptr() *uint16 {
	u := Uint16()
	return &u
}

// Uint16Slice generates a random uint16 slice.
func Uint16Slice() []uint16 {
	x := make([]uint16, randInRange(minSliceLen, maxSliceLen))
	for i := range x {
		x[i] = Uint16()
	}

	return x
}

// Uint32 generates a random uint32 value.
func Uint32() uint32 {
	return rnd.Uint32()
}

// Uint32Ptr generates a random uint32 value and returns the pointer to it.
func Uint32Ptr() *uint32 {
	u := Uint32()
	return &u
}

// Uint32Slice generates a random uint32 slice.
func Uint32Slice() []uint32 {
	x := make([]uint32, randInRange(minSliceLen, maxSliceLen))
	for i := range x {
		x[i] = Uint32()
	}

	return x
}

// Uint64 generates a random uint64 value.
func Uint64() uint64 {
	return rnd.Uint64()
}

// Uint64Ptr generates a random uint64 value and returns the pointer to it.
func Uint64Ptr() *uint64 {
	u := Uint64()
	return &u
}

// Uint64Slice generates a random uint64 slice.
func Uint64Slice() []uint64 {
	x := make([]uint64, randInRange(minSliceLen, maxSliceLen))
	for i := range x {
		x[i] = Uint64()
	}

	return x
}

// Float32 generates a random float32 value.
func Float32() float32 {
	return rnd.Float32()
}

// Float32Ptr generates a random float32 value and returns the pointer to it.
func Float32Ptr() *float32 {
	f := Float32()
	return &f
}

// Float32Slice generates a random float32 slice.
func Float32Slice() []float32 {
	x := make([]float32, randInRange(minSliceLen, maxSliceLen))
	for i := range x {
		x[i] = Float32()
	}

	return x
}

// Float64 generates a random float64 value.
func Float64() float64 {
	return rnd.Float64()
}

// Float64Ptr generates a random float64 value and returns the pointer to it.
func Float64Ptr() *float64 {
	f := Float64()
	return &f
}

// Float64Slice generates a random float64 slice.
func Float64Slice() []float64 {
	x := make([]float64, randInRange(minSliceLen, maxSliceLen))
	for i := range x {
		x[i] = Float64()
	}

	return x
}

// Complex64 generates a random complex64 value.
func Complex64() complex64 {
	r := rnd.Float32()
	i := rnd.Float32()
	return complex(r, i)
}

// Complex64Ptr generates a random complex64 value and returns the pointer to it.
func Complex64Ptr() *complex64 {
	c := Complex64()
	return &c
}

// Complex64Slice generates a random complex64 slice.
func Complex64Slice() []complex64 {
	x := make([]complex64, randInRange(minSliceLen, maxSliceLen))
	for i := range x {
		x[i] = Complex64()
	}

	return x
}

// Complex128 generates a random complex128 value.
func Complex128() complex128 {
	r := rnd.Float64()
	i := rnd.Float64()
	return complex(r, i)
}

// Complex128Ptr generates a random complex128 value and returns the pointer to it.
func Complex128Ptr() *complex128 {
	c := Complex128()
	return &c
}

// Complex128Slice generates a random complex128 slice.
func Complex128Slice() []complex128 {
	x := make([]complex128, randInRange(minSliceLen, maxSliceLen))
	for i := range x {
		x[i] = Complex128()
	}

	return x
}

// Byte generates a random byte value.
func Byte() byte {
	return byte(rnd.Int63())
}

// BytePtr generates a random byte value and returns the pointer to it.
func BytePtr() *byte {
	b := Byte()
	return &b
}

// ByteSlice generates a random byte slice.
func ByteSlice() []byte {
	x := make([]byte, randInRange(minSliceLen, maxSliceLen))
	for i := range x {
		x[i] = Byte()
	}

	return x
}

// Rune generates a random rune value.
func Rune() rune {
	return rune(rnd.Int63())
}

// RunePtr generates a random rune value and returns the pointer to it.
func RunePtr() *rune {
	r := Rune()
	return &r
}

// RuneSlice generates a random rune slice.
func RuneSlice() []rune {
	x := make([]rune, randInRange(minSliceLen, maxSliceLen))
	for i := range x {
		x[i] = Rune()
	}

	return x
}

// Duration generates a random time.Duration value.
func Duration() time.Duration {
	h, m, s := rnd.Intn(100), rnd.Intn(60), rnd.Intn(60)
	d, _ := time.ParseDuration(fmt.Sprintf("%dh%dm%ds", h, m, s))
	return d
}

// DurationPtr generates a random time.Duration value and returns the pointer to it.
func DurationPtr() *time.Duration {
	d := Duration()
	return &d
}

// DurationSlice generates a random time.Duration slice.
func DurationSlice() []time.Duration {
	x := make([]time.Duration, randInRange(minSliceLen, maxSliceLen))
	for i := range x {
		x[i] = Duration()
	}

	return x
}
