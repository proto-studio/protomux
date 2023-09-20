package tokenizer

// TokenType represents the type of a token.
// It can be one of the predefined types or a user-defined type.
type TokenType int

const (
	TokenTypeNil         TokenType = iota // TokenTypeNil represents a nil or undefined token type.
	TokenTypeLiteral                      // TokenTypeLiteral represents a literal token.
	TokenTypeLabel                        // TokenTypeLabel represents a label token.
	TokenTypeWildcard                     // TokenTypeWildcard represents a wildcard token.
	TokenTypeUserDefined TokenType = 10   // TokenTypeUserDefined is the starting point for user-defined token types.
)

// String returns the string representation of the TokenType.
// For user-defined types greater than TokenTypeUserDefined, it returns "USER".
func (tt TokenType) String() string {
	switch tt {
	case TokenTypeNil:
		return "NIL"
	case TokenTypeLiteral:
		return "LITERAL"
	case TokenTypeLabel:
		return "LABEL"
	case TokenTypeWildcard:
		return "WILDCARD"
	default:
		if tt > TokenTypeUserDefined {
			return "USER"
		} else {
			return "UNKNOWN"
		}
	}
}
