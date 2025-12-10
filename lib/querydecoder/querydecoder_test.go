package querydecoder

import (
	"net/url"
	"testing"
	"time"
)

func TestDecode_SetsSupportedTypes(t *testing.T) {
	values := url.Values{
		"string":         {"abc"},
		"strings":        {"a", "b"},
		"int":            {"10"},
		"int8":           {"-8"},
		"int16":          {"-16"},
		"int32":          {"-32"},
		"int64":          {"-64"},
		"uint":           {"12"},
		"uint8":          {"8"},
		"uint16":         {"16"},
		"uint32":         {"32"},
		"uint64":         {"64"},
		"bool":           {"true"},
		"float32":        {"3.14"},
		"float64":        {"2.718"},
		"complex64":      {"(1+2i)"},
		"complex128":     {"(3+4i)"},
		"datetime":       {"2025-01-02T15:04:05Z"},
		"datetimeOffset": {"2025-01-02T15:04:05+09:00"},
		"date":           {"2025-01-02"},
		"time":           {"15:04:05"},
	}

	var dst struct {
		String         string     `query:"string"`
		Strings        []string   `query:"strings"`
		Int            int        `query:"int"`
		Int8           int8       `query:"int8"`
		Int16          int16      `query:"int16"`
		Int32          int32      `query:"int32"`
		Int64          int64      `query:"int64"`
		Uint           uint       `query:"uint"`
		Uint8          uint8      `query:"uint8"`
		Uint16         uint16     `query:"uint16"`
		Uint32         uint32     `query:"uint32"`
		Uint64         uint64     `query:"uint64"`
		Bool           bool       `query:"bool"`
		Float32        float32    `query:"float32"`
		Float64        float64    `query:"float64"`
		Complex64      complex64  `query:"complex64"`
		Complex128     complex128 `query:"complex128"`
		Datetime       time.Time  `query:"datetime"`
		DatetimeOffset time.Time  `query:"datetimeOffset"`
		Date           time.Time  `query:"date"`
		Time           time.Time  `query:"time"`
	}

	if err := Decode(values, &dst); err != nil {
		t.Fatalf("Decode returned error: %v", err)
	}

	if dst.String != "abc" {
		t.Fatalf("string mismatch: %+v", dst)
	}
	if len(dst.Strings) != 2 || dst.Strings[0] != "a" || dst.Strings[1] != "b" {
		t.Fatalf("strings mismatch: %+v", dst)
	}
	if dst.Int != 10 || dst.Int8 != -8 || dst.Int16 != -16 || dst.Int32 != -32 || dst.Int64 != -64 {
		t.Fatalf("int mismatch: %+v", dst)
	}
	if dst.Uint != 12 || dst.Uint8 != 8 || dst.Uint16 != 16 || dst.Uint32 != 32 || dst.Uint64 != 64 {
		t.Fatalf("uint mismatch: %+v", dst)
	}
	if !dst.Bool {
		t.Fatalf("bool mismatch: %+v", dst)
	}
	if dst.Float32 != 3.14 || dst.Float64 != 2.718 {
		t.Fatalf("float mismatch: %+v", dst)
	}
	if dst.Complex64 != complex64(1+2i) || dst.Complex128 != complex128(3+4i) {
		t.Fatalf("complex mismatch: %+v", dst)
	}
	if dst.Datetime.Format(time.RFC3339) != "2025-01-02T15:04:05Z" {
		t.Fatalf("datetime mismatch: %+v", dst)
	}
	if dst.Date.Format(time.DateOnly) != "2025-01-02" {
		t.Fatalf("date mismatch: %+v", dst)
	}
	if dst.Time.Format(time.TimeOnly) != "15:04:05" {
		t.Fatalf("time mismatch: %+v", dst)
	}
}

func TestDecode_InvalidDestination(t *testing.T) {
	if err := Decode(url.Values{}, nil); err == nil {
		t.Fatalf("expected error when dest is nil")
	}

	var notStruct int
	if err := Decode(url.Values{}, &notStruct); err == nil {
		t.Fatalf("expected error when dest is not struct pointer")
	}
}
