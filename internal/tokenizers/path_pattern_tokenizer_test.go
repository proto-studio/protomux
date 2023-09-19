package tokenizers_test

import (
	"bytes"
	"testing"

	"proto.zip/studio/mux/internal/tokenizers"
	"proto.zip/studio/mux/pkg/tokenizer"
)

func TestPathPatternTokenizer(t *testing.T) {
	path := []byte("this/{is}/a/{test}")

	tok := tokenizers.NewPathPatternTokenizer(path)

	if err := expectNextToken("first token", []byte("this"), tokenizer.TokenTypeLiteral, tok); err != nil {
		t.Error(err)
	}

	if err := expectNextToken("second token", []byte("is"), tokenizer.TokenTypeLabel, tok); err != nil {
		t.Error(err)
	}

	if err := expectNextToken("third token", []byte("a"), tokenizer.TokenTypeLiteral, tok); err != nil {
		t.Error(err)
	}

	if err := expectNextToken("fourth token", []byte("test"), tokenizer.TokenTypeLabel, tok); err != nil {
		t.Error(err)
	}

	if err := expectNextToken("last token", nil, tokenizer.TokenTypeNil, tok); err != nil {
		t.Error(err)
	}
}

func TestPathPatternTokenizerLeadingSlash(t *testing.T) {
	path := []byte("/test")
	tok := tokenizers.NewPathPatternTokenizer(path)

	if err := expectNextToken("first token", []byte("test"), tokenizer.TokenTypeLiteral, tok); err != nil {
		t.Error(err)
	}
}

func TestPathPatternTokenizerTrailingSlash(t *testing.T) {
	path := []byte("some/{test}")
	tok := tokenizers.NewPathPatternTokenizer(path)

	if tok.TrailingSlash() {
		t.Error("Expected trailing slash to be false")
	}

	path = []byte("some/{test}/")
	tok = tokenizers.NewPathPatternTokenizer(path)

	if !tok.TrailingSlash() {
		t.Error("Expected trailing slash to be true")
	}

	if err := expectNextToken("first token", []byte("some"), tokenizer.TokenTypeLiteral, tok); err != nil {
		t.Error(err)
	}

	if err := expectNextToken("second token", []byte("test"), tokenizer.TokenTypeLabel, tok); err != nil {
		t.Error(err)
	}

	if err := expectNextToken("last token", nil, tokenizer.TokenTypeNil, tok); err != nil {
		t.Error(err)
	}
}

func TestPathPatternTokenizerDoubleSlash(t *testing.T) {
	path := []byte("some//test")
	tok := tokenizers.NewPathPatternTokenizer(path)

	if err := expectNextToken("first token", []byte("some"), tokenizer.TokenTypeLiteral, tok); err != nil {
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

	if tokenizerErr.Character != '/' {
		t.Errorf("Expected unexpected character to be '/', got '%c'", tokenizerErr.Character)
	}

	expectedPos := bytes.LastIndexByte(path, '/')

	if tokenizerErr.Pos != expectedPos {
		t.Errorf("Expected unexpected position to be %d, got %d", expectedPos, tokenizerErr.Pos)
	}
}

func TestPathPatternTokenizerWhitespace(t *testing.T) {
	path := []byte("{    test    }")
	tok := tokenizers.NewPathPatternTokenizer(path)

	if err := expectNextToken("token", []byte("test"), tokenizer.TokenTypeLabel, tok); err != nil {
		t.Error(err)
	}
}

func TestPathPatternTokenizerNoClose(t *testing.T) {
	path := []byte("{    test    ")
	tok := tokenizers.NewPathPatternTokenizer(path)

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

	expectedPos := len(path)

	if tokenizerErr.Pos != expectedPos {
		t.Errorf("Expected unexpected position to be %d, got %d", expectedPos, tokenizerErr.Pos)
	}
}

var longPathPattern []byte
var shortPathPattern []byte = []byte("this/{is}/a/{path}/{for}/benchmarking/")

func BenchmarkPathPatternTokenizer_6(b *testing.B) {
	for n := 0; n < b.N; n++ {
		tok := tokenizers.NewPathPatternTokenizer(shortPathPattern)
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

func BenchmarkPathPatternTokenizer_600(b *testing.B) {
	// Long path is 100x longer than the short path
	// The resulting string is 60% longer than the maximum safe URL length of approximately 2000 characters
	if longPathPattern == nil {
		substrLen := len(shortPathPattern)
		strLen := substrLen * 100
		longPathPattern = make([]byte, 0, strLen)
		for i := 0; i < strLen; i += substrLen {
			longPathPattern = append(longPathPattern, shortPathPattern...)
		}
		b.ResetTimer()
	}

	for n := 0; n < b.N; n++ {
		tok := tokenizers.NewPathPatternTokenizer(longPathPattern)
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
