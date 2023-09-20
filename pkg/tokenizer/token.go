package tokenizer

// Token represents a token as a slice of bytes.
type Token []byte

// String returns the string representation of the Token.
func (t Token) String() string {
	return string(t)
}
