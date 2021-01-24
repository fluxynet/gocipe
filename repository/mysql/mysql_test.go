package mysql

import (
	"github.com/fluxynet/gocipe"
	"github.com/fluxynet/gocipe/fields"
	"github.com/fluxynet/gocipe/repository"
	"github.com/fluxynet/gocipe/values"
	"reflect"
	"testing"
	"time"
)

func compareSlicesOfInterface(t *testing.T, got, want []interface{}) {
	if len(got) != len(want) {
		t.Errorf("ConditionsToWhere() len(gotArgs) = %d, len(want) %d", len(got), len(want))
		return
	}

	for i := range got {
		if !reflect.DeepEqual(got[i], want[i]) {
			t.Errorf("ConditionsToWhere() got[%d] = %v, want[%d] = %v", i, got[i], i, want[i])
			return
		}
	}
}

func compareQueries(t *testing.T, got, want Query) {
	if got.SQL != want.SQL {
		t.Errorf("got:\n\t%s\nwant:\n\t%s", got.SQL, want.SQL)
		return
	}

	compareSlicesOfInterface(t, got.Args, want.Args)
}

func TestConditionsToWhere(t *testing.T) {
	type args struct {
		c []repository.Condition
	}

	tests := []struct {
		name     string
		args     args
		wantSQL  string
		wantArgs []interface{}
	}{
		{
			name: "Empty",
			args: args{
				c: []repository.Condition{},
			},
			wantSQL:  "",
			wantArgs: nil,
		},
		{
			name: "Equals",
			args: args{
				c: []repository.Condition{
					{
						Property: "name",
						Operator: repository.Equals,
						Value:    "Foo",
						Type:     repository.And,
					},
				},
			},
			wantSQL:  " WHERE `name` = ?",
			wantArgs: []interface{}{"Foo"},
		},
		{
			name: "NotEquals",
			args: args{
				c: []repository.Condition{
					{
						Property: "name",
						Operator: repository.NotEquals,
						Value:    "Foo",
						Type:     repository.And,
					},
				},
			},
			wantSQL:  " WHERE `name` <> ?",
			wantArgs: []interface{}{"Foo"},
		},
		{
			name: "GreaterThan",
			args: args{
				c: []repository.Condition{
					{
						Property: "age",
						Operator: repository.GreaterThan,
						Value:    18,
						Type:     repository.And,
					},
				},
			},
			wantSQL:  " WHERE `age` > ?",
			wantArgs: []interface{}{18},
		},
		{
			name: "GreaterOrEqualTo",
			args: args{
				c: []repository.Condition{
					{
						Property: "age",
						Operator: repository.GreaterOrEqualTo,
						Value:    18,
						Type:     repository.And,
					},
				},
			},
			wantSQL:  " WHERE `age` >= ?",
			wantArgs: []interface{}{18},
		},
		{
			name: "LessThan",
			args: args{
				c: []repository.Condition{
					{
						Property: "age",
						Operator: repository.LessThan,
						Value:    18,
						Type:     repository.And,
					},
				},
			},
			wantSQL:  " WHERE `age` < ?",
			wantArgs: []interface{}{18},
		},
		{
			name: "LessOrEqualTo",
			args: args{
				c: []repository.Condition{
					{
						Property: "age",
						Operator: repository.LessOrEqualTo,
						Value:    18,
						Type:     repository.And,
					},
				},
			},
			wantSQL:  " WHERE `age` <= ?",
			wantArgs: []interface{}{18},
		},
		{
			name: "Combination LessOrEqualTo, Equals",
			args: args{
				c: []repository.Condition{
					{
						Property: "age",
						Operator: repository.LessOrEqualTo,
						Value:    18,
						Type:     repository.And,
					},
					{
						Property: "name",
						Operator: repository.Equals,
						Value:    "Foo",
						Type:     repository.And,
					},
				},
			},
			wantSQL:  " WHERE `age` <= ? AND `name` = ?",
			wantArgs: []interface{}{18, "Foo"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSQL, gotArgs := ConditionsToWhere(tt.args.c)
			if gotSQL != tt.wantSQL {
				t.Errorf("ConditionsToWhere() got = %v, want %v", gotSQL, tt.wantSQL)
			}

			compareSlicesOfInterface(t, gotArgs, tt.wantArgs)
		})
	}
}

func TestCreate(t *testing.T) {
	type args struct {
		entity string
		vals   values.Values
	}

	tests := []struct {
		name string
		args args
		want Query
	}{
		{
			name: "No values",
			args: args{
				entity: "foo",
				vals:   values.Values{},
			},
			want: Query{
				SQL:  "",
				Args: nil,
			},
		},
		{
			name: "No entity",
			args: args{
				entity: "",
				vals: values.FromPairs([]values.Value{
					{Name: "name", Value: "John Doe"},
				}),
			},
			want: Query{
				SQL:  "",
				Args: nil,
			},
		},
		{
			name: "No values, no entity",
			args: args{
				entity: "",
				vals:   values.Values{},
			},
			want: Query{
				SQL:  "",
				Args: nil,
			},
		},
		{
			name: "1 value",
			args: args{
				entity: "foo",
				vals: values.FromPairs([]values.Value{
					{Name: "age", Value: 18},
				}),
			},
			want: Query{
				SQL:  "INSERT INTO `foo` (`age`) VALUES (?)",
				Args: []interface{}{18},
			},
		},
		{
			name: "2 values",
			args: args{
				entity: "bar",
				vals: values.FromPairs([]values.Value{
					{Name: "country", Value: "MU"},
					{Name: "status", Value: true},
				}),
			},
			want: Query{
				SQL:  "INSERT INTO `bar` (`country`,`status`) VALUES (?,?)",
				Args: []interface{}{"MU", true},
			},
		},
		{
			name: "3 values",
			args: args{
				entity: "products",
				vals: values.FromPairs([]values.Value{
					{Name: "name", Value: "Apple"},
					{Name: "price", Value: 3.5},
					{Name: "color", Value: "red"},
				}),
			},
			want: Query{
				SQL:  "INSERT INTO `products` (`name`,`price`,`color`) VALUES (?,?,?)",
				Args: []interface{}{"Apple", 3.5, "red"},
			},
		},
		{
			name: "3 values re-ordered",
			args: args{
				entity: "products",
				vals: values.FromPairs([]values.Value{
					{Name: "name", Value: "Apple"},
					{Name: "color", Value: "red"},
					{Name: "price", Value: 3.5},
				}),
			},
			want: Query{
				SQL:  "INSERT INTO `products` (`name`,`color`,`price`) VALUES (?,?,?)",
				Args: []interface{}{"Apple", "red", 3.5},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Create(tt.args.entity, tt.args.vals)

			compareQueries(t, got, tt.want)
		})
	}
}

func TestDelete(t *testing.T) {
	type args struct {
		entity string
		id     string
	}

	tests := []struct {
		name string
		args args
		want Query
	}{
		{
			name: "No entity",
			args: args{
				entity: "",
				id:     "00000000-0000-0000-0000-000000000001",
			},
			want: Query{
				SQL:  "",
				Args: nil,
			},
		},
		{
			name: "Empty id",
			args: args{
				entity: "",
				id:     "",
			},
			want: Query{
				SQL:  "",
				Args: nil,
			},
		},
		{
			name: "No entity and Empty id",
			args: args{
				entity: "",
				id:     "",
			},
			want: Query{
				SQL:  "",
				Args: nil,
			},
		},
		{
			name: "Non-empty id",
			args: args{
				entity: "products",
				id:     "00000000-0000-0000-0000-00000000000f",
			},
			want: Query{
				SQL:  "DELETE FROM `products` WHERE `id` = ?",
				Args: []interface{}{"00000000-0000-0000-0000-00000000000f"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Delete(tt.args.entity, tt.args.id); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Delete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeleteWhere(t *testing.T) {
	var now = time.Now()

	type args struct {
		entity string
		c      []repository.Condition
	}

	tests := []struct {
		name string
		args args
		want Query
	}{
		{
			name: "No entity",
			args: args{
				entity: "",
				c:      nil,
			},
			want: Query{
				SQL:  "",
				Args: nil,
			},
		},
		{
			name: "No conditions",
			args: args{
				entity: "products",
				c:      nil,
			},
			want: Query{
				SQL:  "DELETE FROM `products`",
				Args: nil,
			},
		},
		{
			name: "No entity and no conditions",
			args: args{
				entity: "",
				c:      nil,
			},
			want: Query{
				SQL:  "",
				Args: nil,
			},
		},
		{
			name: "Entity and 1 condition",
			args: args{
				entity: "persons",
				c: []repository.Condition{
					{
						Property: "age",
						Operator: repository.LessThan,
						Value:    18,
						Type:     repository.And,
					},
				},
			},
			want: Query{
				SQL:  "DELETE FROM `persons` WHERE `age` < ?",
				Args: []interface{}{18},
			},
		},
		{
			name: "Entity and 2 conditions",
			args: args{
				entity: "products",
				c: []repository.Condition{
					{
						Property: "expiry_date",
						Operator: repository.LessOrEqualTo,
						Value:    now,
						Type:     repository.And,
					},
					{
						Property: "in_stock",
						Operator: repository.Equals,
						Value:    true,
						Type:     repository.And,
					},
				},
			},
			want: Query{
				SQL:  "DELETE FROM `products` WHERE `expiry_date` <= ? AND `in_stock` = ?",
				Args: []interface{}{now, true},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DeleteWhere(tt.args.entity, tt.args.c...)
			compareQueries(t, got, tt.want)
		})
	}
}

func TestGet(t *testing.T) {
	type args struct {
		f      fields.Fields
		entity string
		id     string
	}

	tests := []struct {
		name string
		args args
		want Query
	}{
		{
			name: "No entity",
			args: args{
				f: fields.FromPairs([]fields.Field{
					{Name: "name", Kind: gocipe.String},
					{Name: "price", Kind: gocipe.Float64},
				}),
				entity: "",
				id:     "00000000-0000-0000-0000-000000000001",
			},
			want: Query{
				SQL:  "",
				Args: nil,
			},
		},
		{
			name: "No id",
			args: args{
				f: fields.FromPairs([]fields.Field{
					{Name: "name", Kind: gocipe.String},
					{Name: "price", Kind: gocipe.Float64},
				}),
				entity: "products",
				id:     "",
			},
			want: Query{
				SQL:  "",
				Args: nil,
			},
		},
		{
			name: "No entity and no id",
			args: args{
				f: fields.FromPairs([]fields.Field{
					{Name: "name", Kind: gocipe.String},
					{Name: "price", Kind: gocipe.Float64},
				}),
				entity: "",
				id:     "",
			},
			want: Query{
				SQL:  "",
				Args: nil,
			},
		},
		{
			name: "Entity and id",
			args: args{
				f: fields.FromPairs([]fields.Field{
					{Name: "name", Kind: gocipe.String},
					{Name: "price", Kind: gocipe.Float64},
				}),
				entity: "products",
				id:     "00000000-0000-0000-0000-000000000001",
			},
			want: Query{
				SQL:  "SELECT `name`,`price` FROM `products` WHERE `id` = ?",
				Args: []interface{}{"00000000-0000-0000-0000-000000000001"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Get(tt.args.entity, tt.args.f, tt.args.id)
			compareQueries(t, got, tt.want)
		})
	}
}

func TestList(t *testing.T) {
	type args struct {
		entity string
		f      fields.Fields
		p      repository.Pagination
		c      []repository.Condition
	}

	tests := []struct {
		name string
		args args
		want Query
	}{
		{
			name: "no entity",
			args: args{
				entity: "",
				f: fields.FromPairs([]fields.Field{
					{Name: "name", Kind: gocipe.String},
					{Name: "price", Kind: gocipe.Float64},
				}),
				p: repository.Pagination{Order: []repository.OrderBy{
					{Attribute: "age", Sort: repository.Ascending},
				}, Limit: 0},
				c: []repository.Condition{
					{
						Property: "name",
						Operator: repository.Equals,
						Value:    "foo",
					},
				},
			},
			want: Query{
				SQL:  "",
				Args: nil,
			},
		},
		{
			name: "no condition, no pagination",
			args: args{
				entity: "products",
				f: fields.FromPairs([]fields.Field{
					{Name: "name", Kind: gocipe.String},
					{Name: "price", Kind: gocipe.Float64},
				}),
			},
			want: Query{
				SQL:  "SELECT `name`,`price` FROM `products`",
				Args: nil,
			},
		},
		{
			name: "no condition, no order, no offset, limit 10",
			args: args{
				entity: "products",
				f: fields.FromPairs([]fields.Field{
					{Name: "color", Kind: gocipe.String},
					{Name: "stock", Kind: gocipe.Int64},
				}),
				p: repository.Pagination{Limit: 10},
				c: []repository.Condition{},
			},
			want: Query{
				SQL:  "SELECT `color`,`stock` FROM `products` LIMIT 10",
				Args: nil,
			},
		},
		{
			name: "no condition, no order, offset 5",
			args: args{
				entity: "products",
				f: fields.FromPairs([]fields.Field{
					{Name: "name", Kind: gocipe.String},
					{Name: "price", Kind: gocipe.Float64},
				}),
				p: repository.Pagination{Limit: 10, Offset: 5},
				c: []repository.Condition{},
			},
			want: Query{
				SQL:  "SELECT `name`,`price` FROM `products` LIMIT 5,10",
				Args: nil,
			},
		},
		{
			name: "no condition, order desc",
			args: args{
				entity: "products",
				f: fields.FromPairs([]fields.Field{
					{Name: "color", Kind: gocipe.String},
					{Name: "weight", Kind: gocipe.Float64},
				}),
				p: repository.Pagination{Order: []repository.OrderBy{
					{
						Attribute: "age",
						Sort:      repository.Descending,
					},
				}, Limit: 5, Offset: 20},
			},
			want: Query{
				SQL:  "SELECT `color`,`weight` FROM `products` ORDER BY `age` DESC LIMIT 20,5",
				Args: nil,
			},
		},
		{
			name: "no condition, order asc",
			args: args{
				entity: "products",
				f: fields.FromPairs([]fields.Field{
					{Name: "name", Kind: gocipe.String},
					{Name: "price", Kind: gocipe.Float64},
				}),
				p: repository.Pagination{Order: []repository.OrderBy{
					{
						Attribute: "age",
						Sort:      repository.Ascending,
					},
				}, Limit: 5, Offset: 20},
			},
			want: Query{
				SQL:  "SELECT `name`,`price` FROM `products` ORDER BY `age` ASC LIMIT 20,5",
				Args: nil,
			},
		},
		{
			name: "no pagination",
			args: args{
				entity: "products",
				f: fields.FromPairs([]fields.Field{
					{Name: "name", Kind: gocipe.String},
					{Name: "price", Kind: gocipe.Float64},
				}),
				c: []repository.Condition{
					{
						Property: "price",
						Operator: repository.GreaterOrEqualTo,
						Value:    100.50,
					},
				},
			},
			want: Query{
				SQL:  "SELECT `name`,`price` FROM `products` WHERE `price` >= ?",
				Args: []interface{}{100.50},
			},
		},
		{
			name: "condition and pagination",
			args: args{
				entity: "students",
				f: fields.FromPairs([]fields.Field{
					{Name: "name", Kind: gocipe.String},
					{Name: "price", Kind: gocipe.Float64},
				}),
				p: repository.Pagination{Limit: 100, Offset: 500},
				c: []repository.Condition{
					{
						Property: "height",
						Operator: repository.LessOrEqualTo,
						Value:    170,
					},
				},
			},
			want: Query{
				SQL:  "SELECT `name`,`price` FROM `students` WHERE `height` <= ? LIMIT 500,100",
				Args: []interface{}{170},
			},
		},
		{
			name: "condition, pagination and order; no limit",
			args: args{
				entity: "customers",
				f: fields.FromPairs([]fields.Field{
					{Name: "name", Kind: gocipe.String},
					{Name: "age", Kind: gocipe.Int64},
				}),
				p: repository.Pagination{Order: []repository.OrderBy{
					{
						Attribute: "credits",
						Sort:      repository.Descending,
					},
					{
						Attribute: "name",
						Sort:      repository.Ascending,
					},
				}},
				c: []repository.Condition{
					{
						Property: "active",
						Operator: repository.Equals,
						Value:    false,
					},
				},
			},
			want: Query{
				SQL:  "SELECT `name`,`age` FROM `customers` WHERE `active` = ? ORDER BY `credits` DESC, `name` ASC",
				Args: []interface{}{false},
			},
		},
		{
			name: "condition, pagination, order and limit",
			args: args{
				entity: "customers",
				f: fields.FromPairs([]fields.Field{
					{Name: "name", Kind: gocipe.String},
					{Name: "credits", Kind: gocipe.Float64},
				}),
				p: repository.Pagination{Order: []repository.OrderBy{
					{
						Attribute: "credits",
						Sort:      repository.Descending,
					},
					{
						Attribute: "name",
						Sort:      repository.Ascending,
					},
				}, Limit: 50},
				c: []repository.Condition{
					{
						Property: "active",
						Operator: repository.Equals,
						Value:    false,
					},
				},
			},
			want: Query{
				SQL:  "SELECT `name`,`credits` FROM `customers` WHERE `active` = ? ORDER BY `credits` DESC, `name` ASC LIMIT 50",
				Args: []interface{}{false},
			},
		},
		{
			name: "condition, pagination, order, limit and offset",
			args: args{
				entity: "customers",
				f: fields.FromPairs([]fields.Field{
					{Name: "name", Kind: gocipe.String},
					{Name: "credits", Kind: gocipe.Float64},
				}),
				p: repository.Pagination{Order: []repository.OrderBy{
					{
						Attribute: "credits",
						Sort:      repository.Descending,
					},
					{
						Attribute: "name",
						Sort:      repository.Ascending,
					},
				}, Limit: 50, Offset: 1000},
				c: []repository.Condition{
					{
						Property: "active",
						Operator: repository.Equals,
						Value:    false,
					},
				},
			},
			want: Query{
				SQL:  "SELECT `name`,`credits` FROM `customers` WHERE `active` = ? ORDER BY `credits` DESC, `name` ASC LIMIT 1000,50",
				Args: []interface{}{false},
			},
		},
		{
			name: "2 conditions, pagination, order and limit",
			args: args{
				entity: "customers",
				f: fields.FromPairs([]fields.Field{
					{Name: "name", Kind: gocipe.String},
					{Name: "country", Kind: gocipe.String},
				}),
				p: repository.Pagination{Order: []repository.OrderBy{
					{
						Attribute: "credits",
						Sort:      repository.Descending,
					},
					{
						Attribute: "name",
						Sort:      repository.Ascending,
					},
				}, Limit: 50},
				c: []repository.Condition{
					{
						Property: "country",
						Operator: repository.Equals,
						Value:    "MU",
					},
					{
						Property: "active",
						Operator: repository.Equals,
						Value:    false,
					},
				},
			},
			want: Query{
				SQL:  "SELECT `name`,`country` FROM `customers` WHERE `country` = ? AND `active` = ? ORDER BY `credits` DESC, `name` ASC LIMIT 50",
				Args: []interface{}{"MU", false},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := List(tt.args.entity, tt.args.f, tt.args.p, tt.args.c...)
			compareQueries(t, got, tt.want)
		})
	}
}

func TestOperator(t *testing.T) {
	type args struct {
		op repository.ConditionOperator
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Equals",
			args: args{op: repository.Equals},
			want: "=",
		},
		{
			name: "NotEquals",
			args: args{op: repository.NotEquals},
			want: "<>",
		},
		{
			name: "GreaterThan",
			args: args{op: repository.GreaterThan},
			want: ">",
		},
		{
			name: "GreaterOrEqualTo",
			args: args{op: repository.GreaterOrEqualTo},
			want: ">=",
		},
		{
			name: "LessThan",
			args: args{op: repository.LessThan},
			want: "<",
		},
		{
			name: "LessOrEqualTo",
			args: args{op: repository.LessOrEqualTo},
			want: "<=",
		},
		{
			name: "In",
			args: args{op: repository.In},
			want: "IN",
		},
		{
			name: "NotIn",
			args: args{op: repository.NotIn},
			want: "NOT IN",
		},
		{
			name: "Unknown",
			args: args{op: 100},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Operator(tt.args.op); got != tt.want {
				t.Errorf("Operator() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSortToString(t *testing.T) {
	type args struct {
		o repository.OrderSort
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Ascending",
			args: args{o: repository.Ascending},
			want: "ASC",
		},
		{
			name: "Descending",
			args: args{o: repository.Descending},
			want: "DESC",
		},
		{
			name: "Unknown",
			args: args{o: 10},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SortToString(tt.args.o); got != tt.want {
				t.Errorf("SortToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTypeToString(t *testing.T) {
	type args struct {
		t repository.ConditionType
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "And",
			args: args{t: repository.And},
			want: "AND",
		},
		{
			name: "Or",
			args: args{t: repository.Or},
			want: "OR",
		},
		{
			name: "Unknown",
			args: args{t: repository.ConditionType(10)},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TypeToString(tt.args.t); got != tt.want {
				t.Errorf("TypeToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	type args struct {
		entity string
		id     string
		vals   values.Values
	}
	tests := []struct {
		name string
		args args
		want Query
	}{
		{
			name: "no entity",
			args: args{
				entity: "",
				id:     "00000000-0000-0000-0000-000000000001",
				vals: values.FromMap(map[string]interface{}{
					"name": "foobar",
				}),
			},
			want: Query{
				SQL:  "",
				Args: nil,
			},
		},
		{
			name: "no id",
			args: args{
				entity: "animals",
				id:     "",
				vals:   values.Values{},
			},
			want: Query{
				SQL:  "",
				Args: nil,
			},
		},
		{
			name: "no entity, no id",
			args: args{
				entity: "",
				id:     "",
				vals:   values.Values{},
			},
			want: Query{
				SQL:  "",
				Args: nil,
			},
		},
		{
			name: "entity, id, no values",
			args: args{
				entity: "farm",
				id:     "00000000-0000-0000-0000-000000000001",
				vals:   values.Values{},
			},
			want: Query{
				SQL:  "",
				Args: nil,
			},
		},
		{
			name: "entity, id, values 1",
			args: args{
				entity: "farm",
				id:     "00000000-0000-0000-0000-000000000001",
				vals: values.FromMap(map[string]interface{}{
					"name": "Animal",
				}),
			},
			want: Query{
				SQL: "UPDATE `farm` SET `name` = ? WHERE `id` = ?",
				Args: []interface{}{
					"Animal", "00000000-0000-0000-0000-000000000001",
				},
			},
		},
		{
			name: "entity, id, values 2",
			args: args{
				entity: "animals",
				id:     "00000000-0000-0000-0000-000000000002",
				vals: values.FromMap(map[string]interface{}{
					"role": "president",
					"legs": 2,
				}),
			},
			want: Query{
				SQL: "UPDATE `animals` SET `role` = ?, `legs` = ? WHERE `id` = ?",
				Args: []interface{}{
					"president", 2, "00000000-0000-0000-0000-000000000002",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Update(tt.args.entity, tt.args.id, tt.args.vals)
			compareQueries(t, got, tt.want)
		})
	}
}

func TestUpdateWhere(t *testing.T) {
	type args struct {
		entity string
		vals   values.Values
		c      []repository.Condition
	}
	tests := []struct {
		name string
		args args
		want Query
	}{
		{
			name: "no entity",
			args: args{
				entity: "",
				vals: values.FromMap(map[string]interface{}{
					"active": false,
				}),
				c: []repository.Condition{
					{
						Property: "credit",
						Operator: repository.Equals,
						Value:    0,
					},
				},
			},
			want: Query{
				SQL:  "",
				Args: nil,
			},
		},
		{
			name: "no values",
			args: args{
				entity: "customers",
				vals:   values.Values{},
				c: []repository.Condition{
					{
						Property: "credit",
						Operator: repository.Equals,
						Value:    0,
					},
				},
			},
			want: Query{
				SQL:  "",
				Args: nil,
			},
		},
		{
			name: "no entity, no values",
			args: args{
				entity: "",
				vals:   values.Values{},
				c: []repository.Condition{
					{
						Property: "credit",
						Operator: repository.Equals,
						Value:    0,
					},
				},
			},
			want: Query{
				SQL:  "",
				Args: nil,
			},
		},
		{
			name: "no conditions",
			args: args{
				entity: "customers",
				vals: values.FromMap(map[string]interface{}{
					"active": false,
				}),
				c: nil,
			},
			want: Query{
				SQL:  "UPDATE `customers` SET `active` = ?",
				Args: []interface{}{false},
			},
		},
		{
			name: "no entity, no values, no conditions",
			args: args{
				entity: "",
				vals:   values.Values{},
				c:      nil,
			},
			want: Query{
				SQL:  "",
				Args: nil,
			},
		},
		{
			name: "1 condition",
			args: args{
				entity: "customers",
				vals: values.FromMap(map[string]interface{}{
					"active": false,
				}),
				c: []repository.Condition{
					{
						Property: "credit",
						Operator: repository.Equals,
						Value:    0,
					},
				},
			},
			want: Query{
				SQL:  "UPDATE `customers` SET `active` = ? WHERE `credit` = ?",
				Args: []interface{}{false, 0},
			},
		},
		{
			name: "3 conditions",
			args: args{
				entity: "customers",
				vals: values.FromMap(map[string]interface{}{
					"active": true,
				}),
				c: []repository.Condition{
					{
						Property: "credit",
						Operator: repository.NotEquals,
						Value:    0,
					},
					{
						Property: "service_count",
						Operator: repository.GreaterOrEqualTo,
						Value:    1,
					},
				},
			},
			want: Query{
				SQL:  "UPDATE `customers` SET `active` = ? WHERE `credit` <> ? AND `service_count` >= ?",
				Args: []interface{}{true, 0, 1},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := UpdateWhere(tt.args.entity, tt.args.vals, tt.args.c...)
			compareQueries(t, got, tt.want)
		})
	}
}

func TestValuesToSet(t *testing.T) {
	type args struct {
		vals values.Values
	}
	tests := []struct {
		name     string
		args     args
		wantSet  string
		wantArgs []interface{}
	}{
		{
			name: "0 values",
			args: args{
				vals: values.Values{},
			},
			wantSet:  "",
			wantArgs: nil,
		},
		{
			name: "1 value",
			args: args{
				vals: values.FromPairs([]values.Value{
					{Name: "name", Value: "foo"},
				}),
			},
			wantSet:  "SET `name` = ?",
			wantArgs: []interface{}{"foo"},
		},
		{
			name: "2 values",
			args: args{
				vals: values.FromPairs([]values.Value{
					{Name: "name", Value: "foo"},
					{Name: "age", Value: 18},
				}),
			},
			wantSet:  "SET `name` = ?, `age` = ?",
			wantArgs: []interface{}{"foo", 18},
		},
		{
			name: "3 values",
			args: args{
				vals: values.FromPairs([]values.Value{
					{Name: "name", Value: "foo"},
					{Name: "age", Value: 18},
					{Name: "active", Value: true},
				}),
			},
			wantSet:  "SET `name` = ?, `age` = ?, `active` = ?",
			wantArgs: []interface{}{"foo", 18, true},
		},
		{
			name: "3 values reordered",
			args: args{
				vals: values.FromPairs([]values.Value{
					{Name: "active", Value: true},
					{Name: "name", Value: "foo"},
					{Name: "age", Value: 18},
				}),
			},
			wantSet:  "SET `active` = ?, `name` = ?, `age` = ?",
			wantArgs: []interface{}{true, "foo", 18},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSet, gotArgs := ValuesToSet(tt.args.vals)
			if gotSet != tt.wantSet {
				t.Errorf("ValuesToSet() gotSet = %v, want %v", gotSet, tt.wantSet)
			}

			compareSlicesOfInterface(t, gotArgs, tt.wantArgs)
		})
	}
}