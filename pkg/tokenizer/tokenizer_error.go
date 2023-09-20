package tokenizer

import "fmt"

// TokenizerError represents an error encountered during tokenization.
// It contains information about the unexpected character and its position.
type TokenizerError struct {
	Character rune
	Pos       int
}

// Error returns a string representation of the TokenizerError.
// If the character is 0, it indicates an unexpected end of string.
func (e *TokenizerError) Error() string {
	if e.Character == 0 {
		return fmt.Sprintf("unexpected end of string at %d", e.Pos)
	} else {
		return fmt.Sprintf("unexpected character '%c' at %d", e.Character, e.Pos)
	}
}
