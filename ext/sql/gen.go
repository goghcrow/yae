// Generated Code; DO NOT EDIT.

package sql

import "github.com/goghcrow/yae/val"

//go:generate /bin/bash gen.sh

var funs = []*val.Val{
	BETWEEN_NUM_NUM_NUM,
	BETWEEN_TIME_TIME_TIME,
	EQ_BOOL_BOOL,
	EQ_NUM_NUM,
	EQ_STR_STR,
	EQ_TIME_TIME,
	GE_NUM_NUM,
	GE_TIME_TIME,
	GT_NUM_NUM,
	GT_TIME_TIME,
	IN_LIST,
	IS_NULL_A,
	LE_NUM_NUM,
	LE_TIME_TIME,
	LIKE_STR_STR,
	LOGIC_AND_BOOL_BOOL,
	LOGIC_NOT_BOOL,
	LOGIC_OR_BOOL_BOOL,
	LT_NUM_NUM,
	LT_TIME_TIME,
	NE_BOOL_BOOL,
	NE_NUM_NUM,
	NE_STR_STR,
	NE_TIME_TIME,
}

func BuiltIn() []*val.Val {
	return funs
}

