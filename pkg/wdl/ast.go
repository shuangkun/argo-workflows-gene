package wdl

import (
	"encoding/json"
	"fmt"
	"regexp"
)

// The Abstract Syntax Tree (AST) of a WDL document. A domain.Document struct
// represents the direct translation of the input WDL document into the AST.
// There is no topological ordering of the dependency graph at this stage yet.

// An identifier should match pattern /[a-zA-Z][a-zA-Z0-9_]+/
// See: https://github.com/openwdl/wdl/blob/main/versions/1.0/SPEC.md#whitespace-strings-identifiers-constants
type Identifier string

func (this Identifier) IsValid() bool {
	match, err := regexp.MatchString(`[a-zA-Z][a-zA-Z0-9_]+`, string(this))
	if err != nil {
		return false
	}

	return match
}

type Document struct {
	// This is not part of WDL, but is useful to keep track of.
	Url string `json:"url"`

	Version  string    `json:"version"`
	Workflow *Workflow `json:"workflow,omitempty"` // optional

	Imports []Import `json:"imports,omitempty"` // optional
	Structs []Struct `json:"structs,omitempty"` // optional
	Tasks   []Task   `json:"tasks,omitempty"`   // optional
}

func (this Document) String() string {
	out, err := json.Marshal(this)
	if err != nil {
		return "Document<failed to serialize>"
	}

	return fmt.Sprintf("Document<%s>", out)
}

type Workflow struct {
	Name          Identifier        `json:"name"`
	Inputs        []Declaration     `json:"inputs,omitempty"`
	InnerWorkflow InnerWorkflow     `json:"calls,omitempty"`
	Output        []Declaration     `json:"outputs,omitempty"`
	Meta          map[string]string `json:"meta,omitempty"`
	ParameterMeta map[string]string `json:"parameter_meta,omitempty"`
}

type InnerWorkflow struct {
	Name        Identifier    `json:"name"`
	Calls       []Identifier  `json:"calls,omitempty"`
	Scatter     []Declaration `json:"scatter,omitempty"`
	Conditional []Declaration `json:"conditional,omitempty"`
}

type Import struct {
	Url         string                    `json:"url"`
	AbsoluteUrl string                    `json:"absoluteUrl,omitempty"` // not part of WDL, but useful when resolving imports
	As          Identifier                `json:"as,omitempty"`          // optional
	Aliases     map[Identifier]Identifier `json:"aliases,omitempty"`     // optional
}

type Task struct {
	Input   []Declaration              `json:"inputs,omitempty"`
	Command string                     `json:"command"`
	Runtime map[Identifier]IExpression `json:"runtime"`
	Output  []Declaration              `json:"outputs,omitempty"`
}

type Struct struct {
}

type Type struct {
	Optional bool `json:"optional,omitempty"`
}

type ArrayType struct {
}

// Declaration can be bounded (assigned to expression) or unbounded (no assignment).
type Declaration struct {
	Type       Type         `json:"type"`
	Identifier Identifier   `json:"identifier"`
	Expr       *IExpression `json:"expr,omitempty"`
}

type Any interface{}

type IExpression interface {
	GetType() string
	//GetValue() string
}

type UnknownExpr struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

func NewUnknownExpr() UnknownExpr {
	return UnknownExpr{
		Type: "UnknownExpr",
	}
}

func (this UnknownExpr) GetType() string {
	return this.Type
}

func (this UnknownExpr) GetValue() string {
	return this.Value
}

// LorExpr is a logic OR expression.
type LorExpr struct {
	Type  string      `json:"type"`
	Left  IExpression `json:"left"`
	Right IExpression `json:"right"`
}

func NewLorExpr(left IExpression, right IExpression) LorExpr {
	return LorExpr{Type: "LorExpr", Left: left, Right: right}
}

func (this LorExpr) GetType() string {
	return this.Type
}

// LandExpr is a logic AND expression.
type LandExpr struct {
	Type  string      `json:"type"`
	Left  IExpression `json:"left"`
	Right IExpression `json:"right"`
}

func NewLandExpr(left IExpression, right IExpression) LandExpr {
	return LandExpr{Type: "LandExpr", Left: left, Right: right}
}

func (this LandExpr) GetType() string {
	return this.Type
}

// ComparisonExpr is a comparison expression.
type ComparisonExpr struct {
	Type      string      `json:"type"`
	Left      IExpression `json:"left"`
	Operation string      `json:"op"`
	Right     IExpression `json:"right"`
}

func NewComparisonExpr(left IExpression, op string, right IExpression) ComparisonExpr {
	return ComparisonExpr{Type: "ComparisonExpr", Left: left, Operation: op, Right: right}
}

func (this ComparisonExpr) GetType() string {
	return this.Type
}

// BinaryOpExpr is a binary operation expression ( + and - ).
type BinaryOpExpr struct {
	Type      string      `json:"type"`
	Left      IExpression `json:"left"`
	Operation string      `json:"op"`
	Right     IExpression `json:"right"`
}

func NewBinaryOpExpr(left IExpression, op string, right IExpression) BinaryOpExpr {
	return BinaryOpExpr{Type: "BinaryOpExpr", Left: left, Operation: op, Right: right}
}

func (this BinaryOpExpr) GetType() string {
	return this.Type
}

type TerminalExpr struct {
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

func NewTerminalExpr() TerminalExpr {
	return TerminalExpr{Type: "TerminalExpr"}
}

func (this TerminalExpr) GetType() string {
	return this.Type
}

func (this TerminalExpr) GetValue() interface{} {
	return this.Value
}
