package lexer

import (
	"github.com/goghcrow/yae/oper"
	"github.com/goghcrow/yae/token"
)

var keywords = []token.Type{
	// if 是普通函数不是 keyword
	//token.IF,
	//token.THEN,
	//token.ELSE,
	//token.END,
}

var buildInOpers = []token.Type{
	token.DOT,      //.
	token.QUESTION, //?:
}

func newLexicon(ops []oper.Operator) lexicon {
	l := lexicon{}

	l.addRule(str(token.COLON)) // :
	l.addRule(str(token.COMMA)) // ,

	l.addRule(str(token.LEFT_PAREN))    // (
	l.addRule(str(token.RIGHT_PAREN))   // )
	l.addRule(str(token.LEFT_BRACKET))  // [
	l.addRule(str(token.RIGHT_BRACKET)) // ]
	l.addRule(str(token.LEFT_BRACE))    // {
	l.addRule(str(token.RIGHT_BRACE))   // }

	for _, kw := range keywords {
		l.addRule(keyword(kw))
	}

	// 内置的操作符优先级高于自定义操作符
	for _, op := range buildInOpers {
		l.addRule(primOper(op))
	}

	// 自定义操作符
	for _, op := range ops {
		l.addOper(op.Type)
	}

	l.addRule(str(token.TRUE))  // true
	l.addRule(str(token.FALSE)) // false

	// 移除数字前的 [+-]?, [+-]? 被处理成一元操作符, 实际上变成没有负数字面量, 语义不变
	l.addRule(reg(token.NUM, "(?:0|[1-9][0-9]*)(?:[.][0-9]+)+(?:[eE][-+]?[0-9]+)?")) // float
	l.addRule(reg(token.NUM, "(?:0|[1-9][0-9]*)(?:[.][0-9]+)?(?:[eE][-+]?[0-9]+)+")) // float
	l.addRule(reg(token.NUM, "0b(?:0|1[0-1]*)"))                                     // int
	l.addRule(reg(token.NUM, "0x(?:0|[1-9a-fA-F][0-9a-fA-F]*)"))                     // int
	l.addRule(reg(token.NUM, "0o(?:0|[1-7][0-7]*)"))                                 // int
	l.addRule(reg(token.NUM, "(?:0|[1-9][0-9]*)"))                                   // int

	l.addRule(reg(token.STR, "\"(?:[^\"\\\\]*|\\\\[\"\\\\trnbf\\/]|\\\\u[0-9a-fA-F]{4})*\""))
	l.addRule(reg(token.STR, "`[^`]*`")) // raw string

	l.addRule(reg(token.TIME, "'[^`\"']*'"))

	l.addRule(reg(token.NAME, "[a-zA-Z\\p{L}_][a-zA-Z0-9\\p{L}_]*")) // 支持 unicode, 不能以数字开头

	return l
}
