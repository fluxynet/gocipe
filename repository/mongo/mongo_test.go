package mongo

import (
	"reflect"
	"testing"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/fluxynet/gocipe/repository"
)

func compareBsonD(t *testing.T, got, want bson.D) {
	var lg = len(got)
	var lw = len(want)

	if lg != lw {
		t.Errorf("Length got = %d want = %d\n", lg, lw)
		return
	}

	for i := range got {
		if !reflect.DeepEqual(got[i], want[i]) {
			t.Errorf("%d\n\tgot  = %v\n\twant = %v\n", i, got[i], want[i])
			return
		}
	}
}

func TestConditionsToBsonD(t *testing.T) {
	type args struct {
		c []repository.Condition
	}
	tests := []struct {
		name    string
		args    args
		want    bson.D
		wantErr bool
	}{
		{
			name: "Empty",
			args: args{
				c: nil,
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "Invalid Condition",
			args: args{
				c: []repository.Condition{
					{
						Attribute: "name",
						Operator:  repository.ConditionOperator(10),
						Value:     "foo",
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Invalid + Valid Condition",
			args: args{
				c: []repository.Condition{
					{
						Attribute: "name",
						Operator:  repository.ConditionOperator(100),
						Value:     "foo",
					},
					{
						Attribute: "age",
						Operator:  repository.GreaterOrEqualTo,
						Value:     18,
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Invalid Attribute",
			args: args{
				c: []repository.Condition{
					{
						Attribute: "",
						Operator:  repository.Equals,
						Value:     10,
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Invalid + Valid Attribute",
			args: args{
				c: []repository.Condition{
					{
						Attribute: "",
						Operator:  repository.Equals,
						Value:     10,
					},
					{
						Attribute: "name",
						Operator:  repository.Equals,
						Value:     "foo",
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Equals",
			args: args{
				c: []repository.Condition{
					{
						Attribute: "name",
						Operator:  repository.Equals,
						Value:     "foo",
					},
				},
			},
			want: bson.D{
				bson.E{
					Key:   "name",
					Value: bson.M{"$eq": "foo"},
				},
			},
			wantErr: false,
		},
		{
			name: "NotEquals",
			args: args{
				c: []repository.Condition{
					{
						Attribute: "active",
						Operator:  repository.NotEquals,
						Value:     false,
					},
				},
			},
			want: bson.D{
				bson.E{
					Key:   "active",
					Value: bson.M{"$ne": false},
				},
			},
			wantErr: false,
		},
		{
			name: "GreaterThan",
			args: args{
				c: []repository.Condition{
					{
						Attribute: "age",
						Operator:  repository.GreaterThan,
						Value:     18,
					},
				},
			},
			want: bson.D{
				bson.E{
					Key:   "age",
					Value: bson.M{"$gt": 18},
				},
			},
			wantErr: false,
		},
		{
			name: "GreaterOrEqualTo",
			args: args{
				c: []repository.Condition{
					{
						Attribute: "value",
						Operator:  repository.GreaterOrEqualTo,
						Value:     100,
					},
				},
			},
			want: bson.D{
				bson.E{
					Key:   "value",
					Value: bson.M{"$gte": 100},
				},
			},
			wantErr: false,
		},
		{
			name: "LessThan",
			args: args{
				c: []repository.Condition{
					{
						Attribute: "x",
						Operator:  repository.LessThan,
						Value:     0,
					},
				},
			},
			want: bson.D{
				bson.E{
					Key:   "x",
					Value: bson.M{"$lt": 0},
				},
			},
			wantErr: false,
		},
		{
			name: "LessOrEqualTo",
			args: args{
				c: []repository.Condition{
					{
						Attribute: "y",
						Operator:  repository.LessOrEqualTo,
						Value:     50,
					},
				},
			},
			want: bson.D{
				bson.E{
					Key:   "y",
					Value: bson.M{"$lte": 50},
				},
			},
			wantErr: false,
		},
		//{
		//	name: "Like",
		//	args: args{
		//		c: []repository.Condition{
		//			{
		//				Attribute: "name",
		//				Operator:  repository.Like,
		//				Value:     50,
		//			},
		//		},
		//	},
		//	want: bson.D{
		//
		//	},
		//	wantErr: false,
		//},
		//{
		//	name: "In",
		//	args: args{
		//		c: []repository.Condition{
		//
		//		},
		//	},
		//	want: bson.D{
		//
		//	},
		//	wantErr: false,
		//},
		//{
		//	name: "NotIn",
		//	args: args{
		//		c: []repository.Condition{
		//
		//		},
		//	},
		//	want: bson.D{
		//
		//	},
		//	wantErr: false,
		//},
		//{
		//	name: "Combination 1",
		//	args: args{
		//		c: []repository.Condition{
		//
		//		},
		//	},
		//	want: bson.D{
		//
		//	},
		//	wantErr: false,
		//},
		//{
		//	name: "Combination 2",
		//	args: args{
		//		c: []repository.Condition{
		//
		//		},
		//	},
		//	want: bson.D{
		//
		//	},
		//	wantErr: false,
		//},
		//{
		//	name: "Combination 3",
		//	args: args{
		//		c: []repository.Condition{
		//
		//		},
		//	},
		//	want: bson.D{
		//
		//	},
		//	wantErr: false,
		//},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConditionsToBsonD(tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConditionsToBsonD() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			compareBsonD(t, got, tt.want)
		})
	}
}
