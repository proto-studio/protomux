package tokenizers

import "proto.zip/studio/mux/pkg/tokenizer"

type DomainPatternTokenizer struct {
	domain []byte
	pos    int
}

func NewDomainPatternTokenizer(domain []byte) *DomainPatternTokenizer {
	t := &DomainPatternTokenizer{
		domain: domain,
		pos:    len(domain) - 1,
	}
	return t
}

func (t *DomainPatternTokenizer) Next() (tokenizer.Token, tokenizer.TokenType, error) {
	if t.pos == -1 {
		return nil, tokenizer.TokenTypeNil, nil
	}

	// Tokens must start with a dot '.' except the first one, which must never start with a dot
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

	// Variables have the format { label }
	// We're going backwards so we look for the end first
	if t.pos >= 0 && t.domain[t.pos] == '}' {
		t.pos--

		// Eat all trailing whitespace
		for t.pos >= 0 && t.domain[t.pos] == ' ' {
			t.pos--
		}

		start := t.pos + 1

		// Eat label until we hit a bracket or space
		for t.pos >= 0 && t.domain[t.pos] != '{' && t.domain[t.pos] != ' ' {
			t.pos--
		}

		// We hit the end of the label part so store it now
		ret := t.domain[t.pos : start+1]

		// Eat leading whitespace
		for t.pos >= 0 && t.domain[t.pos] == ' ' {
			t.pos--
		}

		// We're past the end (start)
		if t.pos == -1 {
			return nil, tokenizer.TokenTypeNil, &tokenizer.TokenizerError{
				Pos: t.pos,
			}
		}

		// Make sure we're ending on a '{'
		if t.domain[t.pos] != '{' {
			return nil, tokenizer.TokenTypeNil, &tokenizer.TokenizerError{
				Pos:       t.pos,
				Character: rune(t.domain[t.pos]),
			}
		}

		t.pos--
		return ret, tokenizer.TokenTypeLabel, nil

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
