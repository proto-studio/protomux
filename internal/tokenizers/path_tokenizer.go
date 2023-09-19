package tokenizers

import (
	"proto.zip/studio/mux/pkg/tokenizer"
)

type PathTokenizer struct {
	path []byte
	len  int
	pos  int
}

func NewPathTokenizer(path []byte) *PathTokenizer {
	t := &PathTokenizer{
		path: path,
		len:  len(path),
	}
	return t
}

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

func (t *PathTokenizer) TrailingSlash() bool {
	return t.len > 0 && t.path[t.len-1] == '/'
}
