package token

// Fixity Associativity
type Fixity int

//goland:noinspection GoSnakeCaseUsage,GoUnusedConst
const (
	NA Fixity = iota
	PREFIX
	INFIX_N
	INFIX_L
	INFIX_R
	POSTFIX
)
