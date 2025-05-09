package moments

import (
	"encoding/json"
	"path"
	"reflect"
	"strings"

	"github.com/google/uuid"
)

func DefaultIfNil[T any](value *T, defaultValue T) T {
	if value != nil {
		return *value
	}
	return defaultValue
}

func DefaultIfEmpty[T any](value *T, defaultValue T) T {
	if !IsEmpty(value) {
		return *value
	}
	return defaultValue
}
func DeepCopyJson[T any](src T) T {
	var dest T
	bytes, err := json.Marshal(src)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(bytes, &dest)
	if err != nil {
		panic(err)
	}
	return dest
}

func FilterSlice[T any](input []T, fn func(a T) bool) []T {
	output := make([]T, 0)
	for _, item := range input {
		if fn(item) {
			output = append(output, item)
		}
	}
	return output
}

// AnySlice converts a slice of any type to a slice of any type
func AnySlice[T any](slice []T) []any {
	result := make([]any, len(slice))
	for i, v := range slice {
		result[i] = v
	}
	return result
}

func NewRandomUuidString() string {
	return uuid.NewString()
}

// NewSequentialUUIDString returns a new UUID string that is formatted as a UUID v7 string.
// UUID v7 is a time ordered UUID
func NewSequentialUUIDString() string {
	id, err := uuid.NewV7()
	if err != nil {
		panic(err)
	}
	return id.String()
}

// NewSquentialString returns a new UUID string that is formatted as a UUID v7 string.
// UUID v7 is a time ordered UUID
func NewSquentialString() string {
	return strings.ReplaceAll(NewSequentialUUIDString(), "-", "")
}

func MapSlice[TIn any, TOut any](input []TIn, mapFn func(a TIn) TOut) []TOut {
	output := make([]TOut, len(input))
	for i, item := range input {
		output[i] = mapFn(item)
	}
	return output
}

func IsEmpty(object interface{}) bool {

	// get nil case out of the way
	if object == nil {
		return true
	}

	objValue := reflect.ValueOf(object)

	switch objValue.Kind() {
	// collection types are empty when they have no element
	case reflect.Chan, reflect.Map, reflect.Slice:
		return objValue.Len() == 0
	// pointers are empty if nil or if the value they point to is empty
	case reflect.Ptr:
		if objValue.IsNil() {
			return true
		}
		deref := objValue.Elem().Interface()
		return IsEmpty(deref)
	// for all other types, compare against the zero value
	// array types are empty when they match their zero-initialized state
	default:
		zero := reflect.Zero(objValue.Type())
		return reflect.DeepEqual(object, zero.Interface())
	}
}

func TypeName(value any) string {
	ty := reflect.TypeOf(value)
	pkg := path.Base(ty.PkgPath())
	if pkg == "." {
		return ty.Name()
	}
	return pkg + "." + ty.Name()
}
