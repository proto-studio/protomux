package tokenizers_test

import (
	"bytes"
	"testing"

	"proto.zip/studio/mux/internal/tokenizers"
	"proto.zip/studio/mux/pkg/tokenizer"
)

func TestDomainPatternTokenizer(t *testing.T) {
	Domain := []byte("this.is.a.test")

	tok := tokenizers.NewDomainPatternTokenizer(Domain)

	if err := expectNextToken("first token", []byte("test"), tokenizer.TokenTypeLiteral, tok); err != nil {
		t.Error(err)
	}

	if err := expectNextToken("second token", []byte("a"), tokenizer.TokenTypeLiteral, tok); err != nil {
		t.Error(err)
	}

	if err := expectNextToken("third token", []byte("is"), tokenizer.TokenTypeLiteral, tok); err != nil {
		t.Error(err)
	}

	if err := expectNextToken("fourth token", []byte("this"), tokenizer.TokenTypeLiteral, tok); err != nil {
		t.Error(err)
	}

	if err := expectNextToken("last token", nil, tokenizer.TokenTypeNil, tok); err != nil {
		t.Error(err)
	}
}

func TestDomainPatternTokenizerDoubleDot(t *testing.T) {
	path := []byte("some..test")
	tok := tokenizers.NewDomainPatternTokenizer(path)

	if err := expectNextToken("first token", []byte("test"), tokenizer.TokenTypeLiteral, tok); err != nil {
		t.Error(err)
	}

	token, tokType, err := tok.Next()

	tokenizerErr, ok := err.(*tokenizer.TokenizerError)

	if token != nil {
		t.Errorf("Expected token to be nil, got '%s'", token)
		return
	}
	if tokType != tokenizer.TokenTypeNil {
		t.Errorf("Expected token type to be be ``%s``, got '%s'", tokenizer.TokenTypeNil, tokType)
		return
	}
	if !ok {
		t.Errorf("Expected error to be a TokenizerError, got: %v", err)
		return
	}

	if tokenizerErr.Character != '.' {
		t.Errorf("Expected unexpected character to be '/', got '%c'", tokenizerErr.Character)
	}

	expectedPos := bytes.IndexByte(path, '.')

	if tokenizerErr.Pos != expectedPos {
		t.Errorf("Expected unexpected position to be %d, got %d", expectedPos, tokenizerErr.Pos)
	}
}

func TestDomainPatternTokenizerStartingDot(t *testing.T) {
	path := []byte(".test")
	tok := tokenizers.NewDomainPatternTokenizer(path)

	if err := expectNextToken("first token", []byte("test"), tokenizer.TokenTypeLiteral, tok); err != nil {
		t.Error(err)
	}

	token, tokType, err := tok.Next()

	tokenizerErr, ok := err.(*tokenizer.TokenizerError)

	if token != nil {
		t.Errorf("Expected token to be nil, got '%s'", token)
		return
	}
	if tokType != tokenizer.TokenTypeNil {
		t.Errorf("Expected token type to be be ``%s``, got '%s'", tokenizer.TokenTypeNil, tokType)
		return
	}
	if !ok {
		t.Errorf("Expected error to be a TokenizerError, got: %v", err)
		return
	}

	if tokenizerErr.Character != 0 {
		t.Errorf("Expected unexpected character to be 0, got '%c'", tokenizerErr.Character)
	}

	if tokenizerErr.Pos != -1 {
		t.Errorf("Expected unexpected position to be -1, got %d", tokenizerErr.Pos)
	}
}

func TestDomainPatternTokenizerEndingDot(t *testing.T) {
	path := []byte("test.")
	tok := tokenizers.NewDomainPatternTokenizer(path)

	token, tokType, err := tok.Next()

	tokenizerErr, ok := err.(*tokenizer.TokenizerError)

	if token != nil {
		t.Errorf("Expected token to be nil, got '%s'", token)
		return
	}
	if tokType != tokenizer.TokenTypeNil {
		t.Errorf("Expected token type to be be ``%s``, got '%s'", tokenizer.TokenTypeNil, tokType)
		return
	}
	if !ok {
		t.Errorf("Expected error to be a TokenizerError, got: %v", err)
		return
	}

	if tokenizerErr.Character != '.' {
		t.Errorf("Expected unexpected character to be '.', got '%c'", tokenizerErr.Character)
	}

	expectedPos := len(path) - 1
	if tokenizerErr.Pos != expectedPos {
		t.Errorf("Expected unexpected position to be %d, got %d", expectedPos, tokenizerErr.Pos)
	}
}

var longDomainPattern []byte
var shortDomainPattern []byte = []byte("this.{is}.a.{domain}.for.{benchmarking}")

func BenchmarkDomainPatternTokenizer_6(b *testing.B) {
	for n := 0; n < b.N; n++ {
		tok := tokenizers.NewDomainPatternTokenizer(shortDomainPattern)
		for {
			n, _, err := tok.Next()
			if err != nil {
				b.Errorf("Error running benchmark: %s", err)
				return
			}
			if n != nil {
				break
			}
		}
	}
}

func BenchmarkDomainPatternTokenizer_600(b *testing.B) {
	// Long Domain is 100x longer than the short Domain
	// The resulting string is 14x longer than the maximum safe URL length of 253 characters
	if longDomainPattern == nil {
		substrLen := len(shortDomainPattern)
		strLen := substrLen * 100
		longDomainPattern = make([]byte, 0, strLen)
		for i := 0; i < strLen; i += substrLen {
			longDomainPattern = append(longDomainPattern, shortDomainPattern...)
		}
		b.ResetTimer()
	}

	for n := 0; n < b.N; n++ {
		tok := tokenizers.NewDomainPatternTokenizer(longDomainPattern)
		for {
			n, _, err := tok.Next()
			if err != nil {
				b.Errorf("Error running benchmark: %s", err)
				return
			}
			if n != nil {
				break
			}
		}
	}
}
