package tokenizers

import (
	"proto.zip/studio/mux/pkg/tokenizer"
)

// PathTokenizer is responsible for tokenizing paths.
// Unlike PathPatternTokenizer, PathTokenizer does not allow expressions in the path.
type PathTokenizer struct {
	path []byte
	len  int
	pos  int
}

// NewPathTokenizer initializes a new PathTokenizer with the given path.
func NewPathTokenizer(path []byte) *PathTokenizer {
	t := &PathTokenizer{
		path: path,
		len:  len(path),
	}
	return t
}

// Next returns the next token from the path.
// It processes the path from left to right, splitting it at slashes.
// If an error occurs during tokenization, such as encountering unexpected characters or double slashes, it returns a TokenizerError.
func (t *PathTokenizer) Next() (tokenizer.Token, tokenizer.TokenType, error) {
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

// TrailingSlash checks if the path ends with a slash.
func (t *PathTokenizer) TrailingSlash() bool {
	return t.len > 0 && t.path[t.len-1] == '/'
}
