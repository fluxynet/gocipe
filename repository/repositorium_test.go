package repository

import (
	"reflect"
	"testing"

	"github.com/fluxynet/gocipe"
	"github.com/fluxynet/gocipe/fields"
)

func compareOrderBys(t *testing.T, got, want []OrderBy) {
	var (
		lg = len(got)
		lw = len(want)
	)

	if lg != lw {
		t.Errorf("want len = %d\ngot len = %d\n", lw, lg)
		return
	}

	if lg == 0 {
		return
	}

	for i := range got {
		if got[i].Attribute != want[i].Attribute {
			t.Errorf("want = %s\ngot = %s\n", want[i].Attribute, got[i].Attribute)
			return
		}

		if got[i].Sort != want[i].Sort {
			t.Errorf("want = %s\ngot = %s\n", want[i].Sort, got[i].Sort)
			return
		}
	}
}

func compareConditions(t *testing.T, got, want []Condition) {
	var (
		lg = len(got)
		lw = len(want)
	)

	if lg != lw {
		t.Errorf("want len = %d\ngot len = %d\n", lw, lg)
		return
	}

	if lg == 0 {
		return
	}

	for i := range got {
		if got[i].Attribute != want[i].Attribute {
			t.Errorf("want = %s\ngot = %s\n", want[i].Attribute, got[i].Attribute)
			return
		}

		if got[i].Operator != want[i].Operator {
			t.Errorf("want = %s\ngot = %s\n", want[i].Operator, got[i].Operator)
			return
		}

		if !reflect.DeepEqual(got[i].Value, want[i].Value) {
			t.Errorf("want = %v\ngot  = %v\n", want[i].Value, got[i].Value)
			return
		}
	}
}

func TestParseOrderBy(t *testing.T) {
	type args struct {
		s string
		f fields.Fields
	}
	tests := []struct {
		name    string
		args    args
		want    []OrderBy
		wantErr bool
	}{
		{
			name: "Empty",
			args: args{
				s: "",
				f: fields.Fields{},
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "In string not in fields",
			args: args{
				s: "name,age",
				f: fields.FromPairs([]fields.Field{
					{Name: "country", Kind: gocipe.String},
				}),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "In fields not in string",
			args: args{
				s: "",
				f: fields.FromPairs([]fields.Field{
					{Name: "town", Kind: gocipe.String},
				}),
			},
			want:    []OrderBy{},
			wantErr: false,
		},
		{
			name: "1 single asc",
			args: args{
				s: "country",
				f: fields.FromPairs([]fields.Field{
					{Name: "country", Kind: gocipe.String},
					{Name: "town", Kind: gocipe.String},
					{Name: "name", Kind: gocipe.String},
				}),
			},
			want: []OrderBy{
				{Attribute: "country", Sort: Ascending},
			},
			wantErr: false,
		},
		{
			name: "2 ascending",
			args: args{
				s: "country,name",
				f: fields.FromPairs([]fields.Field{
					{Name: "country", Kind: gocipe.String},
					{Name: "town", Kind: gocipe.String},
					{Name: "name", Kind: gocipe.String},
				}),
			},
			want: []OrderBy{
				{Attribute: "country", Sort: Ascending},
				{Attribute: "name", Sort: Ascending},
			},
			wantErr: false,
		},
		{
			name: "1 single desc",
			args: args{
				s: "-country",
				f: fields.FromPairs([]fields.Field{
					{Name: "country", Kind: gocipe.String},
					{Name: "town", Kind: gocipe.String},
					{Name: "name", Kind: gocipe.String},
				}),
			},
			want: []OrderBy{
				{Attribute: "country", Sort: Descending},
			},
			wantErr: false,
		},
		{
			name: "2 descending",
			args: args{
				s: "-country,-town",
				f: fields.FromPairs([]fields.Field{
					{Name: "country", Kind: gocipe.String},
					{Name: "town", Kind: gocipe.String},
					{Name: "name", Kind: gocipe.String},
				}),
			},
			want: []OrderBy{
				{Attribute: "country", Sort: Descending},
				{Attribute: "town", Sort: Descending},
			},
			wantErr: false,
		},
		{
			name: "2 asc 1 desc",
			args: args{
				s: "country,-town,name",
				f: fields.FromPairs([]fields.Field{
					{Name: "country", Kind: gocipe.String},
					{Name: "town", Kind: gocipe.String},
					{Name: "name", Kind: gocipe.String},
				}),
			},
			want: []OrderBy{
				{Attribute: "country", Sort: Ascending},
				{Attribute: "town", Sort: Descending},
				{Attribute: "name", Sort: Ascending},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := OrderByFromString(tt.args.s, tt.args.f)
			if (err != nil) != tt.wantErr {
				t.Errorf("OrderByFromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			compareOrderBys(t, got, tt.want)
		})
	}
}

func TestConditionsFromMap(t *testing.T) {
	type args struct {
		m map[string][]string
		f fields.Fields
	}

	tests := []struct {
		name    string
		args    args
		want    []Condition
		wantErr bool
	}{
		{
			name: "No fields",
			args: args{
				m: map[string][]string{
					"name": {"foobar"},
				},
				f: fields.FromPairs([]fields.Field{}),
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "Empty Map",
			args: args{
				m: map[string][]string{},
				f: fields.FromPairs([]fields.Field{
					{Name: "name", Kind: gocipe.String},
				}),
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "Bool err 1",
			args: args{
				m: map[string][]string{
					"active": {"True"},
				},
				f: fields.FromPairs([]fields.Field{
					{Name: "active", Kind: gocipe.Bool},
				}),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Bool err 2",
			args: args{
				m: map[string][]string{
					"active": {"No"},
				},
				f: fields.FromPairs([]fields.Field{
					{Name: "active", Kind: gocipe.Bool},
				}),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Bool err 3",
			args: args{
				m: map[string][]string{
					"active": {"0"},
				},
				f: fields.FromPairs([]fields.Field{
					{Name: "active", Kind: gocipe.Bool},
				}),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Bool err 4",
			args: args{
				m: map[string][]string{
					"active": {"gt:true"},
				},
				f: fields.FromPairs([]fields.Field{
					{Name: "active", Kind: gocipe.Bool},
				}),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Bool true",
			args: args{
				m: map[string][]string{
					"active": {"true"},
				},
				f: fields.FromPairs([]fields.Field{
					{Name: "active", Kind: gocipe.Bool},
				}),
			},
			want: []Condition{
				{Attribute: "active", Operator: Equals, Value: true},
			},
			wantErr: false,
		},
		{
			name: "Bool false",
			args: args{
				m: map[string][]string{
					"active": {"false"},
				},
				f: fields.FromPairs([]fields.Field{
					{Name: "active", Kind: gocipe.Bool},
				}),
			},
			want: []Condition{
				{Attribute: "active", Operator: Equals, Value: false},
			},
			wantErr: false,
		},
		{
			name: "String",
			args: args{
				m: map[string][]string{
					"name": {"foo"},
				},
				f: fields.FromPairs([]fields.Field{
					{Name: "name", Kind: gocipe.String},
				}),
			},
			want: []Condition{
				{Attribute: "name", Operator: Equals, Value: "foo"},
			},
			wantErr: false,
		},
		{
			name: "String err",
			args: args{
				m: map[string][]string{
					"name": {"foo:"},
				},
				f: fields.FromPairs([]fields.Field{
					{Name: "name", Kind: gocipe.String},
				}),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Int64",
			args: args{
				m: map[string][]string{
					"price": {"12500"},
				},
				f: fields.FromPairs([]fields.Field{
					{Name: "price", Kind: gocipe.Int64},
				}),
			},
			want: []Condition{
				{Attribute: "price", Operator: Equals, Value: int64(12500)},
			},
			wantErr: false,
		},
		{
			name: "Int64 err",
			args: args{
				m: map[string][]string{
					"price": {":12500"},
				},
				f: fields.FromPairs([]fields.Field{
					{Name: "price", Kind: gocipe.Int64},
				}),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Float64",
			args: args{
				m: map[string][]string{
					"pi": {"3.14159265"},
				},
				f: fields.FromPairs([]fields.Field{
					{Name: "pi", Kind: gocipe.Float64},
				}),
			},
			want: []Condition{
				{Attribute: "pi", Operator: Equals, Value: 3.14159265},
			},
			wantErr: false,
		},
		{
			name: "Equals",
			args: args{
				m: map[string][]string{
					"stock": {"eq:100"},
				},
				f: fields.FromPairs([]fields.Field{
					{Name: "stock", Kind: gocipe.Int64},
				}),
			},
			want: []Condition{
				{Attribute: "stock", Operator: Equals, Value: int64(100)},
			},
			wantErr: false,
		},
		{
			name: "NotEquals",
			args: args{
				m: map[string][]string{
					"pi": {"ne:3.14159265"},
				},
				f: fields.FromPairs([]fields.Field{
					{Name: "pi", Kind: gocipe.Float64},
				}),
			},
			want: []Condition{
				{Attribute: "pi", Operator: NotEquals, Value: 3.14159265},
			},
			wantErr: false,
		},
		{
			name: "GreaterThan",
			args: args{
				m: map[string][]string{
					"pi": {"gt:3.14159265"},
				},
				f: fields.FromPairs([]fields.Field{
					{Name: "pi", Kind: gocipe.Float64},
				}),
			},
			want: []Condition{
				{Attribute: "pi", Operator: GreaterThan, Value: 3.14159265},
			},
			wantErr: false,
		},
		{
			name: "GreaterOrEqualTo",
			args: args{
				m: map[string][]string{
					"pi": {"gte:3.14159265"},
				},
				f: fields.FromPairs([]fields.Field{
					{Name: "pi", Kind: gocipe.Float64},
				}),
			},
			want: []Condition{
				{Attribute: "pi", Operator: GreaterOrEqualTo, Value: 3.14159265},
			},
			wantErr: false,
		},
		{
			name: "LessThan",
			args: args{
				m: map[string][]string{
					"stock": {"lt:10"},
				},
				f: fields.FromPairs([]fields.Field{
					{Name: "stock", Kind: gocipe.Int64},
				}),
			},
			want: []Condition{
				{Attribute: "stock", Operator: LessThan, Value: int64(10)},
			},
			wantErr: false,
		},
		{
			name: "LessOrEqualTo",
			args: args{
				m: map[string][]string{
					"created_at": {"lte:2006-01-02T15:04:05Z07:00"},
				},
				f: fields.FromPairs([]fields.Field{
					{Name: "created_at", Kind: gocipe.String},
				}),
			},
			want: []Condition{
				{Attribute: "created_at", Operator: LessOrEqualTo, Value: "2006-01-02T15:04:05Z07:00"},
			},
			wantErr: false,
		},
		{
			name: "[]Values",
			args: args{
				m: map[string][]string{
					"created_at": {"lte:2006-01-02T15:04:05Z07:00", "lte:2006-01-02T15:04:05Z07:00"},
				},
				f: fields.FromPairs([]fields.Field{
					{Name: "created_at", Kind: gocipe.String},
				}),
			},
			want:    []Condition{},
			wantErr: false,
		},
		{
			name: "Combination 1",
			args: args{
				m: map[string][]string{
					"stock": {"lt:10"},
					"pi":    {"gt:3.14159265"},
				},
				f: fields.FromPairs([]fields.Field{
					{Name: "stock", Kind: gocipe.Int64},
					{Name: "pi", Kind: gocipe.Float64},
				}),
			},
			want: []Condition{
				{Attribute: "stock", Operator: LessThan, Value: int64(10)},
				{Attribute: "pi", Operator: GreaterThan, Value: float64(3.14159265)},
			},
			wantErr: false,
		},
		{
			name: "Combination 2",
			args: args{
				m: map[string][]string{
					"created_at": {"lte:2006-01-02T15:04:05Z07:00"},
					"active":     {"true"},
				},
				f: fields.FromPairs([]fields.Field{
					{Name: "active", Kind: gocipe.Bool},
					{Name: "created_at", Kind: gocipe.String},
				}),
			},
			want: []Condition{
				{Attribute: "active", Operator: Equals, Value: true},
				{Attribute: "created_at", Operator: LessOrEqualTo, Value: "2006-01-02T15:04:05Z07:00"},
			},
			wantErr: false,
		},
		{
			name: "Combination 3",
			args: args{
				m: map[string][]string{
					"stock":      {"lt:10"},
					"pi":         {"gt:3.14159265"},
					"created_at": {"lte:2006-01-02T15:04:05Z07:00"},
				},
				f: fields.FromPairs([]fields.Field{
					{Name: "created_at", Kind: gocipe.String},
					{Name: "pi", Kind: gocipe.Float64},
					{Name: "stock", Kind: gocipe.Int64},
				}),
			},
			want: []Condition{
				{Attribute: "created_at", Operator: LessOrEqualTo, Value: "2006-01-02T15:04:05Z07:00"},
				{Attribute: "pi", Operator: GreaterThan, Value: float64(3.14159265)},
				{Attribute: "stock", Operator: LessThan, Value: int64(10)},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConditionsFromMap(tt.args.m, tt.args.f)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConditionsFromMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			compareConditions(t, got, tt.want)
		})
	}
}
