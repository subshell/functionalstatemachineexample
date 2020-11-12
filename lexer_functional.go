package functionalstatemachineexample

import (
	"bufio"
	"errors"
	"io"
	"strings"
)

type stateFunc func(r *bufio.Reader, ch chan<- Token) stateFunc

type runeFunc func(c rune, eof bool, stop func(next stateFunc))

func TokensFunctional(r io.Reader, ch chan<- Token) {
	bufR := bufio.NewReader(r)

	state := parse
	for state != nil {
		state = state(bufR, ch)
	}
}

func parse(r *bufio.Reader, ch chan<- Token) stateFunc {
	c, err := peekRune(r)
	if err != nil {
		if !errors.Is(err, io.EOF) {
			return parseError(err)
		}

		return nil
	}

	if isDigit(c) {
		return parseInt
	}

	switch c {
	case '+':
		return parseOp(Plus, "+")
	case '-':
		return parseOp(Minus, "-")
	case '"':
		return parseString
	}

	return skip
}

func parseInt(r *bufio.Reader, ch chan<- Token) stateFunc {
	buf := strings.Builder{}

	defer emitBuffer(ch, Int, &buf)

	parseRune := func(c rune, eof bool, stop func(next stateFunc)) {
		if eof {
			stop(nil)
			return
		}

		if !isDigit(c) {
			_ = r.UnreadRune()
			stop(parse)
			return
		}

		buf.WriteRune(c)
	}

	return forRunes(r, parseRune)
}

func parseOp(typ string, literal string) stateFunc {
	return func(r *bufio.Reader, ch chan<- Token) stateFunc {
		if _, _, err := r.ReadRune(); err != nil {
			return parseError(err)
		}

		ch <- makeToken(typ, literal)

		return parse
	}
}

func parseString(r *bufio.Reader, ch chan<- Token) stateFunc {
	if _, _, err := r.ReadRune(); err != nil {
		return parseError(err)
	}

	buf := strings.Builder{}

	defer emitBuffer(ch, String, &buf)

	parseRune := func(c rune, eof bool, stop func(next stateFunc)) {
		if eof {
			stop(parseError(errors.New("string not closed"))) //nolint:go-lint
			return
		}

		if c == '"' {
			stop(parse)
			return
		}

		buf.WriteRune(c)
	}

	return forRunes(r, parseRune)
}

func skip(r *bufio.Reader, ch chan<- Token) stateFunc {
	if _, _, err := r.ReadRune(); err != nil {
		return parseError(err)
	}
	return parse
}

func parseError(err error) stateFunc {
	return func(r *bufio.Reader, ch chan<- Token) stateFunc {
		panic(err)
	}
}

func forRunes(r *bufio.Reader, f runeFunc) stateFunc { //nolint:go-lint
	running := true
	nextState := stateFunc(nil)

	stop := func(next stateFunc) {
		running = false
		nextState = next
	}

	for running {
		c, _, err := r.ReadRune()
		if err != nil {
			if !errors.Is(err, io.EOF) {
				nextState = parseError(err)
			} else {
				f(0, true, stop)
			}
			break
		}

		f(c, false, stop)
	}

	return nextState
}

func emitBuffer(ch chan<- Token, typ string, buf *strings.Builder) {
	ch <- makeToken(typ, buf.String())
}

func peekRune(r *bufio.Reader) (rune, error) { //nolint:go-lint
	c, _, err := r.ReadRune()
	if err != nil {
		return 0, err
	}

	if err := r.UnreadRune(); err != nil {
		panic(err)
	}

	return c, nil
}
