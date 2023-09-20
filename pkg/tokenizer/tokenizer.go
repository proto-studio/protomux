// Package tokenizer provides utilities for tokenization of paths and hosts, including token types and their string representations.
package tokenizer

// Tokenizer represents an interface for tokenizing byte sequences into individual tokens.
type Tokenizer interface {
	Next() (Token, TokenType, error)
}
