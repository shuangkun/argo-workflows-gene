package wdl

import (
	"fmt"
	"github.com/antlr/antlr4/runtime/Go/antlr/v4"
	wfv1 "github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/shuangkun/argo-workflows-gene/pkg/parsers/v1_0"
)

func GenerateWorkflow(wdlFile string, parameterFile string) wfv1.Workflow {
	builder, err := NewWdlBuilder(wdlFile)
	if err != nil {
		fmt.Printf("Error occurred when importing \"%v\": %v\n", wdlFile, err)
	}
	document, err := builder.ParseDocument()
	if err != nil {
		fmt.Printf("Failed to parse the WDL document. Reason: %v\n", err)
	}
	wf := GetArgoWorkflow(document)
	return wf
}

func NewWdlV1_0Parser(data string) *v1_0.WdlV1Parser {
	input := antlr.NewInputStream(data)
	lexer := v1_0.NewWdlV1Lexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	p := v1_0.NewWdlV1Parser(stream)
	return p
}

func GetArgoWorkflow(document *Document) wfv1.Workflow {
	wf := wfv1.Workflow{}
	if document.Workflow == nil {
		return wf
	}
	wdlwf := document.Workflow
	wf.Name = string(wdlwf.Name)
	return wf
}
