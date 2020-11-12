package functionalstatemachineexample

const (
	Int    = "int"
	Plus   = "plus"
	Minus  = "minus"
	String = "string"
)

type Token struct {
	Type    string
	Literal string
}
