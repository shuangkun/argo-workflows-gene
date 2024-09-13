package wdl

import (
	"fmt"
	"github.com/antlr/antlr4/runtime/Go/antlr/v4"
	antlrv1_0 "github.com/shuangkun/argo-workflows-gene/pkg/parsers/v1_0"
	log "github.com/sirupsen/logrus"
	urlparse "net/url"
	"path"
	"reflect"
	"strings"
)

type WdlVisitor struct {
	Url     string
	Version string
	*antlrv1_0.BaseWdlV1ParserVisitor
	Reporter IVisitorReporter
}

func NewWdlVisitor(url string, version string) *WdlVisitor {
	reporter := &FmtReporter{}
	return &WdlVisitor{
		Url:      url,
		Version:  version,
		Reporter: reporter,
	}
}

type IVisitorReporter interface {
	Warn(antlr.ParserRuleContext, string)
	Error(antlr.ParserRuleContext, error)
	NotImplemented(fmt.Stringer)
}

// Visit takes in any parse tree and visits the child nodes from there.
func (v *WdlVisitor) Visit(tree antlr.ParseTree) interface{} {
	switch t := tree.(type) {
	case antlrv1_0.IDocumentContext:
		return v.VisitDocument(t)
	default:
		return nil
	}
}

func (v *WdlVisitor) VisitChildren(node antlr.RuleNode) interface{}     { return nil }
func (v *WdlVisitor) VisitTerminal(node antlr.TerminalNode) interface{} { return nil }
func (v *WdlVisitor) VisitErrorNode(node antlr.ErrorNode) interface{}   { return nil }

func (v *WdlVisitor) VisitMap_type(ctx antlrv1_0.IMap_typeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *WdlVisitor) VisitArray_type(ctx antlrv1_0.IArray_typeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *WdlVisitor) VisitPair_type(ctx antlrv1_0.IPair_typeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *WdlVisitor) VisitType_base(ctx antlrv1_0.IType_baseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *WdlVisitor) VisitWdl_type(ctx antlrv1_0.IWdl_typeContext) Type {
	// TODO: parse WDL type
	return Type{
		Optional: true,
	}
}

func (v *WdlVisitor) VisitUnbound_decls(ctx antlrv1_0.IUnbound_declsContext) Declaration {
	identifier := Identifier(ctx.GetText())
	var type_ Type

	for _, child := range ctx.GetChildren() {
		switch tt := child.(type) {
		case antlr.TerminalNode:
			identifier = Identifier(tt.GetText())
		case antlrv1_0.IWdl_typeContext:
			type_ = v.VisitWdl_type(tt)
		}
	}

	return Declaration{
		Type:       type_,
		Identifier: identifier,
		Expr:       nil,
	}
}

func (v *WdlVisitor) VisitBound_decls(ctx antlrv1_0.IBound_declsContext) Declaration {
	identifier := Identifier(ctx.GetText())
	var type_ Type
	var expr IExpression

	for _, child := range ctx.GetChildren() {
		switch tt := child.(type) {
		case antlr.TerminalNode:
			if tt.GetText() != "=" {
				identifier = Identifier(tt.GetText())
			}
		case antlrv1_0.IWdl_typeContext:
			type_ = v.VisitWdl_type(tt)
		case antlrv1_0.IExprContext:
			expr = v.VisitExpr(tt)
		}
	}

	return Declaration{
		Type:       type_,
		Identifier: identifier,
		Expr:       &expr,
	}
}

func (v *WdlVisitor) VisitAny_decls(ctx antlrv1_0.IAny_declsContext) Declaration {
	// First child _should_ be the declaration.
	switch tt := ctx.GetChild(0).(type) {
	case antlrv1_0.IBound_declsContext:
		return v.VisitBound_decls(tt)
	case antlrv1_0.IUnbound_declsContext:
		return v.VisitUnbound_decls(tt)
	}

	// This should not happen. Should we panic?
	return Declaration{}
}

func (v *WdlVisitor) VisitNumber(ctx antlrv1_0.INumberContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *WdlVisitor) VisitExpression_placeholder_option(ctx antlrv1_0.IExpression_placeholder_optionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *WdlVisitor) VisitString_part(ctx antlrv1_0.IString_partContext) string {
	var parts []string

	//var p []antlr.TerminalNode
	//log.Infof("The type is %v\n", reflect.TypeOf(coreCtx).Elem().Name())
	for _, part := range ctx.GetChildren() {
		switch tt := part.(type) {
		case antlr.TerminalNode:
			parts = append(parts, tt.GetText())
		}
		log.Infof("The type is %v\n", reflect.TypeOf(part).Elem().Name())
	}

	return strings.Join(parts, "")
}

func (v *WdlVisitor) VisitString_expr_part(ctx antlrv1_0.IString_expr_partContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *WdlVisitor) VisitString_expr_with_string_part(ctx antlrv1_0.IString_expr_with_string_partContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *WdlVisitor) VisitString(ctx antlrv1_0.IStringContext) string { // interface{} {
	// TODO: parse actual string, which can contain expressions!
	for _, child := range ctx.GetChildren() {
		switch tt := child.(type) {
		case antlrv1_0.IString_partContext:
			log.Infof("IString_partContext: %v", tt.GetText())
			return v.VisitString_part(tt)
		case antlrv1_0.IExpression_placeholder_optionContext:
			log.Infof("IExpression_placeholder_optionContext: %v", tt.GetText())
			fmt.Println(reflect.TypeOf(tt))
		default:
			v.Reporter.NotImplemented(reflect.TypeOf(tt))
		}
	}

	// return v.VisitChildren(ctx)
	return ""
}

func (v *WdlVisitor) VisitPrimitive_literal(ctx antlrv1_0.IPrimitive_literalContext) interface{} {
	log.Infof("Primitive start")
	name := ""
	// TODO: figure out the appropriate return type.
	for _, child := range ctx.GetChildren() {
		switch tt := child.(type) {
		case antlrv1_0.INumberContext:
			return v.VisitNumber(tt)
		case antlrv1_0.IStringContext:
			log.Infof("String is: %v", tt.GetText())
			return v.VisitString(tt)
		case antlr.TerminalNode:
			fmt.Println(tt.GetText())
			name = tt.GetText()
		default:
			fmt.Println(reflect.TypeOf(tt))
		}
	}
	fmt.Println("Primitive end")
	return name
}

func (v *WdlVisitor) VisitExpr(ctx antlrv1_0.IExprContext) IExpression {
	infixCtx := ctx.GetChild(0).(antlrv1_0.IExpr_infixContext)
	log.Infof("infixCtx:", infixCtx.GetText())
	// 	expr
	// 	: expr_infix
	// 	;
	//
	//   expr_infix
	// 	: expr_infix0 #infix0
	// 	;

	infix0Ctx := infixCtx.GetChild(0).(antlrv1_0.IExpr_infix0Context)
	log.Infof("infix0Ctx:", infix0Ctx.GetText())
	return v.VisitInfix0(infix0Ctx)
}

func (v *WdlVisitor) VisitInfix0(ctx antlrv1_0.IExpr_infix0Context) IExpression {

	// expr_infix0
	// : expr_infix0 OR expr_infix1 #lor
	// | expr_infix1 #infix1
	// ;

	if count := ctx.GetChildCount(); count == 1 {
		log.Infof("infix0 contains 1 number of children")
		// No branches
		infix1Ctx := ctx.GetChild(0).(antlrv1_0.IExpr_infix1Context)
		return v.VisitInfix1(infix1Ctx)
	} else if count != 3 {
		v.Reporter.Error(ctx, fmt.Errorf("infix0 contains unexpected number of children"))
		return NewUnknownExpr()
	}

	infix0Ctx := ctx.GetChild(0).(antlrv1_0.IExpr_infix0Context) // 0
	infix1Ctx := ctx.GetChild(2).(antlrv1_0.IExpr_infix1Context) // 2

	return NewLorExpr(
		v.VisitInfix0(infix0Ctx), // left
		v.VisitInfix1(infix1Ctx), // right
	)
}

func (v *WdlVisitor) VisitInfix1(ctx antlrv1_0.IExpr_infix1Context) IExpression {

	// expr_infix1
	// : expr_infix1 AND expr_infix2 #land
	// | expr_infix2 #infix2
	// ;

	if count := ctx.GetChildCount(); count == 1 {
		log.Infof("infix1 contains 1 number of children")
		// No branches
		infix2Ctx := ctx.GetChild(0).(antlrv1_0.IExpr_infix2Context)
		return v.VisitInfix2(infix2Ctx)
	} else if count != 3 {
		v.Reporter.Error(ctx, fmt.Errorf("infix1 contains unexpected number of children"))
		return NewUnknownExpr()
	}

	infix1Ctx := ctx.GetChild(0).(antlrv1_0.IExpr_infix1Context) // 0
	infix2Ctx := ctx.GetChild(2).(antlrv1_0.IExpr_infix2Context) // 2

	return NewLandExpr(
		v.VisitInfix1(infix1Ctx), // left
		v.VisitInfix2(infix2Ctx), // right
	)
}

func (v *WdlVisitor) VisitInfix2(ctx antlrv1_0.IExpr_infix2Context) IExpression {

	// A bunch of comparisons:

	// expr_infix2
	// : expr_infix2 EQUALITY expr_infix3 #eqeq
	// | expr_infix2 NOTEQUAL expr_infix3 #neq
	// | expr_infix2 LTE expr_infix3 #lte
	// | expr_infix2 GTE expr_infix3 #gte
	// | expr_infix2 LT expr_infix3 #lt
	// | expr_infix2 GT expr_infix3 #gt
	// | expr_infix3 #infix3
	// ;

	if count := ctx.GetChildCount(); count == 1 {
		log.Infof("infix2 contains 1 number of children")
		// No branches
		infix3Ctx := ctx.GetChild(0).(antlrv1_0.IExpr_infix3Context)
		return v.VisitInfix3(infix3Ctx)
	} else if count != 3 {
		v.Reporter.Error(ctx, fmt.Errorf("infix2 contains unexpected number of children"))
		return NewUnknownExpr()
	}

	infix2Ctx := ctx.GetChild(0).(antlrv1_0.IExpr_infix2Context) // 0
	op := ctx.GetChild(1).(antlr.TerminalNode).GetText()         // 1
	infix3Ctx := ctx.GetChild(2).(antlrv1_0.IExpr_infix3Context) // 2

	// `op` is one of []string{"==", "!=", "<=", ">=", "<", ">"}.  We could further parse this, but string would be fine for now.

	return NewComparisonExpr(
		v.VisitInfix2(infix2Ctx), // left
		op,                       // operation
		v.VisitInfix3(infix3Ctx), // right
	)
}

func (v *WdlVisitor) VisitInfix3(ctx antlrv1_0.IExpr_infix3Context) IExpression {

	// addition and subtraction

	// expr_infix3
	// : expr_infix3 PLUS expr_infix4 #add
	// | expr_infix3 MINUS expr_infix4 #sub
	// | expr_infix4 #infix4
	// ;

	if count := ctx.GetChildCount(); count == 1 {
		log.Infof("infix3 contains 1 number of children")

		// No branches
		infix4Ctx := ctx.GetChild(0).(antlrv1_0.IExpr_infix4Context)
		return v.VisitInfix4(infix4Ctx)
	} else if count != 3 {
		v.Reporter.Error(ctx, fmt.Errorf("infix3 contains unexpected number of children"))
		return NewUnknownExpr()
	}

	infix3Ctx := ctx.GetChild(0).(antlrv1_0.IExpr_infix3Context) // 0
	op := ctx.GetChild(1).(antlr.TerminalNode).GetText()         // 1
	infix4Ctx := ctx.GetChild(2).(antlrv1_0.IExpr_infix4Context) // 2

	// `op` is one of []string{"+", "-"}.  We could further parse this, but string would be fine for now.

	return NewBinaryOpExpr(
		v.VisitInfix3(infix3Ctx), // left
		op,                       // operation
		v.VisitInfix4(infix4Ctx), // right
	)
}

func (v *WdlVisitor) VisitInfix4(ctx antlrv1_0.IExpr_infix4Context) IExpression {

	// multiplication, division, and mod

	// expr_infix4
	// : expr_infix4 STAR expr_infix5 #mul
	// | expr_infix4 DIVIDE expr_infix5 #divide
	// | expr_infix4 MOD expr_infix5 #mod
	// | expr_infix5 #infix5
	// ;

	if count := ctx.GetChildCount(); count == 1 {
		log.Infof("infix1 contains 1 number of children")
		// No branches
		infix5Ctx := ctx.GetChild(0).(antlrv1_0.IExpr_infix5Context)
		return v.VisitInfix5(infix5Ctx)
	} else if count != 3 {
		v.Reporter.Error(ctx, fmt.Errorf("infix4 contains unexpected number of children"))
		return NewUnknownExpr()
	}

	infix4Ctx := ctx.GetChild(0).(antlrv1_0.IExpr_infix4Context) // 0
	op := ctx.GetChild(1).(antlr.TerminalNode).GetText()         // 1
	infix5Ctx := ctx.GetChild(2).(antlrv1_0.IExpr_infix5Context) // 2

	// `op` is one of []string{"*", "/", "%"}.  We could further parse this, but string would be fine for now.

	return NewBinaryOpExpr(
		v.VisitInfix4(infix4Ctx), // left
		op,                       // operation
		v.VisitInfix5(infix5Ctx), // right
	)
}

func (v *WdlVisitor) VisitInfix5(ctx antlrv1_0.IExpr_infix5Context) IExpression {
	coreCtx := ctx.GetChild(0).(antlrv1_0.IExpr_coreContext)

	// expr_infix5
	// : expr_core
	// ;
	//
	// expr_core
	// : Identifier LPAREN (expr (COMMA expr)* COMMA?)? RPAREN #apply
	// | LBRACK (expr (COMMA expr)* COMMA?)* RBRACK #array_literal
	// | LPAREN expr COMMA expr RPAREN #pair_literal
	// | LBRACE (expr COLON expr (COMMA expr COLON expr)* COMMA?)* RBRACE #map_literal
	// | OBJECT_LITERAL LBRACE (Identifier COLON expr (COMMA Identifier COLON expr)* COMMA?)* RBRACE #object_literal
	// | IF expr THEN expr ELSE expr #ifthenelse
	// | LPAREN expr RPAREN #expression_group
	// | expr_core LBRACK expr RBRACK #at
	// | expr_core DOT Identifier #get_name
	// | NOT expr #negate
	// | (PLUS | MINUS) expr #unarysigned
	// | primitive_literal #primitives
	// | Identifier #left_name
	// ;

	// We probably **have** to use reflection here to get the type of context we're dealing with (the label). I don't
	// think ANTLR4 stores this information in the node, it just assumes that we're working with the original interface
	// so we could type check that way, but we're not.
	log.Infof("The type is %v\n", reflect.TypeOf(coreCtx).Elem().Name())
	log.Infof("The text is %v\n", reflect.TypeOf(coreCtx).Elem().String())
	switch reflect.TypeOf(coreCtx).Elem().Name() {
	case "PrimitivesContext":
		rv := NewTerminalExpr()
		rv.Value = v.VisitPrimitive_literal(coreCtx.GetChild(0).(antlrv1_0.IPrimitive_literalContext))
		log.Infof("rv.Value = %v\n", rv.Value)
		return rv
	}

	return NewUnknownExpr()
}

func (v *WdlVisitor) VisitExpr_infix5(ctx antlrv1_0.IExpr_infix5Context) interface{} {
	return v.VisitChildren(ctx)
}

//
//func (v *WdlVisitor) VisitPair_literal(ctx antlrv1_0.Pair_literalContext) interface{} {
//	return v.VisitChildren(ctx)
//}
//
//func (v *WdlVisitor) VisitUnarysigned(ctx antlrv1_0.IUnarysignedContext) interface{} {
//	return v.VisitChildren(ctx)
//}
//
//func (v *WdlVisitor) VisitApply(ctx antlrv1_0.ApplyContext) interface{} {
//	return v.VisitChildren(ctx)
//}
//
//func (v *WdlVisitor) VisitExpression_group(ctx *parsers.Expression_groupContext) interface{} {
//	return v.VisitChildren(ctx)
//}
//
//func (v *WdlVisitor) VisitPrimitives(ctx *parsers.PrimitivesContext) interface{} {
//	return v.VisitChildren(ctx)
//}
//
//func (v *WdlVisitor) VisitLeft_name(ctx *parsers.Left_nameContext) interface{} {
//	return v.VisitChildren(ctx)
//}
//
//func (v *WdlVisitor) VisitAt(ctx *parsers.AtContext) interface{} {
//	return v.VisitChildren(ctx)
//}
//
//func (v *WdlVisitor) VisitNegate(ctx *parsers.NegateContext) interface{} {
//	return v.VisitChildren(ctx)
//}
//
//func (v *WdlVisitor) VisitMap_literal(ctx *parsers.Map_literalContext) interface{} {
//	return v.VisitChildren(ctx)
//}
//
//func (v *WdlVisitor) VisitIfthenelse(ctx *parsers.IfthenelseContext) interface{} {
//	return v.VisitChildren(ctx)
//}
//
//func (v *WdlVisitor) VisitGet_name(ctx *parsers.Get_nameContext) interface{} {
//	return v.VisitChildren(ctx)
//}
//
//func (v *WdlVisitor) VisitObject_literal(ctx *parsers.Object_literalContext) interface{} {
//	return v.VisitChildren(ctx)
//}
//
//func (v *WdlVisitor) VisitArray_literal(ctx *parsers.Array_literalContext) interface{} {
//	return v.VisitChildren(ctx)
//}

// VisitVersion returns the version as a string. This is a terminal operation.
func (v *WdlVisitor) VisitVersion(ctx antlrv1_0.IVersionContext) string {
	version := ctx.GetText()

	if version != "version"+v.Version {
		v.Reporter.Warn(ctx, fmt.Sprintf("Version mismatch! Using parser for version %s but workflow has version %s", v.Version, version))
	}

	return v.Version
}

// VisitImport_alias returns the import alias as a mapping from the original
// identifer -> the aliased identifier. This is a terminal operation.
func (v *WdlVisitor) VisitImport_alias(ctx antlrv1_0.IImport_aliasContext) (Identifier, Identifier) {
	return Identifier(ctx.GetText()[0]),
		Identifier(ctx.GetText()[1])
}

// VisitImport_as returns the identifier that refers to the import. This is a terminal operation.
func (v *WdlVisitor) VisitImport_as(ctx antlrv1_0.IImport_asContext) Identifier {
	return Identifier(ctx.GetText())
}

func (v *WdlVisitor) VisitImport_doc(ctx antlrv1_0.IImport_docContext) Import {
	var url string
	var as Identifier                      // optional
	aliases := map[Identifier]Identifier{} // optional

	for _, child := range ctx.GetChildren() {
		switch tt := child.(type) {
		case antlrv1_0.IStringContext:
			url = v.VisitString(tt)
		case antlrv1_0.IImport_asContext:
			as = v.VisitImport_as(tt)
		case antlrv1_0.IImport_aliasContext:
			original, alias := v.VisitImport_alias(tt)
			aliases[original] = alias
		default:
			v.Reporter.NotImplemented(reflect.TypeOf(tt))
		}
	}

	var absoluteUrl string
	if strings.Contains(url, "://") {
		absoluteUrl = url
	}

	u, err := urlparse.Parse(v.Url)
	if err == nil {
		u.Path = path.Join(path.Dir(u.Path), url)
		absoluteUrl = u.String()
	}

	return Import{
		Url:         url,
		AbsoluteUrl: absoluteUrl,
		As:          as,
		Aliases:     aliases,
	}
}

//func (v *WdlVisitor) VisitStruct(ctx *parsers.StructContext) interface{} {
//	return v.VisitChildren(ctx)
//}
//
//func (v *WdlVisitor) VisitMeta_value(ctx *parsers.Meta_valueContext) interface{} {
//	return v.VisitChildren(ctx)
//}
//
//func (v *WdlVisitor) VisitMeta_string_part(ctx *parsers.Meta_string_partContext) interface{} {
//	return v.VisitChildren(ctx)
//}
//
//func (v *WdlVisitor) VisitMeta_string(ctx *parsers.Meta_stringContext) interface{} {
//	return v.VisitChildren(ctx)
//}
//
//func (v *WdlVisitor) VisitMeta_array(ctx *parsers.Meta_arrayContext) interface{} {
//	return v.VisitChildren(ctx)
//}
//
//func (v *WdlVisitor) VisitMeta_object(ctx *parsers.Meta_objectContext) interface{} {
//	return v.VisitChildren(ctx)
//}
//
//func (v *WdlVisitor) VisitMeta_object_kv(ctx *parsers.Meta_object_kvContext) interface{} {
//	return v.VisitChildren(ctx)
//}
//
//func (v *WdlVisitor) VisitMeta_kv(ctx *parsers.Meta_kvContext) interface{} {
//	return v.VisitChildren(ctx)
//}
//
//func (v *WdlVisitor) VisitParameter_meta(ctx *parsers.Parameter_metaContext) interface{} {
//	return v.VisitChildren(ctx)
//}
//
//func (v *WdlVisitor) VisitMeta(ctx *parsers.MetaContext) interface{} {
//	return v.VisitChildren(ctx)
//}

func (v *WdlVisitor) VisitTask_runtime_kv(ctx antlrv1_0.ITask_runtime_kvContext) map[Identifier]IExpression {
	runtime := make(map[Identifier]IExpression)
	key := Identifier(ctx.GetChild(0).(antlr.TerminalNode).GetText())
	fmt.Println("kv start")
	//var value IExpression
	var value IExpression
	for _, child := range ctx.GetChildren() {
		switch tt := child.(type) {
		case antlrv1_0.IExprContext:
			log.Infof("IExprContext: %v", tt.GetText())
			value = v.VisitExpr(tt)
		default:
			fmt.Println(reflect.TypeOf(tt))
		}
	}
	fmt.Println("kv end")
	runtime[key] = value
	return runtime
}

func (v *WdlVisitor) VisitTask_runtime(ctx antlrv1_0.ITask_runtimeContext) map[Identifier]IExpression {
	var runtime map[Identifier]IExpression
	for _, child := range ctx.GetChildren() {
		switch tt := child.(type) {
		case antlrv1_0.ITask_runtime_kvContext:
			log.Infof("ITask_runtime_kvContext: %v", tt.GetText())
			runtime = v.VisitTask_runtime_kv(tt)
		}
	}
	return runtime
}

func (v *WdlVisitor) VisitTask_input(ctx antlrv1_0.ITask_inputContext) []Declaration {
	len := 0
	for _, ctx := range ctx.GetChildren() {
		if _, ok := ctx.(antlrv1_0.IAny_declsContext); ok {
			len++
		}
	}

	inputs := make([]Declaration, len)
	idx := 0

	for _, child := range ctx.GetChildren() {
		switch tt := child.(type) {
		case antlrv1_0.IAny_declsContext:
			inputs[idx] = v.VisitAny_decls(tt)
			idx++
		}
	}
	return inputs
}

func (v *WdlVisitor) VisitTask_output(ctx antlrv1_0.ITask_outputContext) []Declaration {
	len := 0
	for _, child := range ctx.GetChildren() {
		if _, ok := child.(antlrv1_0.IBound_declsContext); ok {
			len++
		}
	}
	log.Infof("output len is: %v", len)
	outputs := make([]Declaration, len)
	idx := 0

	for _, child := range ctx.GetChildren() {
		switch tt := child.(type) {
		case antlrv1_0.IBound_declsContext:
			log.Infof("IBound_declsContext: %v", tt.GetText())
			outputs[idx] = v.VisitBound_decls(tt)
			idx++
		}
	}
	return outputs
}

func (v *WdlVisitor) VisitTask_command_string_part(ctx antlrv1_0.ITask_command_string_partContext) string {
	var test string
	for _, child := range ctx.GetChildren() {
		switch tt := child.(type) {
		case antlr.TerminalNode:
			fmt.Println(reflect.TypeOf(tt))
			test = test + tt.GetText()
		//case antlrv1_0.IExprContext:
		//	fmt.Println("expr")
		//	test = test + v.VisitExpr(tt).GetType()
		default:
			fmt.Println(reflect.TypeOf(tt))
		}
	}

	return test
}

func (v *WdlVisitor) VisitTask_command_expr_part(ctx antlrv1_0.ITask_command_expr_partContext) string {
	fmt.Println("expr part")
	var exprString string
	for _, child := range ctx.GetChildren() {
		fmt.Println("expr jdfk;ja;f")
		fmt.Println("test", reflect.TypeOf(child))
		switch tt := child.(type) {
		case antlr.TerminalNode:
			fmt.Println("terminal")
			fmt.Println(tt.GetText())
		case antlrv1_0.IExpression_placeholder_optionContext:
			fmt.Println("placeholder")
		case antlrv1_0.IExprContext:
			fmt.Println("expr")
			fmt.Printf("v.VisitExpr(tt): %v", v.VisitExpr(tt))
			name := v.VisitExpr(tt)
			fmt.Printf("name: %v\n", name)
			if name.GetType() == "TerminalExpr" {
				fmt.Println("TerminalExpr")
				//terminal := name.(TerminalExpr)
				//fmt.Printf("terminal.Value: %v\n", terminal.Value)
				//exprString = "$" + "(" + terminal.Value.(string) + ")"
			}
		}
	}
	fmt.Printf("exprString: %v\n", exprString)
	return exprString
}

func (v *WdlVisitor) VisitTask_command_expr_with_string(ctx antlrv1_0.ITask_command_expr_with_stringContext) string {
	fmt.Println("expr with string")
	command := ""
	for _, child := range ctx.GetChildren() {
		fmt.Printf("command: %v\n", command)
		switch tt := child.(type) {
		case antlr.TerminalNode:
			//command = command + tt.GetText()
		case antlrv1_0.ITask_command_string_partContext:
			fmt.Println("expr with string1")
			command = command + v.VisitTask_command_string_part(tt)
		case antlrv1_0.ITask_command_expr_partContext:
			fmt.Println("expr with string2")
			test := v.VisitTask_command_expr_part(tt)
			fmt.Printf("test expetc: %v\n", test)
			command = command + test
		default:
			fmt.Println(reflect.TypeOf(tt))
		}
	}
	fmt.Printf("command: %v\n", command)
	return command
}

func (v *WdlVisitor) VisitTask_command(ctx antlrv1_0.ITask_commandContext) string {
	var command string
	for _, child := range ctx.GetChildren() {
		switch tt := child.(type) {
		case antlrv1_0.ITask_command_string_partContext:
			fmt.Println("command1")
			command += v.VisitTask_command_string_part(tt)
			fmt.Println(command)
		case antlrv1_0.ITask_command_expr_with_stringContext:
			fmt.Println("command3")
			expr := v.VisitTask_command_expr_with_string(tt)
			command += expr
			fmt.Printf("expr all info: %v", expr)
		default:
			fmt.Println(reflect.TypeOf(tt))
		}
	}
	return command
}

func (v *WdlVisitor) VisitTask_element(ctx antlrv1_0.ITask_elementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *WdlVisitor) VisitTask(ctx antlrv1_0.ITaskContext) Task {
	var inputs []Declaration
	var outputs []Declaration
	var runtime map[Identifier]IExpression
	var command string
	var name Identifier

	for _, child := range ctx.GetChildren() {

		switch tt := child.(type) {
		case antlr.TerminalNode:
			name = Identifier(ctx.GetChildren()[1].(antlr.TerminalNode).GetText())
		case antlrv1_0.ITask_elementContext:
			// Let's flatten this.
			fmt.Println("task chdaild")
			// TODO: we need to support expression and declaration first.
			switch tt := tt.GetChild(0).(type) {
			case antlrv1_0.ITask_inputContext:
				fmt.Println("input")
				inputs = v.VisitTask_input(tt)
			case antlrv1_0.ITask_outputContext:
				log.Infof("output: %v", tt.GetText())
				outputs = v.VisitTask_output(tt)
			case antlrv1_0.ITask_runtimeContext:
				log.Infof("runtime chdaild", tt.GetText())
				runtime = v.VisitTask_runtime(tt)
			case antlrv1_0.ITask_commandContext:
				fmt.Println("command1nndknlk")
				command = v.VisitTask_command(tt)
			default:
				v.Reporter.NotImplemented(reflect.TypeOf(tt))
			}
		}
	}

	return Task{
		Name:    name,
		Input:   inputs,
		Command: command,
		Runtime: runtime,
		Output:  outputs,
	}
}

func (v *WdlVisitor) VisitInner_workflow_element(ctx antlrv1_0.IInner_workflow_elementContext) InnerWorkflow {

	var calls []Identifier
	for _, child := range ctx.GetChildren() {
		switch tt := child.(type) {
		case antlrv1_0.ICallContext:
			fmt.Println("call")
			calls = v.VisitCall(tt)
		}
	}

	return InnerWorkflow{
		Calls: calls,
	}
}

func (v *WdlVisitor) VisitCall_alias(ctx antlrv1_0.ICall_aliasContext) Identifier {
	return ""
}

func (v *WdlVisitor) VisitCall_input(ctx antlrv1_0.ICall_inputContext) Identifier {
	return ""
}

func (v *WdlVisitor) VisitCall_inputs(ctx antlrv1_0.ICall_inputsContext) Identifier {
	fmt.Println("call inputs")
	return "Identifier{}"
}

func (v *WdlVisitor) VisitCall_body(ctx antlrv1_0.ICall_bodyContext) Identifier {
	return ""
}

func (v *WdlVisitor) VisitCall_name(ctx antlrv1_0.ICall_nameContext) Identifier {
	var name Identifier
	name = Identifier(ctx.GetChild(0).(antlr.TerminalNode).GetText())
	return name
}

func (v *WdlVisitor) VisitCall(ctx antlrv1_0.ICallContext) []Identifier {
	fmt.Println("call")
	var calls []Identifier
	for _, child := range ctx.GetChildren() {
		switch tt := child.(type) {
		case antlrv1_0.ICall_inputsContext:
			fmt.Println("call inputs")
			calls = append(calls, v.VisitCall_inputs(tt))
		case antlrv1_0.ICall_inputContext:
			fmt.Println("call inputs")
			calls = append(calls, v.VisitCall_input(tt))
		case antlrv1_0.ICall_nameContext:
			calls = append(calls, v.VisitCall_name(tt))
		case antlrv1_0.ICall_bodyContext:
			calls = append(calls, v.VisitCall_body(tt))
		case antlrv1_0.ICall_aliasContext:
			calls = append(calls, v.VisitCall_alias(tt))
		}
	}
	return calls
}

// func (v *WdlVisitor) VisitScatter(ctx *parsers.ScatterContext) interface{} {
// 	return v.VisitChildren(ctx)
// }

// func (v *WdlVisitor) VisitConditional(ctx *parsers.ConditionalContext) interface{} {
// 	return v.VisitChildren(ctx)
// }

func (v *WdlVisitor) VisitWorkflow_input(ctx antlrv1_0.IWorkflow_inputContext) []Declaration {
	len := 0
	for _, ctx := range ctx.GetChildren() {
		if _, ok := ctx.(antlrv1_0.IAny_declsContext); ok {
			len++
		}
	}

	inputs := make([]Declaration, len)
	idx := 0

	for _, child := range ctx.GetChildren() {
		switch tt := child.(type) {
		case antlrv1_0.IAny_declsContext:
			inputs[idx] = v.VisitAny_decls(tt)
			idx++
		}
	}

	return inputs
}

func (v *WdlVisitor) VisitWorkflow_output(ctx antlrv1_0.IWorkflow_outputContext) []Declaration {
	len := 0
	for _, ctx := range ctx.GetChildren() {
		if _, ok := ctx.(antlrv1_0.IAny_declsContext); ok {
			len++
		}
	}

	outputs := make([]Declaration, len)
	idx := 0

	for _, child := range ctx.GetChildren() {
		switch tt := child.(type) {
		case antlrv1_0.IAny_declsContext:
			outputs[idx] = v.VisitAny_decls(tt)
			idx++
		}
	}
	return outputs
}

//func (v *WdlVisitor) VisitInput(ctx *parsers.InputContext) interface{} {
//	return v.VisitChildren(ctx)
//}
//
//func (v *WdlVisitor) VisitOutput(ctx *parsers.OutputContext) interface{} {
//	return v.VisitChildren(ctx)
//}
//
//func (v *WdlVisitor) VisitInner_element(ctx *parsers.Inner_elementContext) interface{} {
//	return v.VisitChildren(ctx)
//}
//
//func (v *WdlVisitor) VisitParameter_meta_element(ctx *parsers.Parameter_meta_elementContext) interface{} {
//	return v.VisitChildren(ctx)
//}
//
//func (v *WdlVisitor) VisitMeta_element(ctx *parsers.Meta_elementContext) interface{} {
//	return v.VisitChildren(ctx)
//}

func (v *WdlVisitor) VisitWorkflow(ctx antlrv1_0.IWorkflowContext) Workflow {

	log.Infof("GetText(): %v", ctx.GetText())

	var inputs []Declaration
	var outputs []Declaration
	//var calls []Calls
	var inners []InnerWorkflow
	var name Identifier
	for _, child := range ctx.GetChildren() {
		switch tt := child.(type) {
		case antlr.TerminalNode:
			name = Identifier(ctx.GetChildren()[1].(antlr.TerminalNode).GetText())
		case antlrv1_0.IWorkflow_elementContext:
			// Let's flatten this.
			// TODO: we need to support expression and declaration first.
			switch tt := tt.GetChild(0).(type) {
			case antlrv1_0.IWorkflow_inputContext:
				log.Infof("IWorkflow_inputContext: %v", tt.GetText())
				inputs = v.VisitWorkflow_input(tt)
			case antlrv1_0.IInner_workflow_elementContext:
				log.Infof("IInner_workflow_elementContext: %v", tt.GetText())
				inners = append(inners, v.VisitInner_workflow_element(tt))
			case antlrv1_0.IWorkflow_outputContext:
				log.Infof("IWorkflow_outputContext: %v", tt.GetText())
				outputs = v.VisitWorkflow_output(tt)
			default:
				v.Reporter.NotImplemented(reflect.TypeOf(tt))
			}
		default:
			log.Infof("tt: %v", tt)
			v.Reporter.NotImplemented(reflect.TypeOf(tt))
		}
	}

	return Workflow{
		Name:           name,
		Inputs:         inputs,
		Output:         outputs,
		InnerWorkflows: inners,
	}
}

func (v *WdlVisitor) VisitDocument_element(ctx antlrv1_0.IDocument_elementContext) interface{} {
	// We've flattened this structure.
	v.Reporter.Warn(ctx, "VisitDocument_element should not be called.")
	return nil
}

func (v *WdlVisitor) VisitDocument(ctx antlrv1_0.IDocumentContext) Document {
	var version string
	var workflow Workflow // optional
	var imports []Import
	var tasks []Task
	log.Infof("VisitDocument: %v", ctx.GetText())
	for _, ctx := range ctx.GetChildren() {
		switch tt := ctx.(type) {
		// There should be exactly one version statement and one or zero workflows.
		case antlrv1_0.IVersionContext:
			log.Infof("IVersionContext: %v", tt.GetText())
			version = v.VisitVersion(tt)
		case antlrv1_0.IInner_workflow_elementContext:
			log.Infof("IInner_workflow_elementContext: %v", tt.GetText())
		// There could be a number of document elements.
		case antlrv1_0.IWorkflowContext:
			log.Infof("IWorkflowContext: %v", tt.GetText())
			workflow = v.VisitWorkflow(tt)
		// There could be a number of document elements.
		case antlrv1_0.IDocument_elementContext:
			log.Infof("IDocument_elementContext: %v", tt.GetText())
			// Let's flatten this.
			switch tt := tt.GetChild(0).(type) {
			case antlrv1_0.IImport_docContext:
				imports = append(imports, v.VisitImport_doc(tt))
			case antlrv1_0.ITaskContext:
				log.Infof("ITaskContext: %v", tt.GetText())
				tasks = append(tasks, v.VisitTask(tt))
			case antlrv1_0.IWorkflowContext:
				log.Infof("IWorkflowContext: %v", tt.GetText())
				workflow = v.VisitWorkflow(tt)
			// TODO: structs
			default:
				v.Reporter.NotImplemented(reflect.TypeOf(tt))
			}
		default:
			v.Reporter.NotImplemented(reflect.TypeOf(tt))
		}
	}

	return Document{
		Url:      v.Url,
		Version:  version,
		Workflow: &workflow,
		Imports:  imports,
		Tasks:    tasks,
	}
}
