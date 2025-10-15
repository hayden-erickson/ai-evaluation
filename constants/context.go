package constants

// Context keys for storing values in context
type ContextKey string

const (
	ClaimsKey ContextKey = "claims"
	BankKey   ContextKey = "db"
)
