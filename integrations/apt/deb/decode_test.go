package deb_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/abemedia/appcast/integrations/apt/deb"
	"github.com/google/go-cmp/cmp"
)

func TestUnmarshal(t *testing.T) {
	in := `String: test
Hex: 01020304
Int: 1
Int8: 1
Int16: 1
Int32: 1
Uint: 1
Uint8: 1
Uint16: 1
Uint32: 1
Float32: 1.123
Float64: 1.123
Marshaler: test
Date: Tue, 10 Jan 2023 19:04:25 UTC
`

	want := record{
		String:    "test",
		ByteArray: [4]byte{1, 2, 3, 4},
		Int:       1,
		Int8:      1,
		Int16:     1,
		Int32:     1,
		Uint:      1,
		Uint8:     1,
		Uint16:    1,
		Uint32:    1,
		Float32:   1.123,
		Float64:   1.123,
		Marshaler: &marshaler{"test"},
		Date:      time.Date(2023, 1, 10, 19, 4, 25, 0, time.UTC),
	}

	tests := []struct {
		msg  string
		in   string
		want interface{}
	}{
		{
			msg:  "struct",
			in:   in,
			want: want,
		},
		{
			msg:  "struct pointer",
			in:   in,
			want: &want,
		},
		{
			msg:  "struct slice",
			in:   in + "\r\n\r\n" + in,
			want: []record{want, want},
		},
		{
			msg:  "struct pointer slice",
			in:   in + "\n" + in,
			want: []*record{&want, &want},
		},
		{
			msg: "multi-line string",
			in: `String: foo
 bar
 baz
 .
 foobar
`,
			want: record{
				String: "foo\nbar\nbaz\n\nfoobar",
			},
		},
		{
			msg: "multi-line string starting with empty line",
			in: `String:
 foo
 bar
`,
			want: record{
				String: "\nfoo\nbar",
			},
		},
	}

	for _, test := range tests {
		got := reflect.New(reflect.TypeOf(test.want)).Interface()
		err := deb.Unmarshal([]byte(test.in), got)
		if err != nil {
			t.Error(test.msg, err)
		}

		if diff := cmp.Diff(test.want, reflect.ValueOf(got).Elem().Interface()); diff != "" {
			t.Errorf("%s:\n%s", test.msg, diff)
		}
	}
}

func BenchmarkUnmarshal(b *testing.B) {
	type record struct {
		String    string
		ByteArray [4]byte `deb:"Hex"`
		Int       int
	}

	in := []byte(`String: test1
Hex: 01020304
Int: 1
`)

	var v record

	for i := 0; i < b.N; i++ {
		deb.Unmarshal(in, &v)
	}
}