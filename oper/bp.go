package oper

// BP BindingPower, Precedence
// 这里使用 float 是因为可以更精细定义自定义操作符的优先级
// e.g. 如果需要区分前后缀操作符优先级, 可以自己调整
type BP float32

//goland:noinspection GoSnakeCaseUsage
const (
	BP_NONE       BP = iota
	BP_LEFT_BRACE    // {
	BP_COND          // ?:
	BP_LOGIC_OR      // ||
	BP_LOGIC_AND     // &&
	BP_EQ            // == !=
	BP_CMP           // < > <= >=
	BP_TERM          // + -
	BP_FACTOR        // * / %
	BP_EXP           // ^
	BP_PREFIX        // - !
	BP_POSTFIX
	BP_CALL   // ()
	BP_MEMBER // . []
)
