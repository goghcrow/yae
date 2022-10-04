package oper

// BP BindingPower, Precedence
type BP float32

//goland:noinspection GoSnakeCaseUsage
const (
	BP_NONE       BP = iota
	BP_LEFT_BRACE    // {
	BP_COND          // ?:
	BP_LOGIC_OR      // ||
	BP_LOGIC_AND     // &&
	BP_EQ            // == !=
	BP_COMP          // < > <= >=
	BP_TERM          // + -
	BP_FACTOR        // * / %
	BP_EXP           // ^
	BP_PREFIX        // - !
	BP_POSTFIX
	BP_CALL   // ()
	BP_MEMBER // . []
)
