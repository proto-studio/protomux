package tokenizer

import "fmt"

type TokenizerError struct {
	Character rune
	Pos       int
}

func (e *TokenizerError) Error() string {
	if e.Character == 0 {
		return fmt.Sprintf("unexpedted end of string at %d", e.Pos)
	} else {
		return fmt.Sprintf("unexpedted character '%c' at %d", e.Character, e.Pos)
	}
}
