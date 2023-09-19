package tokenizers

import "proto.zip/studio/mux/pkg/tokenizer"

type DomainTokenizer struct {
	domain []byte
	pos    int
}

func NewDomainTokenizer(domain []byte) *DomainTokenizer {
	t := &DomainTokenizer{
		domain: domain,
		pos:    len(domain) - 1,
	}
	return t
}

func (t *DomainTokenizer) Next() (tokenizer.Token, tokenizer.TokenType, error) {
	if t.pos == -1 {
		return nil, tokenizer.TokenTypeNil, nil
	}

	// Tokens must start with a dot '.' except the first one
	if t.pos == len(t.domain)-1 {
		if t.domain[t.pos] == '.' {
			return nil, tokenizer.TokenTypeNil, &tokenizer.TokenizerError{
				Pos:       t.pos,
				Character: rune(t.domain[t.pos]),
			}
		}
	} else if t.domain[t.pos] != '.' {
		return nil, tokenizer.TokenTypeNil, &tokenizer.TokenizerError{
			Pos:       t.pos,
			Character: rune(t.domain[t.pos]),
		}
	} else {
		t.pos--
	}

	start := t.pos

	for t.pos >= 0 && t.domain[t.pos] != '.' {
		t.pos--
	}

	if t.pos == start {
		if t.pos == -1 {
			// Domain starts with a '.'
			return nil, tokenizer.TokenTypeNil, &tokenizer.TokenizerError{
				Pos: t.pos,
			}
		} else {
			// Domain has a double dot '..'
			return nil, tokenizer.TokenTypeNil, &tokenizer.TokenizerError{
				Pos:       t.pos,
				Character: rune(t.domain[t.pos]),
			}
		}
	}

	ret := t.domain[t.pos+1 : start+1]

	return ret, tokenizer.TokenTypeLiteral, nil
}
