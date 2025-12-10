package querydecoder

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"time"
)

// Decode は url.Values を構造体へマッピングする。
// 構造体のフィールドに `query:"name"` タグを指定して利用する。
func Decode(values url.Values, dest interface{}) error {
	if values == nil {
		values = url.Values{}
	}
	if dest == nil {
		return fmt.Errorf("querydecoder: dest must not be nil")
	}
	v := reflect.ValueOf(dest)
	if v.Kind() != reflect.Pointer || v.IsNil() {
		return fmt.Errorf("querydecoder: dest must be non-nil pointer to struct")
	}
	structVal := v.Elem()
	if structVal.Kind() != reflect.Struct {
		return fmt.Errorf("querydecoder: dest must point to struct")
	}
	structType := structVal.Type()
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		tag := field.Tag.Get("query")
		if tag == "" || !field.IsExported() {
			continue
		}
		vals, ok := values[tag]
		if !ok || len(vals) == 0 {
			continue
		}
		if err := assignValue(structVal.Field(i), vals); err != nil {
			return fmt.Errorf("querydecoder: field %s: %w", field.Name, err)
		}
	}
	return nil
}

func assignValue(field reflect.Value, values []string) error {
	if !field.CanSet() {
		return fmt.Errorf("cannot set field")
	}
	last := values[len(values)-1]
	switch field.Kind() {

	case reflect.String:
		field.SetString(last)

	case reflect.Slice:
		if field.Type().Elem().Kind() != reflect.String {
			return fmt.Errorf("unsupported slice type %s", field.Type())
		}
		field.Set(reflect.ValueOf(values))

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		bitSize := field.Type().Bits()
		n, err := strconv.ParseInt(last, 10, bitSize)
		if err != nil {
			return err
		}
		field.SetInt(n)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		bitSize := field.Type().Bits()
		n, err := strconv.ParseUint(last, 10, bitSize)
		if err != nil {
			return err
		}
		field.SetUint(n)

	case reflect.Bool:
		b, err := strconv.ParseBool(last)
		if err != nil {
			return err
		}
		field.SetBool(b)

	case reflect.Float32, reflect.Float64:
		bitSize := field.Type().Bits()
		f, err := strconv.ParseFloat(last, bitSize)
		if err != nil {
			return err
		}
		field.SetFloat(f)

	case reflect.Complex64, reflect.Complex128:
		bitSize := field.Type().Bits()
		c, err := strconv.ParseComplex(last, bitSize)
		if err != nil {
			return err
		}
		field.SetComplex(c)
	case reflect.Struct:
		if field.Type() == reflect.TypeOf(time.Time{}) {
			t, err := parseTime(last)
			if err != nil {
				return err
			}
			field.Set(reflect.ValueOf(t))
			return nil
		}
		fallthrough
	default:
		return fmt.Errorf("unsupported kind %s", field.Kind())
	}
	return nil
}

func parseTime(value string) (time.Time, error) {
	// 優先的に RFC3339 を試し、それ以外に date/time のみも許容する。
	if t, err := time.Parse(time.RFC3339, value); err == nil {
		return t, nil
	}
	if t, err := time.Parse(time.DateOnly, value); err == nil {
		return t, nil
	}
	if t, err := time.Parse(time.TimeOnly, value); err == nil {
		return t, nil
	}
	return time.Time{}, fmt.Errorf("querydecoder: invalid time value %q", value)
}
