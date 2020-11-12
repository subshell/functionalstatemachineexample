package functionalstatemachineexample

func isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func makeToken(typ string, literal string) Token {
	return Token{
		Type:    typ,
		Literal: literal,
	}
}
