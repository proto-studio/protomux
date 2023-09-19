package tokenizer

type TokenType int

const (
	TokenTypeNil TokenType = iota
	TokenTypeLiteral
	TokenTypeLabel
	TokenTypeWildcard
	TokenTypeUserDefined TokenType = 10
)

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
