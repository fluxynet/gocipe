package mysql

import (
	"strconv"
	"strings"

	"github.com/fluxynet/gocipe/types"
	"github.com/fluxynet/gocipe/types/fields"
	"github.com/fluxynet/gocipe/types/fields/entity"
	"github.com/fluxynet/gocipe/values"

	"github.com/fluxynet/gocipe/repository"
)

// Query represents a mysql to be executed
type Query struct {
	SQL  string
	Args []interface{}
}

// SelectFieldNames in mysql select format
func SelectFieldNames(f fields.Fields) string {
	var (
		b  strings.Builder
		it = f.Iterator()
		n  = f.Length()
	)

	for i := 1; it.Next(); i++ {
		b.WriteString("`")
		b.WriteString(it.Field().Name)
		b.WriteString("`")

		if i != n {
			b.WriteString(",")
		}
	}

	return b.String()
}

// Get generates Query for a SELECT operation (by id)
func Get(entity entity.Entity, id string) Query {
	var (
		name = entity.Name()
		f    = entity.Fields()
		n    = f.Length()
	)

	if n == 0 || name == "" || id == "" {
		return Query{}
	}

	return Query{
		SQL:  "SELECT " + SelectFieldNames(f) + " FROM `" + name + "` WHERE `id` = ?",
		Args: []interface{}{id},
	}
}

// Operator returns the mysql equivalent of a ConditionOperator
func Operator(op repository.ConditionOperator) string {
	switch op {
	case repository.Equals:
		return "="
	case repository.NotEquals:
		return "<>"
	case repository.GreaterThan:
		return ">"
	case repository.GreaterOrEqualTo:
		return ">="
	case repository.LessThan:
		return "<"
	case repository.LessOrEqualTo:
		return "<="
	case repository.Like:
		return "LIKE"
	case repository.In:
		return "IN"
	case repository.NotIn:
		return "NOT IN"
	}

	return ""
}

// TypeToString returns the mysql equivalent of condition types (AND / OR)
func TypeToString(t repository.ConditionType) string {
	switch t {
	case repository.And:
		return "AND"
	case repository.Or:
		return "OR"
	}

	return ""
}

// SortToString returns the mysql equivalent of OrderSort - Ascending / Descending order
func SortToString(o repository.OrderSort) string {
	switch o {
	case repository.Ascending:
		return "ASC"
	case repository.Descending:
		return "DESC"
	}

	return ""
}

// ConditionsToWhere returns the `WHERE` segment and arguments of a mysql query. includes preceding space and where.
// string part is empty string if no condition passed
// args is empty slice if no condition passed
func ConditionsToWhere(c []repository.Condition) (string, []interface{}) {
	var t = len(c)
	if t == 0 {
		return "", nil
	}

	var where strings.Builder
	var args = make([]interface{}, t)

	where.WriteString(" WHERE ")

	t -= 1
	for i := range c {
		where.WriteString("`")
		where.WriteString(c[i].Attribute)
		where.WriteString("` ")
		where.WriteString(Operator(c[i].Operator))
		where.WriteString(" ?")

		if i != t {
			where.WriteString(" AND ")
		}

		args[i] = c[i].Value
	}

	return where.String(), args
}

func PaginationToOrderBy(p repository.Pagination) string {
	var (
		b strings.Builder
		l = len(p.Order)
	)

	if l != 0 {
		b.WriteString(" ORDER BY ")
	}

	l -= 1
	for i := range p.Order {
		b.WriteString("`")
		b.WriteString(p.Order[i].Attribute)
		b.WriteString("` ")
		b.WriteString(SortToString(p.Order[i].Sort))

		if i != l {
			b.WriteString(", ")
		}
	}

	if p.Limit == 0 {
		return b.String()
	}

	if p.Offset == 0 {
		b.WriteString(" LIMIT ")
		b.WriteString(strconv.Itoa(p.Limit))
	} else {
		b.WriteString(" LIMIT ")
		b.WriteString(strconv.Itoa(p.Offset))
		b.WriteString(",")
		b.WriteString(strconv.Itoa(p.Limit))
	}

	return b.String()
}

// List returns a list of entities retrieved from mysql based on conditions
func List(entity entity.Entity, p repository.Pagination, c ...repository.Condition) Query {
	var (
		name = entity.Name()
		f    = entity.Fields()
		q    = Query{}
	)

	if name == "" || f.IsEmpty() {
		return q
	}

	var where, pagination string

	pagination = PaginationToOrderBy(p)

	where, q.Args = ConditionsToWhere(c)
	q.SQL = "SELECT " + SelectFieldNames(f) + " FROM `" + name + "`" + where + pagination

	return q
}

// Delete generates Query for a DELETE operation (by id)
func Delete(named repository.Named, id string) Query {
	var name = named.Name()
	if name == "" || id == "" {
		return Query{}
	}

	return Query{
		SQL:  "DELETE FROM `" + name + "` WHERE `id` = ?",
		Args: []interface{}{id},
	}
}

// DeleteWhere generates Query for a DELETE operation (based on 1 or more conditions)
func DeleteWhere(named repository.Named, c ...repository.Condition) Query {
	var name = named.Name()
	if name == "" {
		return Query{}
	}

	var sql = "DELETE FROM `" + name + "`"
	var where, args = ConditionsToWhere(c)

	var q = Query{
		SQL:  sql + where,
		Args: args,
	}
	return q
}

// Create generates Query for an INSERT INTO operation
func Create(named repository.Named, vals *values.Values) Query {
	var name = named.Name()

	if name == "" || vals.IsEmpty() {
		return Query{}
	}

	var (
		n = vals.Length()

		q = Query{
			Args: make([]interface{}, n),
		}

		m = make([]string, n)
		p = make([]string, n)
	)

	var it = vals.Iterator()
	for i := 0; it.Next(); i++ {
		v := it.Value()
		m[i] = "`" + v.Name + "`"
		p[i] = "?"
		q.Args[i] = v.Value
	}

	q.SQL = "INSERT INTO `" + name + "` (" + strings.Join(m, ",") + ") VALUES (" + strings.Join(p, ",") + ")"

	return q
}

// ValuesToSet accepts 1 or more values and returns (SET field1 = ?, field2 = ?) and args
func ValuesToSet(vals *values.Values) (set string, args []interface{}) {
	if vals.IsEmpty() {
		return "", nil
	}

	var (
		n  = vals.Length()
		s  = make([]string, n)
		it = vals.Iterator()
	)

	args = make([]interface{}, n)

	for i := 0; it.Next(); i++ {
		c := it.Value()
		s[i] = "`" + c.Name + "` = ?"
		args[i] = c.Value
	}

	return "SET " + strings.Join(s, ", "), args
}

// Update generates Query for an UPDATE ... WHERE id = ? query
func Update(named repository.Named, id string, vals *values.Values) Query {
	var name = named.Name()

	if name == "" || id == "" || vals.IsEmpty() {
		return Query{}
	}

	var set, args = ValuesToSet(vals)

	return Query{
		SQL:  "UPDATE `" + name + "` " + set + " WHERE `id` = ?",
		Args: append(args, id),
	}
}

// UpdateWhere generates Query for an UPDATE ... WHERE ... query
func UpdateWhere(named repository.Named, vals *values.Values, c ...repository.Condition) Query {
	var name = named.Name()

	if name == "" || vals.IsEmpty() {
		return Query{}
	}

	var set, args = ValuesToSet(vals)
	var where, argw = ConditionsToWhere(c)

	return Query{
		SQL:  "UPDATE `" + name + "` " + set + where,
		Args: append(args, argw...),
	}
}

// GetScanDest returns a slice of memory locations appropriate for scanning values row by row
func GetScanDest(f fields.Fields) []interface{} {
	var (
		it  = f.Iterator()
		dst = make([]interface{}, f.Length())
	)

	for i := 0; it.Next(); i++ {
		dst[i] = types.New(it.Field().Kind)
	}

	return dst
}
