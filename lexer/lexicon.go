package lex

import "github.com/goghcrow/yae/token"

// lexicon Lexical grammar, 按顺序匹配, 匹配到第一个停止
var lexicon = []rule{
	str(token.COLON), // :
	str(token.COMMA), // ,

	str(token.LEFT_PAREN),    // (
	str(token.RIGHT_PAREN),   // )
	str(token.LEFT_BRACKET),  // [
	str(token.RIGHT_BRACKET), // ]
	str(token.LEFT_BRACE),    // {
	str(token.RIGHT_BRACE),   // }

	str(token.DOT),      // .
	str(token.QUESTION), // ?

	// 注释掉该行即可支持 if(bool, x, x) 风格!!!
	//keyword(token.IF),
	keyword(token.THEN),
	keyword(token.ELSE),
	keyword(token.END),

	keywordAlias(token.LOGIC_NOT, token.NOT.Name()),
	keywordAlias(token.LOGIC_AND, token.AND.Name()),
	keywordAlias(token.LOGIC_OR, token.OR.Name()),

	str(token.PLUS),  // +
	str(token.MINUS), // -

	str(token.MUL), // *
	str(token.DIV), // /
	str(token.EXP), // ^
	str(token.MOD), // %

	str(token.LOGIC_OR),  // ||
	str(token.LOGIC_AND), // &&
	str(token.LE),        // <=
	str(token.LT),        // <
	str(token.GE),        // >=
	str(token.GT),        // >
	str(token.EQ),        // ==
	str(token.NE),        // !=
	str(token.LOGIC_NOT), // !

	str(token.NULL),  // null
	str(token.TRUE),  // true
	str(token.FALSE), // false

	// 移除数字前的 [+-]?, lex没有使用最长路径来匹配, +- 被优先匹配成操作符了
	// 如果优先匹配数字的话, 1-1, 会被分成 1,-1, 需要修一遍 tokseq
	// so, 不支持 + 数字, - 数字 处理成一元操作符
	//reg(token.NUM, "(?:0|[1-9][0-9]*)(?:[.][0-9]+)?"),
	reg(token.NUM, "(?:0|[1-9][0-9]*)(?:[.][0-9]+)+(?:[eE][-+]?[0-9]+)?"), // float
	reg(token.NUM, "(?:0|[1-9][0-9]*)(?:[.][0-9]+)?(?:[eE][-+]?[0-9]+)+"), // float
	reg(token.NUM, "0b(?:0|1[0-1]*)"),                                     // int
	reg(token.NUM, "0x(?:0|[1-9a-fA-F][0-9a-fA-F]*)"),                     // int
	reg(token.NUM, "0o(?:0|[1-7][0-7]*)"),                                 // int
	reg(token.NUM, "(?:0|[1-9][0-9]*)"),                                   // int

	reg(token.STR, "\"(?:[^\"\\\\]*|\\\\[\"\\\\trnbf\\/]|\\\\u[0-9a-fA-F]{4})*\""),
	reg(token.STR, "`[^`]*`"), // raw string

	reg(token.TIME, "'[^`\"']*'"),

	reg(token.NAME, "[a-zA-Z\\p{L}_][a-zA-Z0-9\\p{L}_]*"), // 支持 unicode, 不能以数字开头
}
