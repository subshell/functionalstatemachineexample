package functionalstatemachineexample

import (
	"bufio"
	"errors"
	"io"
	"strings"
)

func TokensTraditional(r io.Reader, ch chan<- Token) {
	inInt := false
	intBuf := strings.Builder{}

	inString := false
	stringBuf := strings.Builder{}

	bufR := bufio.NewReader(r)
	for {
		c, _, err := bufR.ReadRune()
		if err != nil {
			if !errors.Is(err, io.EOF) {
				panic(err)
			}

			if inInt {
				ch <- makeToken(Int, intBuf.String())
			}

			if inString {
				panic(errors.New("string not closed")) //nolint:go-lint
			}

			break
		}

		if inInt {
			if isDigit(c) {
				intBuf.WriteRune(c)
				continue
			}

			ch <- makeToken(Int, intBuf.String())
			inInt = false
			intBuf.Reset()
		}

		if inString {
			if c != '"' {
				stringBuf.WriteRune(c)
				continue
			}

			ch <- makeToken(String, stringBuf.String())
			inString = false
			stringBuf.Reset()
			continue
		}

		if isDigit(c) {
			inInt = true
			intBuf.WriteRune(c)
			continue
		}

		switch c {
		case '+':
			ch <- makeToken(Plus, "+")
		case '-':
			ch <- makeToken(Minus, "-")
		case '"':
			inString = true
		}
	}
}
