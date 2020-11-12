package functionalstatemachineexample

import (
	"io"
	"strings"
	"testing"

	"github.com/matryer/is"
)

var (
	tests = []struct {
		input    string
		expected []Token
	}{
		{
			`1 + 23 - 345 "56 + 7" "" 3 - 2`,
			makeTokens(
				Int, "1",
				Plus, "+",
				Int, "23",
				Minus, "-",
				Int, "345",
				String, "56 + 7",
				String, "",
				Int, "3",
				Minus, "-",
				Int, "2",
			),
		},
	}
)

func TestTokensTraditional(t *testing.T) {
	testTokens(TokensTraditional, t)
}

func TestTokensFunctional(t *testing.T) {
	testTokens(TokensFunctional, t)
}

func testTokens(lexerFunc func(r io.Reader, ch chan<- Token), t *testing.T) {
	t.Helper()

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			is := is.New(t)

			r := strings.NewReader(test.input)
			ch := make(chan Token, 100)
			lexerFunc(r, ch)
			close(ch)

			tokens := []Token{}
			for t := range ch {
				tokens = append(tokens, t)
			}

			is.Equal(tokens, test.expected)
		})
	}
}

func makeTokens(v ...string) []Token {
	tokens := make([]Token, len(v)/2)
	for i := 0; i < len(v)/2; i++ {
		tokens[i] = makeToken(v[i*2], v[i*2+1])
	}
	return tokens
}
