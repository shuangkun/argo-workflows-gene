package wdl

import (
	"fmt"
	"github.com/antlr/antlr4/runtime/Go/antlr/v4"
	domain "github.com/shuangkun/argo-workflows-gene/pkg/parsers/v1_0"
)

type WdlBuilder struct {
	Url      string
	Version  string
	Document Document
}

type CustomErrorListener struct {
	*antlr.DefaultErrorListener
	Errors []error
}

func NewWdlBuilder(url string) (*WdlBuilder, error) {
	version, err := GetVersion(url)
	if err != nil {
		return nil, err
	}

	if !IsIn(version, []string{"1.0"}) {
		return nil, fmt.Errorf("Unrecognized or Not supported WDL version: %s", version)
	}

	builder := &WdlBuilder{Url: url, Version: version}
	return builder, nil
}

func (this *WdlBuilder) ParseDocument() (*Document, error) {
	data, err := ReadString(this.Url)
	if err != nil {
		return nil, fmt.Errorf("Failed to read from URI: %w", err)
	}

	visitor := NewWdlVisitor(this.Url, this.Version)
	errListener := &CustomErrorListener{}
	var tree domain.IDocumentContext

	switch this.Version {
	case "1.0":
		p := NewWdlV1_0Parser(data)
		p.RemoveErrorListeners()
		p.AddErrorListener(errListener)
		tree = p.Document()
	default:
		return nil, fmt.Errorf("Unknown WDL version: %v", this.Version)
	}

	if len(errListener.Errors) > 0 {
		for _, err := range errListener.Errors {
			fmt.Printf("%v\n", err)
		}

		return nil, fmt.Errorf("SyntaxError detected. Please fix your WDL.")
	}

	document := visitor.VisitDocument(tree)
	return &document, nil
}
