package factory

import (
	"errors"
	"fmt"
	"reflect"
)

// Populate receives the pointer to a struct and populates all fields, including the nested fields, with random values.
// If skipUnsupported set to false, an error will be returned in case a type is not supported.
func Populate(s interface{}, skipUnsupported bool) error {
	val, err := validateStruct(s)
	if err != nil {
		return err
	}

	return iterateOnFields(val, skipUnsupported, set)
}

func validateStruct(s interface{}) (reflect.Value, error) {
	if s == nil {
		return reflect.Value{}, errors.New("nil: you should pass a pointer to a struct type")
	}

	v := reflect.ValueOf(s) // reflect.Value --> v.Type(), v.Kind(), v.NumField()
	t := reflect.TypeOf(s)  // reflect.Type --> t.Kind(), t.Name(), t.NumField()

	// A pointer to a struct should be passed
	if t.Kind() != reflect.Ptr {
		return reflect.Value{}, errors.New("non-pointer type: you should pass a pointer to a struct type")
	}

	// Navigate to the pointer value
	v = v.Elem()
	t = t.Elem()

	if t.Kind() != reflect.Struct {
		return reflect.Value{}, errors.New("non-struct type: you should pass a pointer to a struct type")
	}

	return v, nil
}

func iterateOnFields(vStruct reflect.Value, skipUnsupported bool, handle func(reflect.Value) error) error {
	// Iterate over struct fields
	for i := 0; i < vStruct.NumField(); i++ {
		v := vStruct.Field(i)
		t := v.Type()

		// Recursively, iterate on nested structs
		// Nested structs do not need to have the `flag` tag and can be not settable.
		if isNestedStruct(t) {
			if err := iterateOnFields(v, skipUnsupported, handle); err != nil {
				return err
			}
			continue
		}

		// Skip unsupported fields
		if !isTypeSupported(t) {
			if skipUnsupported {
				continue
			}
			return fmt.Errorf("unsupported type: %s", t)
		}

		if err := handle(v); err != nil {
			return err
		}
	}

	return nil
}

func isTypeSupported(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.String:
		return true
	case reflect.Bool:
		return true
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	case reflect.Float32, reflect.Float64:
		return true
	case reflect.Complex64, reflect.Complex128:
		return true
	case reflect.Ptr, reflect.Array, reflect.Slice:
		return isTypeSupported(t.Elem())
	case reflect.Map:
		return isTypeSupported(t.Key()) && isTypeSupported(t.Elem())
	case reflect.Struct:
		return isStructSupported(t)
	case reflect.Interface:
		return isInterfaceSupported(t)
	default:
		return false
	}
}

func isStructSupported(t reflect.Type) bool {
	return (t.PkgPath() == "time" && t.Name() == "Time") ||
		(t.PkgPath() == "net/url" && t.Name() == "URL")
}

func isInterfaceSupported(t reflect.Type) bool {
	return false
}

func isNestedStruct(t reflect.Type) bool {
	if t.Kind() != reflect.Struct {
		return false
	}

	if isStructSupported(t) {
		return false
	}

	return true
}

func set(v reflect.Value) error {
	t := v.Type()

	switch t.Kind() {
	case reflect.Ptr:
		ptrV := reflect.New(t.Elem())
		x, err := generate(t.Elem())
		if err != nil {
			return err
		}
		ptrV.Elem().Set(x)
		v.Set(ptrV)

	case reflect.Array:
		for i := 0; i < v.Len(); i++ {
			x, err := generate(t.Elem())
			if err != nil {
				return err
			}
			v.Index(i).Set(x)
		}

	case reflect.Slice:
		l := randInRange(minSliceLen, maxSliceLen)
		sliceV := reflect.MakeSlice(t, 0, l)
		for i := 0; i < l; i++ {
			x, err := generate(t.Elem())
			if err != nil {
				return err
			}
			sliceV = reflect.Append(sliceV, x)
		}
		v.Set(sliceV)

	case reflect.Map:
		l := randInRange(minMapLen, maxMapLen)
		mapV := reflect.MakeMapWithSize(t, l)
		for i := 0; i < l; i++ {
			x, err := generate(t.Key())
			if err != nil {
				return err
			}
			y, err := generate(t.Elem())
			if err != nil {
				return err
			}
			mapV.SetMapIndex(x, y)
		}
		v.Set(mapV)

	default:
		x, err := generate(t)
		if err != nil {
			return err
		}
		v.Set(x)
	}

	return nil
}

func generate(t reflect.Type) (reflect.Value, error) {
	var x interface{}

	switch t.Kind() {
	case reflect.String:
		x = String()
	case reflect.Bool:
		x = Bool()
	case reflect.Int:
		x = Int()
	case reflect.Int8:
		x = Int8()
	case reflect.Int16:
		x = Int16()
	case reflect.Int32:
		x = Int32()
	case reflect.Int64:
		if t.PkgPath() == "time" && t.Name() == "Duration" {
			x = Duration()
		} else {
			x = Int64()
		}
	case reflect.Uint:
		x = Uint()
	case reflect.Uint8:
		x = Uint8()
	case reflect.Uint16:
		x = Uint16()
	case reflect.Uint32:
		x = Uint32()
	case reflect.Uint64:
		x = Uint64()
	case reflect.Float32:
		x = Float32()
	case reflect.Float64:
		x = Float64()
	case reflect.Complex64:
		x = Complex64()
	case reflect.Complex128:
		x = Complex128()
	case reflect.Struct:
		if t.PkgPath() == "time" && t.Name() == "Time" {
			x = Time()
		} else if t.PkgPath() == "net/url" && t.Name() == "URL" {
			x = URL()
		}
	}

	if x == nil {
		return reflect.Zero(t), fmt.Errorf("unsupported type: %s", t)
	}

	return reflect.ValueOf(x), nil
}
