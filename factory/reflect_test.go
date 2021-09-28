package factory

import (
	"errors"
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type Struct struct {
	String     string
	Bool       bool
	Int        int
	Int8       int8
	Int16      int16
	Int32      int32
	Int64      int64
	Uint       uint
	Uint8      uint8
	Uint16     uint16
	Uint32     uint32
	Uint64     uint64
	Float32    float32
	Float64    float64
	Complex64  complex64
	Complex128 complex128
	Byte       byte
	Rune       rune
	Duration   time.Duration
	Time       time.Time
	URL        url.URL

	Pointer struct {
		String     *string
		Bool       *bool
		Int        *int
		Int8       *int8
		Int16      *int16
		Int32      *int32
		Int64      *int64
		Uint       *uint
		Uint8      *uint8
		Uint16     *uint16
		Uint32     *uint32
		Uint64     *uint64
		Float32    *float32
		Float64    *float64
		Complex64  *complex64
		Complex128 *complex128
		Byte       *byte
		Rune       *rune
		Duration   *time.Duration
		Time       *time.Time
		URL        *url.URL
	}

	Array struct {
		String     [2]string
		Bool       [2]bool
		Int        [2]int
		Int8       [2]int8
		Int16      [2]int16
		Int32      [2]int32
		Int64      [2]int64
		Uint       [2]uint
		Uint8      [2]uint8
		Uint16     [2]uint16
		Uint32     [2]uint32
		Uint64     [2]uint64
		Float32    [2]float32
		Float64    [2]float64
		Complex64  [2]complex64
		Complex128 [2]complex128
		Byte       [2]byte
		Rune       [2]rune
		Duration   [2]time.Duration
		Time       [2]time.Time
		URL        [2]url.URL
	}

	Slice struct {
		String     []string
		Bool       []bool
		Int        []int
		Int8       []int8
		Int16      []int16
		Int32      []int32
		Int64      []int64
		Uint       []uint
		Uint8      []uint8
		Uint16     []uint16
		Uint32     []uint32
		Uint64     []uint64
		Float32    []float32
		Float64    []float64
		Complex64  []complex64
		Complex128 []complex128
		Byte       []byte
		Rune       []rune
		Duration   []time.Duration
		Time       []time.Time
		URL        []url.URL
	}

	Map struct {
		String     map[string]string
		Bool       map[bool]bool
		Int        map[int]int
		Int8       map[int8]int8
		Int16      map[int16]int16
		Int32      map[int32]int32
		Int64      map[int64]int64
		Uint       map[uint]uint
		Uint8      map[uint8]uint8
		Uint16     map[uint16]uint16
		Uint32     map[uint32]uint32
		Uint64     map[uint64]uint64
		Float32    map[float32]float32
		Float64    map[float64]float64
		Complex64  map[complex64]complex64
		Complex128 map[complex128]complex128
		Byte       map[byte]byte
		Rune       map[rune]rune
		Duration   map[time.Duration]time.Duration
		Time       map[time.Time]time.Time
		URL        map[url.URL]url.URL
	}
}

func TestPopulate(t *testing.T) {
	tests := []struct {
		name            string
		s               interface{}
		skipUnsupported bool
		expectedError   error
	}{
		{
			"Nil",
			nil,
			false,
			errors.New("nil: you should pass a pointer to a struct type"),
		},
		{
			"NonStruct",
			new(string),
			false,
			errors.New("non-struct type: you should pass a pointer to a struct type"),
		},
		{
			"NonPointer",
			struct{}{},
			false,
			errors.New("non-pointer type: you should pass a pointer to a struct type"),
		},
		{
			"UnsupportedType",
			&struct{ Error error }{},
			false,
			errors.New("unsupported type: error"),
		},
		{
			"SkipUnsupportedType",
			&struct{ Error error }{},
			true,
			nil,
		},
		{
			"Supported",
			&Struct{},
			false,
			nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := Populate(tc.s, tc.skipUnsupported)

			if tc.expectedError != nil {
				assert.Equal(t, tc.expectedError, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateStruct(t *testing.T) {
	tests := []struct {
		name          string
		v             interface{}
		expectedError error
	}{
		{
			"Nil",
			nil,
			errors.New("nil: you should pass a pointer to a struct type"),
		},
		{
			"NonStruct",
			new(string),
			errors.New("non-struct type: you should pass a pointer to a struct type"),
		},
		{
			"NonPointer",
			struct{}{},
			errors.New("non-pointer type: you should pass a pointer to a struct type"),
		},
		{
			"OK",
			new(struct{}),
			nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v, err := validateStruct(tc.v)

			if tc.expectedError == nil {
				assert.NotNil(t, v)
				assert.NoError(t, err)
			} else {
				assert.Empty(t, v)
				assert.Equal(t, tc.expectedError, err)
			}
		})
	}
}

func TestIterateOnFields(t *testing.T) {
	tests := []struct {
		name               string
		v                  interface{}
		skipUnsupported    bool
		expectedError      error
		expectedFieldCount int
	}{
		{
			"UnsupportedType",
			struct{ Error error }{},
			false,
			errors.New("unsupported type: error"),
			0,
		},
		{
			"SkipUnsupportedType",
			struct{ Error error }{},
			true,
			nil,
			0,
		},
		{
			"OK",
			Struct{},
			false,
			nil,
			105,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fieldCount := 0

			vStruct := reflect.ValueOf(tc.v)
			err := iterateOnFields(vStruct, tc.skipUnsupported, func(v reflect.Value) error {
				fieldCount++
				return nil
			})

			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedFieldCount, fieldCount)
		})
	}
}

func TestIsTypeSupported(t *testing.T) {
	var s Struct

	tests := []struct {
		name              string
		v                 interface{}
		expectedSupported bool
	}{
		{"String", s.String, true},
		{"Bool", s.Bool, true},
		{"Int", s.Int, true},
		{"Int8", s.Int8, true},
		{"Int16", s.Int16, true},
		{"Int32", s.Int32, true},
		{"Int64", s.Int64, true},
		{"Uint", s.Uint, true},
		{"Uint8", s.Uint8, true},
		{"Uint16", s.Uint16, true},
		{"Uint32", s.Uint32, true},
		{"Uint64", s.Uint64, true},
		{"Float32", s.Float32, true},
		{"Float64", s.Float64, true},
		{"Complex64", s.Complex64, true},
		{"Complex128", s.Complex128, true},
		{"Byte", s.Byte, true},
		{"Rune", s.Rune, true},
		{"Duration", s.Duration, true},
		{"Time", s.Time, true},
		{"URL", s.URL, true},
		{"StringPointer", s.Pointer.String, true},
		{"BoolPointer", s.Pointer.Bool, true},
		{"IntPointer", s.Pointer.Int, true},
		{"Int8Pointer", s.Pointer.Int8, true},
		{"Int16Pointer", s.Pointer.Int16, true},
		{"Int32Pointer", s.Pointer.Int32, true},
		{"Int64Pointer", s.Pointer.Int64, true},
		{"UintPointer", s.Pointer.Uint, true},
		{"Uint8Pointer", s.Pointer.Uint8, true},
		{"Uint16Pointer", s.Pointer.Uint16, true},
		{"Uint32Pointer", s.Pointer.Uint32, true},
		{"Uint64Pointer", s.Pointer.Uint64, true},
		{"Float32Pointer", s.Pointer.Float32, true},
		{"Float64Pointer", s.Pointer.Float64, true},
		{"Complex64Pointer", s.Pointer.Complex64, true},
		{"Complex128Pointer", s.Pointer.Complex128, true},
		{"BytePointer", s.Pointer.Byte, true},
		{"RunePointer", s.Pointer.Rune, true},
		{"DurationPointer", s.Pointer.Duration, true},
		{"TimePointer", s.Pointer.Time, true},
		{"URLPointer", s.Pointer.URL, true},
		{"StringArray", s.Array.String, true},
		{"BoolArray", s.Array.Bool, true},
		{"IntArray", s.Array.Int, true},
		{"Int8Array", s.Array.Int8, true},
		{"Int16Array", s.Array.Int16, true},
		{"Int32Array", s.Array.Int32, true},
		{"Int64Array", s.Array.Int64, true},
		{"UintArray", s.Array.Uint, true},
		{"Uint8Array", s.Array.Uint8, true},
		{"Uint16Array", s.Array.Uint16, true},
		{"Uint32Array", s.Array.Uint32, true},
		{"Uint64Array", s.Array.Uint64, true},
		{"Float32Array", s.Array.Float32, true},
		{"Float64Array", s.Array.Float64, true},
		{"Complex64Array", s.Array.Complex64, true},
		{"Complex128Array", s.Array.Complex128, true},
		{"ByteArray", s.Array.Byte, true},
		{"RuneArray", s.Array.Rune, true},
		{"DurationArray", s.Array.Duration, true},
		{"TimeArray", s.Array.Time, true},
		{"URLArray", s.Array.URL, true},
		{"StringSlice", s.Slice.String, true},
		{"BoolSlice", s.Slice.Bool, true},
		{"IntSlice", s.Slice.Int, true},
		{"Int8Slice", s.Slice.Int8, true},
		{"Int16Slice", s.Slice.Int16, true},
		{"Int32Slice", s.Slice.Int32, true},
		{"Int64Slice", s.Slice.Int64, true},
		{"UintSlice", s.Slice.Uint, true},
		{"Uint8Slice", s.Slice.Uint8, true},
		{"Uint16Slice", s.Slice.Uint16, true},
		{"Uint32Slice", s.Slice.Uint32, true},
		{"Uint64Slice", s.Slice.Uint64, true},
		{"Float32Slice", s.Slice.Float32, true},
		{"Float64Slice", s.Slice.Float64, true},
		{"Complex64Slice", s.Slice.Complex64, true},
		{"Complex128Slice", s.Slice.Complex128, true},
		{"ByteSlice", s.Slice.Byte, true},
		{"RuneSlice", s.Slice.Rune, true},
		{"DurationSlice", s.Slice.Duration, true},
		{"TimeSlice", s.Slice.Time, true},
		{"URLSlice", s.Slice.URL, true},
		{"StringMap", s.Map.String, true},
		{"BoolMap", s.Map.Bool, true},
		{"IntMap", s.Map.Int, true},
		{"Int8Map", s.Map.Int8, true},
		{"Int16Map", s.Map.Int16, true},
		{"Int32Map", s.Map.Int32, true},
		{"Int64Map", s.Map.Int64, true},
		{"UintMap", s.Map.Uint, true},
		{"Uint8Map", s.Map.Uint8, true},
		{"Uint16Map", s.Map.Uint16, true},
		{"Uint32Map", s.Map.Uint32, true},
		{"Uint64Map", s.Map.Uint64, true},
		{"Float32Map", s.Map.Float32, true},
		{"Float64Map", s.Map.Float64, true},
		{"Complex64Map", s.Map.Complex64, true},
		{"Complex128Map", s.Map.Complex128, true},
		{"ByteMap", s.Map.Byte, true},
		{"RuneMap", s.Map.Rune, true},
		{"DurationMap", s.Map.Duration, true},
		{"TimeMap", s.Map.Time, true},
		{"URLMap", s.Map.URL, true},
		{"NotSupported", make(chan struct{}), false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			supported := isTypeSupported(reflect.TypeOf(tc.v))

			assert.Equal(t, tc.expectedSupported, supported)
		})
	}
}

func TestIsStructSupported(t *testing.T) {
	tests := []struct {
		name     string
		v        interface{}
		expected bool
	}{
		{"NotSupported", struct{}{}, false},
		{"Time", time.Time{}, true},
		{"URL", url.URL{}, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tStruct := reflect.TypeOf(tc.v)

			assert.Equal(t, tc.expected, isStructSupported(tStruct))
		})
	}
}

func TestIsInterfaceSupported(t *testing.T) {
	tests := []struct {
		name     string
		v        interface{}
		expected bool
	}{
		{"NotSupported", errors.New(""), false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tInterface := reflect.TypeOf(tc.v)

			assert.Equal(t, tc.expected, isInterfaceSupported(tInterface))
		})
	}
}

func TestIsNestedStruct(t *testing.T) {
	vStruct := reflect.ValueOf(struct {
		Int    int
		Time   time.Time
		URL    url.URL
		Nested struct {
			String string
		}
	}{})

	vInt := vStruct.FieldByName("Int")
	assert.False(t, isNestedStruct(vInt.Type()))

	vTime := vStruct.FieldByName("Time")
	assert.False(t, isNestedStruct(vTime.Type()))

	vURL := vStruct.FieldByName("URL")
	assert.False(t, isNestedStruct(vURL.Type()))

	vNested := vStruct.FieldByName("Nested")
	assert.True(t, isNestedStruct(vNested.Type()))
}

func TestSet(t *testing.T) {
	t.Run("Supported", func(t *testing.T) {
		s := struct {
			String      string
			StringPtr   *string
			StringArray [2]string
			StringSlice []string
			StringMap   map[string]string
		}{}

		vStruct := reflect.ValueOf(&s).Elem()
		for i := 0; i < vStruct.NumField(); i++ {
			v := vStruct.Field(i)
			err := set(v)

			assert.NoError(t, err)
			assert.NotZero(t, v)
		}
	})

	t.Run("Unsupported", func(t *testing.T) {
		s := struct {
			Error      error
			ErrorPtr   *error
			ErrorArray [2]error
			ErrorSlice []error
			ErrorMap1  map[error]error
			ErrorMap2  map[string]error
		}{}

		vStruct := reflect.ValueOf(&s).Elem()
		for i := 0; i < vStruct.NumField(); i++ {
			v := vStruct.Field(i)
			err := set(v)

			assert.Error(t, err)
			assert.EqualError(t, err, "unsupported type: error")
		}
	})
}

func TestGenerate(t *testing.T) {
	tests := []struct {
		name          string
		v             interface{}
		expectedError error
	}{
		{"String", "", nil},
		{"Bool", false, nil},
		{"Int", int(0), nil},
		{"Int8", int8(0), nil},
		{"Int16", int16(0), nil},
		{"Int32", int32(0), nil},
		{"Int64", int64(0), nil},
		{"Uint", uint(0), nil},
		{"Uint8", uint8(0), nil},
		{"Uint16", uint16(0), nil},
		{"Uint32", uint32(0), nil},
		{"Uint64", uint64(0), nil},
		{"Float32", float32(0), nil},
		{"Float64", float64(0), nil},
		{"Complex64", complex64(0), nil},
		{"Complex128", complex128(0), nil},
		{"Byte", byte(0), nil},
		{"Rune", rune(0), nil},
		{"Duration", time.Duration(0), nil},
		{"Time", time.Time{}, nil},
		{"URL", url.URL{}, nil},
		{"UnsupportedStruct", time.Location{}, errors.New("unsupported type: time.Location")},
		{"UnsupportedType", errors.New("error"), errors.New("unsupported type: *errors.errorString")},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			v, err := generate(reflect.TypeOf(tc.v))

			if tc.expectedError != nil {
				assert.Equal(t, tc.expectedError, err)
			} else {
				assert.NoError(t, err)
				assert.NotZero(t, v)
				assert.NotNil(t, v.Interface())
			}
		})
	}
}
