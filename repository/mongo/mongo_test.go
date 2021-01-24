package mongo

import (
	"github.com/fluxynet/gocipe/repository"
	"github.com/fluxynet/gocipe/values"
	"go.mongodb.org/mongo-driver/bson"
	"reflect"
	"testing"
)

func TestValuesToBsonM(t *testing.T) {
	type args struct {
		vals values.Values
	}
	tests := []struct {
		name string
		args args
		want bson.M
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValuesToBsonM(tt.args.vals); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValuesToBsonM() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConditionsToBsonD(t *testing.T) {
	type args struct {
		c []repository.Condition
	}
	tests := []struct {
		name string
		args args
		want bson.D
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConditionsToBsonD(tt.args.c); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConditionsToBsonD() = %v, want %v", got, tt.want)
			}
		})
	}
}
