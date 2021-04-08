package fields

import (
	"testing"

	"github.com/fluxynet/gocipe/types"
)

type node struct {
	name string
	kind types.Type
}

func checkValuesPresent(t *testing.T, want []node, got Fields) {
	var lw = len(want)
	var lg = got.Length()

	if lw != lg {
		t.Errorf("wantLength = %d\ngotLength = %d\n", lw, lg)
		return
	} else if lg == 0 {
		return
	}

	for i := range want {
		v := got.TypeOf(want[i].name)

		if v == "" {
			t.Errorf("not found: %s\n", want[i].name)
			return
		}

		if want[i].kind != v {
			t.Errorf("want: %s got: %s\n", want[i].name, v)
			return
		}
	}
}

func TestFromMap(t *testing.T) {
	type args struct {
		m map[string]types.Type
	}

	tests := []struct {
		name string
		args args
		want []node
	}{
		{
			name: "empty",
			args: args{
				m: map[string]types.Type{},
			},
			want: nil,
		},
		{
			name: "1 kind",
			args: args{
				m: map[string]types.Type{
					"name": types.String,
				},
			},
			want: []node{
				{name: "name", kind: types.String},
			},
		},
		{
			name: "2 values",
			args: args{
				m: map[string]types.Type{
					"name":      types.String,
					"is_active": types.Bool,
				},
			},
			want: []node{
				{name: "name", kind: types.String},
				{name: "is_active", kind: types.Bool},
			},
		},
		{
			name: "2 values swapped",
			args: args{
				m: map[string]types.Type{
					"is_active": types.Bool,
					"name":      types.String,
				},
			},
			want: []node{
				{name: "is_active", kind: types.Bool},
				{name: "name", kind: types.String},
			},
		},
		{
			name: "one of each kind",
			args: args{
				m: map[string]types.Type{
					"aBool":    types.Bool,
					"aString":  types.String,
					"aInteger": types.Int64,
					"aFloat":   types.Float64,
				},
			},
			want: []node{
				{
					name: "aBool",
					kind: types.Bool,
				},
				{
					name: "aString",
					kind: types.String,
				},
				{
					name: "aInteger",
					kind: types.Int64,
				},
				{
					name: "aFloat",
					kind: types.Float64,
				},
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

func TestFields_String(t *testing.T) {
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
				{name: "name", kind: types.String},
			},
			want: "name:string",
		},
		{
			name: "2 values",
			nodes: []node{
				{name: "name", kind: types.String},
				{name: "is_active", kind: types.Bool},
			},
			want: "name:string, is_active:bool",
		},
		{
			name: "2 values swapped",
			nodes: []node{
				{name: "is_active", kind: types.Bool},
				{name: "name", kind: types.String},
			},
			want: "is_active:bool, name:string",
		},
		{
			name: "one of each kind",
			nodes: []node{
				{name: "aBool", kind: types.Bool},
				{name: "aString", kind: types.String},
				{name: "aInt64", kind: types.Int64},
				{name: "aFloat64", kind: types.Float64},
			},
			want: "aBool:bool, aString:string, aInt64:int64, aFloat64:float64",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				nodes  []node
				fields Fields
			)

			for k := range tt.nodes {
				nodes = append(nodes, node{
					name: tt.nodes[k].name,
					kind: tt.nodes[k].kind,
				})

				fields.Set(tt.nodes[k].name, tt.nodes[k].kind)
			}

			var lw = len(nodes)
			var lg = fields.Length()

			if lw != lg {
				t.Errorf("wantLength = %d\ngotLength = %d\n", lw, lg)
				return
			} else if lg == 0 {
				return
			}

			raw := fields.String()
			if tt.want != raw {
				t.Errorf("w> %s\ng> %s\n", tt.want, raw)
				return
			}
		})
	}
}
