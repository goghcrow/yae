package ext

import (
	"strings"
	"testing"

	"github.com/goghcrow/yae/conv"
	"github.com/goghcrow/yae/ext/sql"
	"github.com/goghcrow/yae/parser"
	"github.com/goghcrow/yae/parser/ast"
	"github.com/goghcrow/yae/parser/lexer"
	"github.com/goghcrow/yae/parser/oper"
	"github.com/goghcrow/yae/parser/pos"
	"github.com/goghcrow/yae/trans"
	"github.com/goghcrow/yae/util"
	"github.com/goghcrow/yae/val"
)

func TestVar(t *testing.T) {
	type model struct {
		UserType     int    `yae:"user_type"`
		UserID       int    `yae:"user_id"`
		DepartmentID int    `yae:"department_id"`
		UserRole     string `yae:"user_role"`

		CurrentUserID int `yae:"current_user_id"`
	}

	type ctx struct {
		CurrentUserID int `yae:"current_user_id"`
	}

	compiled := compileCriteria(`user_type == 1 && user_id == current_user_id && (department_id.in([1,2,3]) || user_role != "admin")`, model{})

	actual, _ := compiled(ctx{42})
	expected := "`user_type` = 1 AND `user_id` = 42 AND (`department_id` IN (1, 2, 3) OR `user_role` <> \"admin\")"
	if actual != expected {
		t.Errorf("expect %s actula %s", expected, actual)
	}

	actual, _ = compiled(ctx{100})
	expected = "`user_type` = 1 AND `user_id` = 100 AND (`department_id` IN (1, 2, 3) OR `user_role` <> \"admin\")"
	if actual != expected {
		t.Errorf("expect %s actula %s", expected, actual)
	}
}

func TestSimple(t *testing.T) {
	type model struct {
		types string `yae:"biz_type"`
		id    int
	}
	expr := CondGroup{
		LogicalOper: AND,
		Conds: []Criteria{
			Cond{
				Field:    "biz_type",
				Operator: sql.EQ,
				Operands: []ast.Expr{
					ast.Str(`"'xxx"`, pos.Unknown),
				},
			},
			Cond{
				Field:    "id",
				Operator: sql.IN,
				Operands: []ast.Expr{
					ast.List([]ast.Expr{
						ast.Num("1", pos.Unknown),
						ast.Num("2", pos.Unknown),
						ast.Num("3", pos.Unknown),
					}, pos.Unknown),
				},
			},
		},
	}

	of := conv.MustTypeEnvOf(model{})
	f := CompileToSql(expr, of)
	actual, _ := f(val.NewEnv())
	expected := "`biz_type` = \"'xxx\" AND `id` IN (1, 2, 3)"
	if actual != expected {
		t.Errorf("expect %s actula %s", expected, actual)
	}
}

func TestMixed(t *testing.T) {
	criteria := parseCriteria(`(a > 1 or b < 2 and c >=3 and d <= 4 or e == 5 and s.like("hello%")) and f.between(1,100) or not isnull(f)`)
	actual, _ := CompileToSql(criteria, conv.MustTypeEnvOf(struct {
		a, b, c, d, e int
		s             string
		f             int
	}{}))(val.NewEnv())
	expected := "(`a` > 1 OR `b` < 2 AND `c` >= 3 AND `d` <= 4 OR `e` = 5 AND `s` LIKE \"hello%\") AND `f` BETWEEN 1 AND 100 OR NOT `f` IS NULL"
	if actual != expected {
		t.Errorf("expect %s actula %s", expected, actual)
	}
}

func compileCriteria(s string, ty interface{}) func(v interface{}) (string, error) {
	return CompileToSql(parseCriteria(s), conv.MustTypeEnvOf(ty))
}

func parse(s string) ast.Expr {
	ops := oper.BuiltIn()
	toks := lexer.NewLexer(ops).Lex(s)
	return trans.Desugar(parser.NewParser(ops).Parse(toks))
}

func parseCriteria(s string) Criteria { return expr2criteria(parse(s)) }

func expr2criteria(expr ast.Expr) Criteria { return expr2criteria0(expr).(Criteria) }

// for test
func expr2criteria0(expr ast.Expr) interface{} /*ast.Expr | Criteria*/ {
	switch e := expr.(type) {
	case *ast.CallExpr:
		ident, ok := e.Callee.(*ast.IdentExpr)
		util.Assert(ok, "expect ident actual %s", e.Callee)

		callee := strings.ToUpper(ident.Name)
		switch callee {
		case oper.LOGIC_AND, sql.AND:
			lhs, lok := expr2criteria0(e.Args[0]).(Criteria)
			rhs, rok := expr2criteria0(e.Args[1]).(Criteria)
			util.Assert(lok && rok, "expect criteria")
			return CondGroup{
				LogicalOper: AND,
				Conds:       []Criteria{lhs, rhs},
			}
		case oper.LOGIC_OR, sql.OR:
			lhs, lok := expr2criteria0(e.Args[0]).(Criteria)
			rhs, rok := expr2criteria0(e.Args[1]).(Criteria)
			util.Assert(lok && rok, "expect criteria")
			return CondGroup{
				LogicalOper: OR,
				Conds:       []Criteria{lhs, rhs},
			}
		case oper.LOGIC_NOT, sql.NOT:
			lhs, lok := expr2criteria0(e.Args[0]).(Criteria)
			util.Assert(lok, "expect criteria")
			return CondGroup{
				LogicalOper: NOT,
				Conds:       []Criteria{lhs},
			}
		case oper.EQ:
			return makeCond(e, sql.EQ)
		case oper.NE:
			return makeCond(e, sql.NE)
		case oper.GT:
			return makeCond(e, sql.GT)
		case oper.GE:
			return makeCond(e, sql.GE)
		case oper.LT:
			return makeCond(e, sql.LT)
		case oper.LE:
			return makeCond(e, sql.LE)
		default:
			return makeCond(e, callee)
		}
	//case *ast.IdentExpr:
	//case *ast.UnaryExpr:
	//case *ast.BinaryExpr:
	//case *ast.TenaryExpr:
	//case *ast.SubscriptExpr:
	//case *ast.MemberExpr:
	//case *ast.StrExpr:
	//case *ast.NumExpr:
	//case *ast.TimeExpr:
	//case *ast.Criteria:
	//case *ast.ListExpr:
	//case *ast.MapExpr:
	//case *ast.ObjExpr:
	case *ast.GroupExpr:
		return expr2criteria0(e.SubExpr)
	default:
		return e
	}
}

func makeCond(e *ast.CallExpr, oper string) Cond {
	field, ok := e.Args[0].(*ast.IdentExpr)
	util.Assert(ok, "expect ident actual %s", e.Args[0])
	return Cond{
		Field:    field.Name,
		Operator: oper,
		Operands: e.Args[1:],
	}
}
