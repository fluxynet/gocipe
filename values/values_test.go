package values

import (
	"bytes"
	"github.com/fluxynet/gocipe"
	"github.com/fluxynet/gocipe/fields"
	"io"
	"reflect"
	"testing"
)

type node struct {
	key   string
	value interface{}
}

func checkValuesPresent(t *testing.T, want []node, got *Values) {
	var lw = len(want)
	var lg int
	if got != nil {
		lg = got.Length()
	}

	if lw != lg {
		t.Errorf("wantLength = %d\ngotLength = %d\n", lw, lg)
		return
	} else if lg == 0 {
		return
	}

	for i := range want {
		v := got.Get(want[i].key)

		if v == nil {
			t.Errorf("not found: %s\n", want[i].key)
			return
		}

		if !reflect.DeepEqual(want[i].value, v.Value) {
			t.Errorf("want: [%T] %v\ngot: [%T] %v\n", want[i].value, want[i].value, v.Value, v.Value)
			return
		}
	}
}

func compareValues(t *testing.T, want, got *Values) {
	if (want == nil) != (got == nil) {
		t.Errorf("nils: want = %t got = %t", want == nil, got == nil)
		return
	}

	var lw, lg int

	if want != nil {
		lw = want.Length()
	}

	if got != nil {
		lg = got.Length()
	}

	if lw != lg {
		t.Errorf("wantLength = %d\ngotLength = %d\n", lw, lg)
		return
	} else if lg == 0 {
		return
	}

	var wit = want.Iterator()
	for wit.Next() {
		w := wit.Value()
		g := got.Get(w.Name)

		if w.Name != g.Name {
			t.Errorf("want: %s\ngot: %s\n", w.Name, g.Name)
			return
		}

		if !reflect.DeepEqual(w.Value, g.Value) {
			t.Errorf("want: [%T] %v\ngot: [%T] %v\n", w.Value, w.Value, g.Value, g.Value)
			return
		}
	}
}

func compareValuesToNodes(t *testing.T, want []node, got Values) {
	var lw = len(want)
	var lg = got.Length()

	if lw != lg {
		t.Errorf("wantLength = %d\ngotLength = %d\n", lw, lg)
		return
	} else if lg == 0 {
		return
	}

	var it = got.Iterator()
	for i := 0; it.Next(); i++ {
		v := it.Value()

		if want[i].key != v.Name {
			t.Errorf("(%d)\nwant: %s\ngot: %s\n", i, want[i].key, v.Name)
			return
		}

		if reflect.DeepEqual(want[i].value, v.Value) {
			t.Errorf("(%d)\nwant: [%T] %v\ngot: [%T] %v\n", i, want[i].value, want[i].value, v.Value, v.Value)
			return
		}
	}
}

func TestFromMap(t *testing.T) {
	type args struct {
		m map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want []node
	}{
		{
			name: "empty",
			args: args{
				m: nil,
			},
			want: nil,
		},
		{
			name: "1 value",
			args: args{
				m: map[string]interface{}{
					"name": "Bond",
				},
			},
			want: []node{
				{
					key:   "name",
					value: "Bond",
				},
			},
		},
		{
			name: "2 values",
			args: args{
				m: map[string]interface{}{
					"name":  "Bond",
					"agent": 007,
				},
			},
			want: []node{
				{
					key:   "name",
					value: "Bond",
				},
				{
					key:   "agent",
					value: 007,
				},
			},
		},
		{
			name: "2 values swapped",
			args: args{
				m: map[string]interface{}{
					"agent": 007,
					"name":  "Bond",
				},
			},
			want: []node{
				{
					key:   "agent",
					value: 007,
				},
				{
					key:   "name",
					value: "Bond",
				},
			},
		},
		{
			name: "one of each kind",
			args: args{
				m: map[string]interface{}{
					"aBool":    bool(true),
					"aString":  string("foo"),
					"aInt":     int(-53),
					"aInt8":    int8(-120),
					"aInt16":   int16(-32760),
					"aInt32":   int32(-2147483640),
					"aInt64":   int64(-922337203685477580),
					"aUint":    uint(53),
					"aUint8":   uint8(120),
					"aUint16":  uint16(32760),
					"aUint32":  uint32(2147483640),
					"aUint64":  uint64(922337203685477580),
					"aByte":    byte('a'),
					"aRune":    rune('👍'),
					"aFloat32": float32(3.14),
					"aFloat64": float64(3.141592653589793238),
				},
			},
			want: []node{
				{key: "aBool", value: bool(true)},
				{key: "aInt8", value: int8(-120)},
				{key: "aUint8", value: uint8(120)},
				{key: "aInt16", value: int16(-32760)},
				{key: "aString", value: string("foo")},
				{key: "aInt32", value: int32(-2147483640)},
				{key: "aInt", value: int(-53)},
				{key: "aInt64", value: int64(-922337203685477580)},
				{key: "aUint", value: uint(53)},
				{key: "aUint16", value: uint16(32760)},
				{key: "aFloat64", value: float64(3.141592653589793238)},
				{key: "aUint64", value: uint64(922337203685477580)},
				{key: "aByte", value: byte('a')},
				{key: "aRune", value: rune('👍')},
				{key: "aUint32", value: uint32(2147483640)},
				{key: "aFloat32", value: float32(3.14)},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FromMap(tt.args.m)
			checkValuesPresent(t, tt.want, got)
		})
	}
}

func TestFromJSON(t *testing.T) {
	type args struct {
		r string
		f fields.Fields
	}
	tests := []struct {
		name    string
		args    args
		want    *Values
		wantErr bool
	}{
		{
			name: "empty json",
			args: args{
				r: `{}`,
				f: fields.FromMap(map[string]gocipe.Type{}),
			},
			want:    FromMap(map[string]interface{}{}),
			wantErr: false,
		},
		{
			name: "empty fields",
			args: args{
				r: `{"name": "foo"}`,
				f: fields.FromMap(map[string]gocipe.Type{}),
			},
			want:    FromMap(map[string]interface{}{}),
			wantErr: false,
		},
		{
			name: "empty",
			args: args{
				r: `{}`,
				f: fields.FromMap(map[string]gocipe.Type{}),
			},
			want:    FromMap(map[string]interface{}{}),
			wantErr: false,
		},
		{
			name: "invalid json",
			args: args{
				r: `{"name": "foo", }`,
				f: fields.FromMap(map[string]gocipe.Type{
					"name": gocipe.String,
				}),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "1 field",
			args: args{
				r: `{"name": "foo"}`,
				f: fields.FromMap(map[string]gocipe.Type{
					"name": gocipe.String,
				}),
			},
			want: FromMap(map[string]interface{}{
				"name": "foo",
			}),
			wantErr: false,
		},
		{
			name: "2 fields",
			args: args{
				r: `{"country": "MU", "island": true}`,
				f: fields.FromMap(map[string]gocipe.Type{
					"country": gocipe.String,
					"island":  gocipe.Bool,
				}),
			},
			want: FromMap(map[string]interface{}{
				"country": "MU",
				"island":  true,
			}),
			wantErr: false,
		},
		{
			name: "one of each kind",
			args: args{
				r: `
{
    "aBool": true,
    "aInteger": 120,
    "aString": "foo",
    "aFloat": 3.141592653589793238
}
				`,
				f: fields.FromMap(map[string]gocipe.Type{
					"aBool":    gocipe.Bool,
					"aInteger": gocipe.Int64,
					"aString":  gocipe.String,
					"aFloat":   gocipe.Float64,
				}),
			},
			want: FromMap(map[string]interface{}{
				"aBool":    true,
				"aInteger": 120,
				"aString":  "foo",
				"aFloat":   float64(3.141592653589793238),
			}),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FromJSON(io.NopCloser(bytes.NewReader([]byte(tt.args.r))), tt.args.f, false)

			if (err != nil) && !tt.wantErr {
				t.Errorf("unwanted error: %v\n", tt.wantErr)
				return
			} else if (err == nil) && tt.wantErr {
				t.Errorf("error did not occur as wanted\n")
				return
			}

			compareValues(t, got, tt.want)
		})
	}
}

func Test_iterator_Next(t *testing.T) {
	tests := []struct {
		name  string
		nodes []node
		want  string
	}{
		{
			name: "empty",
		},
		{
			name: "1 kind",
			nodes: []node{
				{key: "name", value: gocipe.String},
			},
			want: `name:"string"`,
		},
		{
			name: "2 values",
			nodes: []node{
				{key: "name", value: gocipe.String},
				{key: "is_active", value: gocipe.Bool},
			},
			want: `name:"string", is_active:"bool"`,
		},
		{
			name: "2 values swapped",
			nodes: []node{
				{key: "is_active", value: gocipe.Bool},
				{key: "name", value: gocipe.String},
			},
			want: `is_active:"bool", name:"string"`,
		},
		{
			name: "one of each kind",
			nodes: []node{
				{key: "aBool", value: gocipe.Bool},
				{key: "aString", value: gocipe.String},
				{key: "aInt64", value: gocipe.Int64},
				{key: "aFloat64", value: gocipe.Float64},
			},
			want: `aBool:"bool", aString:"string", aInt64:"int64", aFloat64:"float64"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				nodes  []node
				values Values
			)

			for k := range tt.nodes {
				nodes = append(nodes, node{
					key:   tt.nodes[k].key,
					value: tt.nodes[k].value,
				})

				values.Set(tt.nodes[k].key, tt.nodes[k].value)
			}

			var lw = len(nodes)
			var lg = values.Length()

			if lw != lg {
				t.Errorf("wantLength = %d\ngotLength = %d\n", lw, lg)
				return
			} else if lg == 0 {
				return
			}

			raw := values.String()
			if tt.want != raw {
				t.Errorf("w> %s\ng> %s\n", tt.want, raw)
				return
			}
		})
	}
}
