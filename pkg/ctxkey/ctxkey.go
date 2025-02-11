// In a separate package (e.g., pkg/ctxkey/keys.go)
package ctxkey

type contextKey struct {
	name string
}

// UserIDKey is used to store/retrieve user ID from context
var UserIDKey = &contextKey{"user_id"}
