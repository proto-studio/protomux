package tokenizers_test

import (
	"bytes"
	"testing"

	"proto.zip/studio/mux/internal/tokenizers"
	"proto.zip/studio/mux/pkg/tokenizer"
)

func TestPathTokenizer(t *testing.T) {
	path := []byte("this/is/a/test")

	tok := tokenizers.NewPathTokenizer(path)

	if err := expectNextToken("first token", []byte("this"), tokenizer.TokenTypeLiteral, tok); err != nil {
		t.Error(err)
	}

	if err := expectNextToken("second token", []byte("is"), tokenizer.TokenTypeLiteral, tok); err != nil {
		t.Error(err)
	}

	if err := expectNextToken("third token", []byte("a"), tokenizer.TokenTypeLiteral, tok); err != nil {
		t.Error(err)
	}

	if err := expectNextToken("fourth token", []byte("test"), tokenizer.TokenTypeLiteral, tok); err != nil {
		t.Error(err)
	}

	if err := expectNextToken("last token", nil, tokenizer.TokenTypeNil, tok); err != nil {
		t.Error(err)
	}
}

func TestPathTokenizerLeadingSlash(t *testing.T) {
	path := []byte("/test")
	tok := tokenizers.NewPathTokenizer(path)

	if err := expectNextToken("first token", []byte("test"), tokenizer.TokenTypeLiteral, tok); err != nil {
		t.Error(err)
	}
}

func TestPathTokenizerTrailingSlash(t *testing.T) {
	path := []byte("some/test")
	tok := tokenizers.NewPathTokenizer(path)

	if tok.TrailingSlash() {
		t.Error("Expected trailing slash to be false")
	}

	path = []byte("some/test/")
	tok = tokenizers.NewPathTokenizer(path)

	if !tok.TrailingSlash() {
		t.Error("Expected trailing slash to be true")
	}

	if err := expectNextToken("first token", []byte("some"), tokenizer.TokenTypeLiteral, tok); err != nil {
		t.Error(err)
	}

	if err := expectNextToken("second token", []byte("test"), tokenizer.TokenTypeLiteral, tok); err != nil {
		t.Error(err)
	}

	if err := expectNextToken("last token", nil, tokenizer.TokenTypeNil, tok); err != nil {
		t.Error(err)
	}
}

func TestPathTokenizerDoubleSlash(t *testing.T) {
	path := []byte("some//test")
	tok := tokenizers.NewPathTokenizer(path)

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

var longPath []byte
var shortPath []byte = []byte("this/is/a/path/for/benchmarking/")

func BenchmarkPathTokenizer_6(b *testing.B) {
	for n := 0; n < b.N; n++ {
		tok := tokenizers.NewPathTokenizer(shortPath)
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

func BenchmarkPathTokenizer_600(b *testing.B) {
	// Long path is 100x longer than the short path
	// The resulting string is 60% longer than the maximum safe URL length of approximately 2000 characters
	if longPath == nil {
		substrLen := len(shortPath)
		strLen := substrLen * 100
		longPath = make([]byte, 0, strLen)
		for i := 0; i < strLen; i += substrLen {
			longPath = append(longPath, shortPath...)
		}
		b.ResetTimer()
	}

	for n := 0; n < b.N; n++ {
		tok := tokenizers.NewPathTokenizer(longPath)
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
