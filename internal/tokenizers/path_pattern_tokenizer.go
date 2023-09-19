package tokenizers

import (
	"proto.zip/studio/mux/pkg/tokenizer"
)

type PathPatternTokenizer struct {
	path []byte
	len  int
	pos  int
}

func NewPathPatternTokenizer(path []byte) *PathPatternTokenizer {
	t := &PathPatternTokenizer{
		path: path,
		len:  len(path),
	}
	return t
}

func (t *PathPatternTokenizer) Next() (tokenizer.Token, tokenizer.TokenType, error) {
	if t.pos == t.len {
		return nil, tokenizer.TokenTypeNil, nil
	}

	// Initial character must be a slash except at the start of the slice
	if t.path[t.pos] == '/' {
		t.pos++
	} else if t.pos != 0 {
		return nil, tokenizer.TokenTypeNil, &tokenizer.TokenizerError{
			Pos:       t.pos,
			Character: rune(t.path[t.pos]),
		}
	}

	// Variables have the format { label }
	if t.pos < t.len && t.path[t.pos] == '{' {
		t.pos++

		// Eat all leading whitespace
		for t.pos < t.len && t.path[t.pos] == ' ' {
			t.pos++
		}

		start := t.pos

		// Eat label until we hit a bracket or space
		for t.pos < t.len && t.path[t.pos] != '}' && t.path[t.pos] != ' ' {
			t.pos++
		}

		// We hit the end of the label part so store it now
		ret := t.path[start:t.pos]

		// Eat trailing whitespace
		for t.pos < t.len && t.path[t.pos] == ' ' {
			t.pos++
		}

		// We're past the end
		if t.pos == t.len {
			return nil, tokenizer.TokenTypeNil, &tokenizer.TokenizerError{
				Pos: t.pos,
			}
		}

		// Make sure we're ending on a '}'
		if t.path[t.pos] != '}' {
			return nil, tokenizer.TokenTypeNil, &tokenizer.TokenizerError{
				Pos:       t.pos,
				Character: rune(t.path[t.pos]),
			}
		}

		t.pos++
		return ret, tokenizer.TokenTypeLabel, nil
	}

	// Not a variable, must be a litteral
	// Read until we hit a slash
	start := t.pos

	for t.pos < t.len && t.path[t.pos] != '/' {
		t.pos++
	}

	// We didn't progress, we're either at the end or had a double slash
	if t.pos == start {
		if t.pos == t.len {
			return nil, tokenizer.TokenTypeNil, nil
		} else {
			return nil, tokenizer.TokenTypeNil, &tokenizer.TokenizerError{
				Pos:       t.pos,
				Character: rune(t.path[t.pos]),
			}
		}
	}

	ret := t.path[start:t.pos]
	return ret, tokenizer.TokenTypeLiteral, nil
}

func (t *PathPatternTokenizer) TrailingSlash() bool {
	return t.len > 0 && t.path[t.len-1] == '/'
}
