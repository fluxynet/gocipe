package tenant

import (
	"context"
	"net/http"
)

// Key name used in context
const Key = "tenant"

type tenant string

// Get returns a tenant name from context
func Get(ctx context.Context) string {
	var v = ctx.Value(tenant(Key))

	if v == nil {
		return ""
	}

	var name, ok = v.(tenant)
	if ok {
		return string(name)
	}

	return ""
}

// Set assigns a tenant in context
func Set(ctx context.Context, name string) context.Context {
	return context.WithValue(ctx, tenant(Key), name)
}

// StaticHttp middleware to set tenant as a static value
func StaticHttp(name string) func(next http.Handler) http.Handler {
	return staticHttp{Name: name}.Middleware
}

// staticHttp provides a middleware to set tenant as a static value
type staticHttp struct {
	Name string
}

func (s staticHttp) Middleware(next http.Handler) http.Handler {
	var fn = func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(Set(r.Context(), s.Name))
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
