package functionalstatemachineexample

import (
	"io"
	"strings"
	"testing"
)

var Tok Token

const input = `1 + 2 - 3 "5 + 7" "" 3 - 2`

func BenchmarkTokensTraditional(b *testing.B) {
	benchmarkTokens(TokensTraditional, b)
}

func BenchmarkTokensFunctional(b *testing.B) {
	benchmarkTokens(TokensFunctional, b)
}

func benchmarkTokens(lexerFunc func(r io.Reader, ch chan<- Token), b *testing.B) {
	b.Helper()
	b.StopTimer()

	for i := 0; i < b.N; i++ {
		r := strings.NewReader(input)
		ch := make(chan Token, 100)

		b.StartTimer()
		lexerFunc(r, ch)
		b.StopTimer()

		close(ch)

		Tok = <-ch
	}
}
