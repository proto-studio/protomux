package tokenizer

type Tokenizer interface {
	Next() (Token, TokenType, error)
}
