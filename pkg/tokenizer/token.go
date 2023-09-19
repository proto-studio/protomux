package tokenizer

type Token []byte

func (t Token) String() string {
	return string(t)
}
